package genppk

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/colt3k/nglog/ng"
)

type PPKEC struct {
	PrivateFilename string
	PublicFilename  string
	PrivateKey      *ecdsa.PrivateKey
	PublicKey       *ecdsa.PublicKey
	PrivatePEM      string
	privateDAT      []byte
}

func (p *PPKEC) GenerateECKeys(size int) {

	if size == 0 {
		size = 256
	}
	var err error
	var privateKey *ecdsa.PrivateKey
	switch size {
	case 224:
		privateKey, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
		if err != nil {
			ng.Logf(ng.ERROR, "%+v", err)
		}
	case 256:
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			ng.Logf(ng.ERROR, "%+v", err)
		}
	case 384:
		privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			ng.Logf(ng.ERROR, "%+v", err)
		}
	case 521:
		privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			ng.Logf(ng.ERROR, "%+v", err)
		}
	}

	p.PrivateKey = privateKey
	p.PublicKey = &privateKey.PublicKey
}

func (p *PPKEC) ExportECPrivateKeyAsPemString() string {
	privkeyBytes, err := x509.MarshalECPrivateKey(p.PrivateKey)
	if err != nil {
		log.Fatal("issue marshalling ec priv key ", err)
	}

	privkeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	return string(privkeyPEM)
}
func (p *PPKEC) ExportECPublicKeyAsPemString() string {
	PubASN1, err := x509.MarshalPKIXPublicKey(&p.PrivateKey.PublicKey)
	if err != nil {
		ng.Logf(ng.ERROR, "%+v", err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: PubASN1,
	})
	return string(pubBytes)
}
func (p *PPKEC) SavePrivateKeyAsPEM() {
	outFile, err := os.Create(p.PrivateFilename)
	if err != nil {
		ng.Logf(ng.ERROR, "issue creating file %s\n%+v", p.PrivateFilename, err)
	}

	defer outFile.Close()
	privBytes, err := x509.MarshalECPrivateKey(p.PrivateKey)
	if err != nil {
		log.Fatal("issue marshalling ec priv key ", err)
	}
	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	}

	err = pem.Encode(outFile, privateKey)
	if err != nil {
		ng.Logf(ng.ERROR, "issue encoding key\n%+v", err)
	}
}

func (p *PPKEC) SavePublicKeyAsPEM() {
	PubASN1, err := x509.MarshalPKIXPublicKey(&p.PrivateKey.PublicKey)
	if err != nil {
		ng.Logf(ng.ERROR, "%+v", err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: PubASN1,
	})

	ioutil.WriteFile(p.PublicFilename+".pub", pubBytes, 0644)
}

func (p *PPKEC) LoadPrivate() {
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
	priv, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic("failed to parse PKCS1 encoded private key: " + err.Error())
	}
	p.PrivateKey = priv
}

func (p *PPKEC) LoadPublic() {
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
	p.PublicKey = pub.(*ecdsa.PublicKey)
}

// ****************************************
func (p *PPKEC) ParsePrivateKeyFromPemStr(privPEM string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func (p *PPKEC) ParsePublicKeyFromPemStr(pubPEM string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *ecdsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, fmt.Errorf("key type is not EC")
}