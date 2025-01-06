package auth

import (
	"backend/config"
	"backend/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	var user models.User
	db := config.GetDB()
	err = db.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", loginReq.Username).
		Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err == sql.ErrNoRows {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println("Ошибка базы данных:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Message: "Успешный вход",
		Token:   "dummy-token",
	})
}
