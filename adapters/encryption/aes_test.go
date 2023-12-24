package encryption_test

import (
	"testing"

	"github.com/JakubBizewski/jakubme_links/adapters/encryption"
)

func TestAESEncryptionService(t *testing.T) {
	existingCipheredText := "FHpJ/KgdGWRxtLUfZjZSA3HDkpU6VcsQ6ntoTYEDtvWz9saMjHCZHAKSOtO4M0xohZe45YS3ZarVCAbM9RPpcFY="
	existingPlainText := "this text was encrypted on 24.12.2023"

	encryptionKey := "12345678901234567890123456789012"

	encryptionService := encryption.CreateAESEncryptionService(encryptionKey)

	t.Run("EncryptsAndDecrypts", func(t *testing.T) {
		plainText := "hello world"

		cipherText, err := encryptionService.Encrypt(plainText)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if plainText == cipherText {
			t.Errorf("Expected %s to be different than %s", plainText, cipherText)
		}

		decryptedText, err := encryptionService.Decrypt(cipherText)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if plainText != decryptedText {
			t.Errorf("Expected %s, got %s", plainText, decryptedText)
		}
	})

	t.Run("DecryptsExistingCipheredText", func(t *testing.T) {
		decryptedText, err := encryptionService.Decrypt(existingCipheredText)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if existingPlainText != decryptedText {
			t.Errorf("Expected %s, got %s", existingPlainText, decryptedText)
		}
	})
}
