package driven

type EncryptionService interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}
