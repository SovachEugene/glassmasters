package clients

import (
	"backend/config"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Обработчик для изменения поля "processed"
func MarkAsProcessedHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем подключение к базе данных
	db := config.GetDB()
	if db == nil {
		log.Fatal("Подключение к базе данных не инициализировано")
	}

	// Извлекаем ID клиента из параметров URL
	vars := mux.Vars(r)
	clientID := vars["id"]

	// Выполняем SQL-запрос для обновления поля "processed"
	_, err := db.Exec("UPDATE clients SET processed = 1 WHERE id = ?", clientID)
	if err != nil {
		log.Println("Ошибка обновления поля 'processed':", err)
		http.Error(w, "Не удалось обновить клиента", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Клиент успешно обработан",
	})
}
