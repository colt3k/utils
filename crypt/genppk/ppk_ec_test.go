package genppk

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestPPK_ECPrivateKeys(t *testing.T) {
	p := &PPKEC{PublicFilename: "mypubec", PrivateFilename: "myprivec"}
	p.GenerateECKeys(256)

	pemPrivStr := p.ExportECPrivateKeyAsPemString()
	pemPubStr := p.ExportECPublicKeyAsPemString()
	fmt.Println(pemPrivStr)
	fmt.Println(pemPubStr)

	privKey, err := p.ParsePrivateKeyFromPemStr(pemPrivStr)
	if err != nil {
		fmt.Printf("error parsing priv: %v", err)
		t.FailNow()
	}
	pubKey, err := p.ParsePublicKeyFromPemStr(pemPubStr)
	if err != nil {
		fmt.Printf("error parsing pub: %v", err)
		t.FailNow()
	}

	// ****************************** SIGNING ****************************************************

	// SignPKCS1v15
	message := []byte("This is the message to be signed!")
	hash := sha1.New()
	io.WriteString(hash, string(message))
	hashed := hash.Sum(nil)

	r, s, serr := ecdsa.Sign(rand.Reader, privKey, hashed)
	if serr != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)

	fmt.Printf("Signature : %x\n", signature)

	// ******************************* Verify *********************************************
	// Verify
	verifystatus := ecdsa.Verify(pubKey, hashed, r, s)
	fmt.Println(verifystatus) // should be true

	fmt.Println()
}
