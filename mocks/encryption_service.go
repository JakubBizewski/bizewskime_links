package mocks

type MockEncryptionService struct{}

func (s *MockEncryptionService) Encrypt(plainText string) (string, error) {
	return plainText, nil
}

func (s *MockEncryptionService) Decrypt(cipherText string) (string, error) {
	return cipherText, nil
}
