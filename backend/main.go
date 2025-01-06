package main

import (
	"log"
	"net/http"
	"os"

	"backend/config"
	"backend/handlers/auth"
	"backend/handlers/clients"
	"backend/handlers/translations"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const defaultServerAddr = ":SERVER_ADDR" // Default server address

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load server address from environment variables
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = defaultServerAddr
	}

	// Initialize the database
	config.InitDB()

	// Create routes
	r := mux.NewRouter()

	// API routes
	registerClientRoutes(r)
	registerAuthRoutes(r)
	registerTranslationRoutes(r)

	// Apply CORS middleware
	http.Handle("/", enableCORS(r))

	// Start the server
	log.Printf("Server running at http://localhost%s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

// registerClientRoutes sets up routes for client-related endpoints
func registerClientRoutes(r *mux.Router) {
	r.HandleFunc("/api/clients", clients.GetClientsHandler).Methods("GET")
	r.HandleFunc("/api/clients", clients.CreateClientHandler).Methods("POST")
	r.HandleFunc("/api/clients/{id}/processed", clients.MarkAsProcessedHandler).Methods("POST")
}

// registerAuthRoutes sets up routes for authentication
func registerAuthRoutes(r *mux.Router) {
	r.HandleFunc("/api/login", auth.LoginHandler).Methods("POST")
}

// registerTranslationRoutes sets up routes for translations
func registerTranslationRoutes(r *mux.Router) {
	r.HandleFunc("/api/translations", translations.GetLatestTranslationHandler).Methods("GET")
}

// enableCORS adds CORS headers to responses
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
