package auth

import (
	"msn/pkg/common/fault"
	"msn/pkg/utils/crypto"
	"msn/pkg/utils/validation"
	"net/http"
	"time"
)

func ValidateCredentials(email, password string) error {
	if err := validation.ValidateEmail(email); err != nil {
		return fault.NewBadRequest(err.Error())
	}
	if err := validation.ValidatePassword(password); err != nil {
		return fault.NewBadRequest(err.Error())
	}
	return nil
}

func ValidateUser(email, password, hashedPassword string, userDeletedAt *time.Time) error {
	if !crypto.PasswordMatches(password, hashedPassword) {
		return fault.NewUnauthorized("invalid credentials")
	}

	if userDeletedAt != nil {
		return fault.New(
			"user must be active to login",
			fault.WithHTTPCode(http.StatusUnauthorized),
			fault.WithTag(fault.DISABLED_USER),
		)
	}

	return nil
}
