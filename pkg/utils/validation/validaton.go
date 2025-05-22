package validation

import (
	"fmt"
	"net/mail"
	"unicode"
)

func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	var hasUpper, hasLower, hasNumber, hasSymbol bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsSymbol(c), unicode.IsPunct(c):
			hasSymbol = true
		}
	}
	if !hasUpper || !hasLower || !hasNumber || !hasSymbol {
		return fmt.Errorf("password must contain uppercase, lowercase, number and symbol")
	}
	return nil
}
