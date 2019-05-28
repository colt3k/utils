package genppk

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestPPK_GetPrivateKey(t *testing.T) {
	p := &PPK{PublicFilename:"mypub", PrivateFilename:"mypriv"}
	p.GenerateKeys(0)

	pemPrivStr := p.ExportPrivateKeyAsPemString()
	pemPubStr := p.ExportPublicKeyAsPemString()
	fmt.Println(pemPrivStr)
	fmt.Println(pemPubStr)

	privKey,err := p.ParsePrivateKeyFromPemStr(pemPrivStr)
	if err != nil {
		fmt.Printf("error parsing priv: %v", err)
		t.FailNow()
	}
	pubKey, err := p.ParsePublicKeyFromPemStr(pemPubStr)
	if err != nil {
		fmt.Printf("error parsing pub: %v", err)
		t.FailNow()
	}

	// ********************* EncryptOAEP (Optimal Asymmetric Encryption Padding) ************************
	msg := []byte("The secret message!")
	label := []byte("")
	sha1hash := sha1.New()
	encryptedmsg, err := rsa.EncryptOAEP(sha1hash, rand.Reader, pubKey, msg, label)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Printf("OAEP encrypted [%s] to \n[%x]\n", string(msg), encryptedmsg)
	fmt.Println()

	// DecryptOAEP
	decryptedmsg, err := rsa.DecryptOAEP(sha1hash, rand.Reader, privKey, encryptedmsg, label)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("OAEP decrypted [%x] to \n[%s]\n", encryptedmsg, decryptedmsg)
	fmt.Println()

	// ****************************** PKCS1v15 - Deemed less secure than OAEP ***************************
	// EncryptPKCS1v15
	encryptedPKCS1v15, errPKCS1v15 := rsa.EncryptPKCS1v15(rand.Reader, pubKey, msg)
	if errPKCS1v15 != nil {
		fmt.Println(errPKCS1v15)
		os.Exit(1)
	}
	fmt.Printf("PKCS1v15 encrypted [%s] to \n[%x]\n", string(msg), encryptedPKCS1v15)
	fmt.Println()
	// DecryptPKCS1v15
	decryptedPKCS1v15, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, encryptedPKCS1v15)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("PKCS1v15 decrypted [%x] to \n[%s]\n", encryptedPKCS1v15, decryptedPKCS1v15)
	fmt.Println()

	// ****************************** SIGNING ****************************************************

	// SignPKCS1v15
	var h crypto.Hash
	message := []byte("This is the message to be signed!")
	hash := sha1.New()
	io.WriteString(hash, string(message))
	hashed := hash.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, h, hashed)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("PKCS1v15 Signature : %x\n", signature)


	// ******************************* VerifyPKCS1v15 *********************************************
	err = rsa.VerifyPKCS1v15(pubKey, h, hashed, signature)
	if err != nil {
		fmt.Println("VerifyPKCS1v15 failed")
		os.Exit(1)
	} else {
		fmt.Println("VerifyPKCS1v15 successful")
	}
	fmt.Println()


	//p.ShowPrivateKey()
	//p.SavePrivateKeyAsPEM()
	//p.SavePublicKeyAsPEM()

	// Load and decode keys
	//p.LoadPrivate()
	//p.LoadPublic()
}
