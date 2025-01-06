package config

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser создает пользователя с хешированным паролем
func CreateUser(username, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	db := GetDB()
	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, string(passwordHash))
	if err != nil {
		log.Println("Ошибка при создании пользователя:", err)
		return err
	}

	log.Println("Пользователь успешно создан:", username)
	return nil
}
