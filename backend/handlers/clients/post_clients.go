package clients

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"
)

func CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	var clientRequest models.ClientRequest

	// Декодируем тело запроса
	err := json.NewDecoder(r.Body).Decode(&clientRequest)
	if err != nil {
		log.Println("Ошибка декодирования JSON:", err)
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	// Нормализуем номер телефона
	clientRequest.Phone, err = utils.NormalizePhoneNumber(clientRequest.Phone)
	if err != nil {
		log.Println("Ошибка нормализации номера телефона:", err)
		http.Error(w, "Некорректный номер телефона", http.StatusBadRequest)
		return
	}

	// Сохраняем клиента в базу данных
	db := config.GetDB()
	_, err = db.Exec(
		"INSERT INTO clients (name, email, phone, message) VALUES (?, ?, ?, ?)",
		clientRequest.Name, clientRequest.Email, clientRequest.Phone, clientRequest.Message,
	)
	if err != nil {
		log.Println("Ошибка сохранения клиента:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Клиент успешно добавлен"}`))
}
