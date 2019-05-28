package encrypt

//Crypter interface for custom cryptography
type Crypter interface {
	Encrypt() []byte
	Decrypt() []byte
	Validate() bool
}

type KeyGenerator interface {
	Generate() []byte
}