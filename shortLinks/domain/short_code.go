package domain

import (
	"errors"
	"math/rand"
)

const randomCodeCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	ErrEmptyShortCode = errors.New("Empty short code")
)

type ShortCode string

func CreateRandomShortCode(length int) ShortCode {
	randomCode := make([]byte, length)
	for i := range randomCode {
		randomCode[i] = randomCodeCharset[rand.Intn(len(randomCodeCharset))]
	}

	return ShortCode(randomCode)
}

func CreateShortCode(shortCode string) (ShortCode, error) {
	if shortCode == "" {
		return "", ErrEmptyShortCode
	}

	return ShortCode(shortCode), nil
}
