package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

type AESEncryptionService struct {
	key []byte
}

func CreateAESEncryptionService(key string) *AESEncryptionService {
	return &AESEncryptionService{key: []byte(key)}
}

func (s *AESEncryptionService) Encrypt(plainText string) (string, error) {
	c, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce, err := createNonceWithRandomBytes(gcm.NonceSize())
	if err != nil {
		return "", err
	}

	cipheredText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipheredText), nil
}

func createNonceWithRandomBytes(size int) ([]byte, error) {
	nonce := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return nonce, nil
}

func (s *AESEncryptionService) Decrypt(cipherBase64 string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, bareCipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, bareCipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
