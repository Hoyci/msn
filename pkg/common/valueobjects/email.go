package valueobjects

import (
	"errors"
	"net/mail"
)

type Email struct {
	Value string
}

func NewEmail(value string) (Email, error) {
	if _, err := mail.ParseAddress(value); err != nil {
		return Email{}, errors.New("invalid email format")
	}
	return Email{Value: value}, nil
}
