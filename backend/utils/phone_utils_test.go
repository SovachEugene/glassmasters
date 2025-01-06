package utils

import "testing"

func TestNormalizePhoneNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"0982276859", "+380982276859", false},
		{"+38(098) 227 68 59", "+380982276859", false},
		{"380982276859", "+380982276859", false},
		{"invalid", "", true},
	}

	for _, test := range tests {
		result, err := NormalizePhoneNumber(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("Для входных данных '%s' ожидалась ошибка: %v, но получено: %v", test.input, test.hasError, err)
		}
		if result != test.expected {
			t.Errorf("Для входных данных '%s' ожидался результат '%s', но получено '%s'", test.input, test.expected, result)
		}
	}
}
