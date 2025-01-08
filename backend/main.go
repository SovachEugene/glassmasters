package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"backend/config"
	"backend/handlers/auth"
	"backend/handlers/clients"
	"backend/handlers/translations"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const defaultServerAddr = ":8080" // Default server address

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Получаем адрес сервера из переменных окружения
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = defaultServerAddr
	}

	// Проверяем строку подключения к базе данных
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		log.Fatal("DB_DSN is not set in the environment")
	}

	// Инициализируем базу данных (но без миграций)
	config.InitDB()

	// Создаем маршрутизатор
	r := mux.NewRouter()

	// Регистрируем маршруты
	registerClientRoutes(r)
	registerAuthRoutes(r)
	registerTranslationRoutes(r)

	// Добавляем CORS middleware
	http.Handle("/", enableCORS(r))

	// Логируем запуск сервера
	log.Printf("Server running at http://localhost%s", serverAddr)

	// Канал для обработки сигнала завершения работы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		if err := http.ListenAndServe(serverAddr, nil); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-stop
	log.Println("Shutting down server...")
}

// registerClientRoutes устанавливает маршруты для работы с клиентами
func registerClientRoutes(r *mux.Router) {
	r.HandleFunc("/api/clients", clients.GetClientsHandler).Methods("GET")
	r.HandleFunc("/api/clients", clients.CreateClientHandler).Methods("POST")
	r.HandleFunc("/api/clients/{id}/processed", clients.MarkAsProcessedHandler).Methods("POST")
}

// registerAuthRoutes устанавливает маршруты для работы с аутентификацией
func registerAuthRoutes(r *mux.Router) {
	r.HandleFunc("/api/login", auth.LoginHandler).Methods("POST")
}

// registerTranslationRoutes устанавливает маршруты для работы с переводами
func registerTranslationRoutes(r *mux.Router) {
	r.HandleFunc("/api/translations", translations.GetLatestTranslationHandler).Methods("GET")
}

// enableCORS добавляет CORS-заголовки в ответы
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
