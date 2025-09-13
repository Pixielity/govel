package interfaces

type EncrypterInterface interface {
	Encrypt(value interface{}, serialize bool) (string, error)
	EncryptString(value string) (string, error)
	Decrypt(payload string, unserialize bool) (interface{}, error)
	DecryptString(payload string) (string, error)
}