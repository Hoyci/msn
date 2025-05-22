package valueobjects

import (
	"errors"
	"msn/pkg/utils/crypto"
)

type Password struct {
	Hash string
}

func NewPassword(plaintext string) (Password, error) {
	if len(plaintext) < 8 {
		return Password{}, errors.New("password must be at least 8 characters")
	}
	hash, err := crypto.HashPassword(plaintext)
	if err != nil {
		return Password{}, err
	}
	return Password{Hash: hash}, nil
}

func (p Password) Matches(plaintext string) bool {
	return crypto.PasswordMatches(plaintext, p.Hash)
}
