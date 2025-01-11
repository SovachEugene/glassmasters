package main

import (
	"backend/config"
	"backend/handlers/auth"
	"backend/handlers/clients"
	"backend/handlers/translations"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const (
	defaultServerAddr = ":8080" // Default server address
	migrationsUpDir   = "./migrations/migrations_up"
	migrationsDownDir = "./migrations/migrations_down"
)

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

	// Инициализируем базу данных
	config.InitDB()
	defer config.CloseDB()

	// Выполняем миграции
	if err := runMigrations(config.GetDB(), migrationsUpDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

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

// runMigrations выполняет миграции из указанной папки
func runMigrations(db *sql.DB, migrationsDir string) error {
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			filePath := filepath.Join(migrationsDir, file.Name())
			query, err := ioutil.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
			}

			log.Printf("Running migration: %s", file.Name())
			if _, err := db.Exec(string(query)); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
			}
		}
	}

	log.Println("Migrations executed successfully")
	return nil
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
