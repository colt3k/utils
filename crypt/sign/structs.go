package sign

import (
	"crypto/cipher"
	"encoding/json"

	"github.com/pkg/errors"
)

// Key holds encryption and message authentication keys for a repository. It is stored
// encrypted and authenticated as a JSON data structure in the Data field of the Key
// structure.
type Key struct {
	MACKey        `json:"mac"`
	EncryptionKey `json:"encrypt"`
}

// EncryptionKey is key used for encryption
type EncryptionKey [32]byte

// MACKey is used to sign (authenticate) data.
type MACKey struct {
	K [16]byte // for AES-128
	R [16]byte // for Poly1305

	masked bool // remember if the MAC key has already been masked
}

// statically ensure that *Key implements crypto/cipher.AEAD
var _ cipher.AEAD = &Key{}

// NonceSize returns the size of the nonce that must be passed to Seal
// and Open.
func (k *Key) NonceSize() int {
	return ivSize
}

type jsonMACKey struct {
	K []byte `json:"k"`
	R []byte `json:"r"`
}

// MarshalJSON converts the MACKey to JSON.
func (m *MACKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonMACKey{K: m.K[:], R: m.R[:]})
}

// UnmarshalJSON fills the key m with data from the JSON representation.
func (m *MACKey) UnmarshalJSON(data []byte) error {
	j := jsonMACKey{}
	err := json.Unmarshal(data, &j)
	if err != nil {
		return errors.Wrap(err, "Unmarshal")
	}
	copy(m.K[:], j.K)
	copy(m.R[:], j.R)

	return nil
}

// MarshalJSON converts the EncryptionKey to JSON.
func (k *EncryptionKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(k[:])
}

// UnmarshalJSON fills the key k with data from the JSON representation.
func (k *EncryptionKey) UnmarshalJSON(data []byte) error {
	d := make([]byte, AESKeyLength)
	err := json.Unmarshal(data, &d)
	if err != nil {
		return errors.Wrap(err, "Unmarshal")
	}
	copy(k[:], d)

	return nil
}
