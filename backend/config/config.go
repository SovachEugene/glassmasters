package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv загружает переменные окружения из файла .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found: %v", err)
	}
}

// GetEnv получает значение переменной окружения или возвращает значение по умолчанию
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
