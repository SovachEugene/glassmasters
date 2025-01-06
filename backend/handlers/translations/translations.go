package translations

import (
	"backend/config"
	"log"
	"net/http"
)

func GetLatestTranslationHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр `lang` из URL
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		http.Error(w, "Language not specified", http.StatusBadRequest)
		return
	}

	// Подключаемся к базе данных
	db := config.GetDB()

	// Запрос для получения последнего JSON по `lang`
	query := `
        SELECT json_data 
        FROM translations 
        WHERE lang = ? 
        ORDER BY id DESC 
        LIMIT 1
    `

	var jsonData string
	err := db.QueryRow(query, lang).Scan(&jsonData)
	if err != nil {
		log.Printf("Error fetching translation for lang %s: %v", lang, err)
		http.Error(w, "Error fetching translation", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type и возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}
