package security

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func NewHasher() *BcryptHasher { return &BcryptHasher{} }

func (h *BcryptHasher) Hash(ctx context.Context, password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (h *BcryptHasher) Compare(ctx context.Context, password, passwordHash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	return false, err
}
