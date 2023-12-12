package model

import "math/rand"

const randomCodeCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ShortLink struct {
	ShortCode string
	TargetUrl string
}

func NewRandomShortLink(targetUrl string, shortCodeLen int) ShortLink {
	return ShortLink{
		ShortCode: generateRandomShortCode(shortCodeLen),
		TargetUrl: targetUrl,
	}
}

func generateRandomShortCode(shortCodeLen int) string {
	randomCode := make([]byte, shortCodeLen)
	for i := range randomCode {
		randomCode[i] = randomCodeCharset[rand.Intn(len(randomCodeCharset))]
	}

	return string(randomCode)
}
