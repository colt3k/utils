package h_mac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

const VaultHeader = "$ANSIBLE_VAULT;1.1;AES256"

type Key struct {
	CipherKey []byte
	HMACKey   []byte
	IV        []byte
}
type Secret struct {
	Salt []byte
	HMAC []byte
	Data []byte
}

func Hash(hmacKey, data []byte) []byte {
	// Hash the secret content
	hash := hmac.New(sha256.New, hmacKey)
	hash.Write(data)
	hashSum := hash.Sum(nil)

	return hashSum
}

func EncodeSecretAnsibleVault(secret *Secret, key *Key) (string, error) {
	hmacEncrypt := hmac.New(sha256.New, key.HMACKey)
	hmacEncrypt.Write(secret.Data)
	hexSalt := hex.EncodeToString(secret.Salt)
	hexHmac := hmacEncrypt.Sum(nil)
	hexCipher := hex.EncodeToString(secret.Data)

	combined := strings.Join([]string{
		string(hexSalt),
		hex.EncodeToString([]byte(hexHmac)),
		string(hexCipher),
	}, "\n")

	result := strings.Join([]string{
		VaultHeader,
		wrapText(hex.EncodeToString([]byte(combined))),
	}, "\n")

	return result, nil
}

func DecodeSecretAnsibleVault(input string) (*Secret, error) {

	lines := strings.SplitN(input, "\n", 3)
	if len(lines) != 3 {
		return nil, errors.New("invalid secret")
	}

	salt, err := hex.DecodeString(lines[0])
	if err != nil {
		return nil, err
	}

	hmac, err := hex.DecodeString(lines[1])
	if err != nil {
		return nil, err
	}

	data, err := hex.DecodeString(lines[2])
	if err != nil {
		return nil, err
	}

	return &Secret{salt, hmac, data}, nil
}


func wrapText(text string) string {
	src := []byte(text)
	result := []byte{}

	for i := 0; i < len(src); i++ {
		if i > 0 && i%80 == 0 {
			result = append(result, '\n')
		}
		result = append(result, src[i])
	}

	return string(result)
}

func HexDecode(input string) (string, error) {
	input = strings.TrimSpace(input)
	input = strings.Replace(input, "\r", "", -1)
	input = strings.Replace(input, "\n", "", -1)

	decoded, err := hex.DecodeString(input)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

func CheckDigest(secret *Secret, key *Key) error {
	hash := hmac.New(sha256.New, key.HMACKey)
	hash.Write(secret.Data)
	if !hmac.Equal(hash.Sum(nil), secret.HMAC) {
		return errors.New("invalid password")
	}
	return nil
}