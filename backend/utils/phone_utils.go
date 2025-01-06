package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// NormalizePhoneNumber нормализует номер телефона к формату +380XXXXXXXXX
func NormalizePhoneNumber(input string) (string, error) {
	// Удаляем все, кроме цифр
	re := regexp.MustCompile(`[^\d]`)
	onlyDigits := re.ReplaceAllString(input, "")

	// Проверяем длину номера
	if len(onlyDigits) < 9 || len(onlyDigits) > 12 {
		return "", fmt.Errorf("Некорректный номер телефона: %s", input)
	}

	// Добавляем код страны, если номер начинается с "0"
	if strings.HasPrefix(onlyDigits, "0") {
		onlyDigits = "38" + onlyDigits
	}

	// Добавляем "+", если номер начинается с "38"
	if strings.HasPrefix(onlyDigits, "38") {
		onlyDigits = "+" + onlyDigits
	}

	// Проверяем итоговый формат
	if len(onlyDigits) != 13 || !strings.HasPrefix(onlyDigits, "+380") {
		return "", fmt.Errorf("Некорректный номер после нормализации: %s", onlyDigits)
	}

	return onlyDigits, nil
}
