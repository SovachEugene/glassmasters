package clients

import (
	"backend/config"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func GetClientsHandler(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()

	rows, err := db.Query("SELECT id, name, email, phone, message, city, country, processed FROM clients where processed =1")
	if err != nil {
		log.Println("Ошибка выполнения SQL-запроса:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var clients []map[string]interface{}
	for rows.Next() {
		var id int
		var name, email, phone, message string
		var city, country sql.NullString
		var processed bool

		err = rows.Scan(&id, &name, &email, &phone, &message, &city, &country, &processed)
		if err != nil {
			log.Println("Ошибка обработки строки:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		client := map[string]interface{}{
			"id":        id,
			"name":      name,
			"email":     email,
			"phone":     phone,
			"message":   message,
			"city":      nullStringToString(city),
			"country":   nullStringToString(country),
			"processed": processed,
		}
		clients = append(clients, client)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
