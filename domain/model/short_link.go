package model

import "math/rand"

const randomCodeCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ShortLink struct {
	ShortCode string
	TargetURL string
}

func CreateRandomShortLink(targetURL string, shortCodeLen int) ShortLink {
	return ShortLink{
		ShortCode: generateRandomShortCode(shortCodeLen),
		TargetURL: targetURL,
	}
}

//nolint:gosec // we don't need cryptographically secure random numbers here
func generateRandomShortCode(shortCodeLen int) string {
	randomCode := make([]byte, shortCodeLen)
	for i := range randomCode {
		randomCode[i] = randomCodeCharset[rand.Intn(len(randomCodeCharset))]
	}

	return string(randomCode)
}
