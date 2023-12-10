package domain

import (
	"errors"
	"strings"
)

var (
	ErrEmptyUrl   = errors.New("Empty URL")
	ErrInvalidUrl = errors.New("Invalid URL")
)

type Url string

func CreateUrl(url string) (Url, error) {
	if url == "" {
		return "", ErrEmptyUrl
	}

	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		return "", ErrInvalidUrl
	}

	return Url(url), nil
}
