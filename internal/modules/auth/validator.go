package auth

import (
	user "msn/internal/modules/user"
	"msn/pkg/common/fault"
	"msn/pkg/utils/crypto"
	"net/http"
	"net/mail"
	"unicode"
)

func ValidateCredentials(email, password string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fault.NewBadRequest("invalid email format")
	}

	if len(password) < 8 {
		return fault.NewBadRequest("password must be at least 8 characters")
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
		return fault.NewBadRequest("password must contain uppercase, lowercase, numbers and symbol")
	}

	return nil
}

func ValidateUser(email, password string, user *user.User) error {
	if !crypto.PasswordMatches(password, user.Password()) {
		return fault.NewUnauthorized("invalid credentials")
	}

	if user.DeletedAt() != nil {
		return fault.New(
			"user must be active to login",
			fault.WithHTTPCode(http.StatusUnauthorized),
			fault.WithTag(fault.DISABLED_USER),
		)
	}

	return nil
}
