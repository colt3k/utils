package aesvault

//https://docs.ansible.com/ansible/2.4/vault.html#vault-format
import (
	"fmt"
	"strings"
	"testing"

	"github.com/colt3k/utils/hash/h_mac"

	"github.com/colt3k/utils/crypt"
)

var (
	cipherText = `$ANSIBLE_VAULT;1.1;AES256
37373736393734333235616532653938373163376463303032656336313366303431356338643665
3036646236346132303039646337373137383065323363360a313736333037333465386361623365
34666433366364666464343230356131393465353165343634643135393035313235323534636434
6335373432633963360a623833626161353366313366633864333764666532663232316437373462
3336`
	password = "password"

	a *AES
)

func init() {
	a = &AES{
		AESKeyLength:   32,
		SaltLength:     32,
		IVLength:       16,
		IterationCount: 10000,
		Plaintext:      "secret",
	}
}
func TestAES_EncryptAnsibleVault(t *testing.T) {

}

func TestAES_DecryptAnsibleVault(t *testing.T) {

}

func encrypt() {
	salt := crypt.GenSalt(nil, a.SaltLength)

	k := PBKDF2KeySHA256(&password, salt, a)
	a.CipherKey = k.CipherKey
	a.HMACKey = k.HMACKey
	a.IV = k.IV
	crypted, err := a.Encrypt()
	if err != nil {
		panic(err)
	}

	// Hash the secret content
	hashSum := h_mac.Hash(a.HMACKey, crypted)

	// Encode the secret payload
	s,err := h_mac.EncodeSecretAnsibleVault(&h_mac.Secret{Data: crypted, Salt: salt, HMAC: hashSum}, k)
	if err != nil {
		panic(err)
	}

	a.CipherText = string(s)
	fmt.Printf("Encrypted: |%s|\n", s)
}

func decrypt() {

	fmt.Printf("Decrypting: |%s|\n",a.CipherText)
	lines := strings.Split(string(a.CipherText), "\n")
	// Valid secret must include header and body
	if len(lines) < 2 {
		fmt.Println("less than 2 lines")
	}

	// Validate the vault file format
	if strings.TrimSpace(lines[0]) != h_mac.VaultHeader {
		fmt.Println("invalid vault header")
	}

	decoded, err := h_mac.HexDecode(strings.Join(lines[1:], "\n"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decoded: |%s|\n",decoded)

	secret, err := h_mac.DecodeSecretAnsibleVault(decoded)
	if err != nil {
		panic(err)
	}

	key := PBKDF2KeySHA256(&password, secret.Salt, a)
	a.CipherKey = key.CipherKey
	a.HMACKey = key.HMACKey
	a.IV = key.IV
	a.CipherText = string(secret.Data)
	if err := h_mac.CheckDigest(secret, key); err != nil {
		panic(err)
	}

	s, err := a.Decrypt()
	if err != nil {
		panic(err)
	}

	fmt.Println("plaintext", s)
}

func TestAES_Validate(t *testing.T) {
	encrypt()
	decrypt()
}