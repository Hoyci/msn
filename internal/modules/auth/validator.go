package auth

import (
	user "msn/internal/modules/user"
	"msn/pkg/common/fault"
	"msn/pkg/utils/crypto"
	"msn/pkg/utils/validation"
	"net/http"
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

func ValidateUser(email, password string, user *user.User) error {
	if !crypto.PasswordMatches(password, user.Password) {
		return fault.NewUnauthorized("invalid credentials")
	}

	if user.DeletedAt != nil {
		return fault.New(
			"user must be active to login",
			fault.WithHTTPCode(http.StatusUnauthorized),
			fault.WithTag(fault.DISABLED_USER),
		)
	}

	return nil
}
