package genppk

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/colt3k/nglog/ng"
	"golang.org/x/crypto/ssh"
)

type PPK struct {
	PrivateFilename string
	PublicFilename  string
	PrivateKey      *rsa.PrivateKey
	PublicKey       *rsa.PublicKey
	PrivatePEM      string
	privateDAT      []byte
}

func (p *PPK) GenerateKeys(size int) error {
	if size == 0 {
		size = 2048
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		ng.Logf(ng.ERROR, "%+v", err)
	}

	//Validate Key
	err = privateKey.Validate()
	if err != nil {
		return err
	}

	p.PrivateKey = privateKey
	p.PublicKey = &privateKey.PublicKey

	return nil
}

func (p *PPK) ExportPrivateKeyAsPemString() string {
	privkeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(p.PrivateKey),	// ASN.1 DER format
		},
	)
	return string(privkeyPEM)
}
func (p *PPK) ExportPublicKeyAsPemString() string {
	PubASN1, err := x509.MarshalPKIXPublicKey(&p.PrivateKey.PublicKey)
	if err != nil {
		ng.Logf(ng.ERROR, "%+v", err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: PubASN1,
	})
	return string(pubBytes)
}

func (p *PPK) ShowPrivateKey() {

	var byt bytes.Buffer
	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(p.PrivateKey),
	}

	err := pem.Encode(&byt, privateKey)
	if err != nil {
		ng.Logf(ng.ERROR, "issue encoding key\n%+v", err)
	}
	p.PrivatePEM = byt.String()
	//ng.Logf(ng.INFO, "private: %s", byt.String())
}
func (p *PPK) SavePrivateKeyAsPEM() {
	outFile, err := os.Create(p.PrivateFilename)
	if err != nil {
		ng.Logf(ng.ERROR, "issue creating file %s\n%+v", p.PrivateFilename, err)
	}

	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Headers:nil,
		Bytes: x509.MarshalPKCS1PrivateKey(p.PrivateKey),
	}

	err = pem.Encode(outFile, privateKey)
	if err != nil {
		ng.Logf(ng.ERROR, "issue encoding key\n%+v", err)
	}
}

func (p *PPK) SavePublicKey() {
	pubKey := &p.PrivateKey.PublicKey
	pub, _ := ssh.NewPublicKey(pubKey)
	err := ioutil.WriteFile(p.PublicFilename+".pub", ssh.MarshalAuthorizedKey(pub), 0777)
	if err != nil {
		ng.Logf(ng.ERROR, "issue writing public key\n%+v", err)
	}
}

func (p *PPK) SavePublicKeyAsPEM() {
	PubASN1, err := x509.MarshalPKIXPublicKey(&p.PrivateKey.PublicKey)
	if err != nil {
		ng.Logf(ng.ERROR, "%+v", err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: PubASN1,
	})

	ioutil.WriteFile(p.PublicFilename+".pub", pubBytes, 0644)
}

func (p *PPK) LoadPrivate() {
	f, err := os.Open(p.PrivateFilename)
	if err != nil {
		ng.Logf(ng.FATAL, "error opening file")
	}
	defer f.Close()

	dat, err := ioutil.ReadAll(f)
	if err != nil {
		ng.Logf(ng.FATAL, "error reading file")
	}

	block, _ := pem.Decode(dat)
	if block == nil {
		panic("failed to parse PEM block containing the private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic("failed to parse PKCS1 encoded private key: " + err.Error())
	}
	p.PrivateKey = priv
}

func (p *PPK) LoadPublic() {
	f, err := os.Open(p.PublicFilename + ".pub")
	if err != nil {
		ng.Logf(ng.FATAL, "error opening file")
	}
	defer f.Close()

	dat, err := ioutil.ReadAll(f)
	if err != nil {
		ng.Logf(ng.FATAL, "error reading file")
	}

	block, _ := pem.Decode(dat)
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}

	//switch pub := pub.(type) {
	//case *rsa.PublicKey:
	//	fmt.Println("pub is of type RSA:", pub)
	//case *dsa.PublicKey:
	//	fmt.Println("pub is of type DSA:", pub)
	//case *ecdsa.PublicKey:
	//	fmt.Println("pub is of type ECDSA:", pub)
	//default:
	//	panic("unknown type of public key")
	//}
	p.PublicKey = pub.(*rsa.PublicKey)
}

// ****************************************
func (p *PPK) ParsePrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func (p *PPK) ExportPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)

	return string(pubkey_pem), nil
}

func (p *PPK) ParsePublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, fmt.Errorf("key type is not RSA")
}

func (p *PPK) RSAEncrypt(origData []byte) ([]byte, error) {

	//fmt.Println("Modulus : ", p.PublicKey.N.String())
	//fmt.Println(">>> ", p.PublicKey.N)
	//fmt.Printf("Modulus(Hex) : %X\n", p.PublicKey.N)
	//fmt.Println("Public Exponent : ", p.PublicKey.E)
	return rsa.EncryptPKCS1v15(rand.Reader, p.PublicKey, origData)
}

func (p *PPK) RsaDecrypt(ciphertext []byte) ([]byte, error) {

	return rsa.DecryptPKCS1v15(rand.Reader, p.PrivateKey, ciphertext)
}
