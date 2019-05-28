package bcryp

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type BCrypPass struct {
	pass []byte
	hash []byte
}

func (b *BCrypPass) Encrypt() []byte {
	bs, err := bcrypt.GenerateFromPassword(b.pass, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}
	return bs
}

func (b *BCrypPass) Validate() bool {
	err := bcrypt.CompareHashAndPassword(b.hash, b.pass)
	if err != nil {
		return false
	}
	return true
}
func (b *BCrypPass) Decrypt() []byte {

	return nil
}

func New(int, int, int, int, *string, []byte, []byte) {

}
