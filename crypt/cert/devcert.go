// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'certs/cert.pem' and 'certs/key.pem' and will overwrite existing files.

// some parts from https://github.com/jsha/minica/blob/master/main.go
package cert

/*
Generate self cert for local development testing on TLS
*/
import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/netut"
)

var (
	workPath = ""
	KeyPath  = ""
	CertPath = ""
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "ECDSA PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

type Cert struct {
	path       string        //Optional certificate path after '$pwd/deployments/tls/', (default '$pwd/deployments/tls/certs')
	org        string        //Optional organization, 											(default Acme Co)
	domains    []string      //Optional domain, 												(default localhost/127.0.0.1)
	validFrom  string        //Creation date formatted as Jan 1 15:04:05 2011					(default time.Now())
	ecdsaCurve string        //ECDSA curve to use to generate a key. Valid values are P224, P256 (recommended), P384, P521
	validFor   time.Duration //Duration that certificate is valid for 							(default 1 yr)
	isCA       bool          //whether this cert should be its own Certificate Authority		(default false)
	rsaBits    int           //Size of RSA key to generate. Ignored if ecdsaCurve is set		(default 2048)
}

/*
New Certificate/Key
	path...			where to store once created
	organization...	"Acme Co"
	domainIPCombo..	"localhost/127.0.0.1"
	validFrom...	"Jan 1 15:04:05 2011"
	ecdsaCurve...	empty for RSA otherwise P256 recommended
	validFor...		some time.X Duration
	isCA...			become it's own certificate authority
	rsaBits...		default 2048, ignored if using ecdsaCurve instead of RSA
*/
func New(path, organization, domainIPCombo, validFrom, ecdsaCurve string, validFor *time.Duration, isCA bool, rsaBits int) *Cert {
	t := new(Cert)
	if len(path) > 0 {
		t.path = path
	} else {
		t.path = "certs"
	}
	if len(organization) > 0 {
		t.org = organization
	} else {
		t.org = "Acme Co"
	}
	if len(domainIPCombo) > 0 {
		t.domains = strings.Split(strings.TrimSpace(domainIPCombo), ",")
	}
	if len(validFrom) > 0 {
		t.validFrom = validFrom
	}
	t.ecdsaCurve = ecdsaCurve
	if validFor != nil {
		t.validFor = *validFor
	} else {
		t.validFor = 365 * 24 * time.Hour
	}
	t.isCA = isCA
	if rsaBits != 0 {
		t.rsaBits = rsaBits
	} else {
		t.rsaBits = 2048
	}
	// setup path
	pwd, err := os.Getwd()
	if err != nil {
		log.Logln(log.ERROR, "unable to find current directory")
	}
	tlsPath := filepath.Join(pwd, "deployments", "tls")
	workPath = filepath.Join(tlsPath, t.path)

	CertPath = filepath.Join(workPath, "cert.pem")
	KeyPath = filepath.Join(workPath, "key.pem")

	return t
}

/*
CreateCert used to create self signed certificates
domain			Comma-separated hostnames and IPs to generate a certificate for (default localhost/127.0.0.1)
certValidFrom	Creation date formatted as Jan 1 15:04:05 2011 					(default time.now() )
certValidFor	Duration that certificate is valid for 							(default 1 yr)
rsaBits			Size of RSA key to generate. Ignored if --ecdsa-curve is set	(default 2048)
ecdsaCurve		ECDSA curve to use to generate a key. Valid values are P224, P256 (recommended), P384, P521
isCA			whether this cert should be its own Certificate Authority		(default false)

*/
func (c *Cert) CreateCert(override bool) {

	_, keyErr := ioutil.ReadFile(KeyPath)
	_, certErr := ioutil.ReadFile(CertPath)
	if (os.IsNotExist(keyErr) && os.IsNotExist(certErr)) || override {
		priv := c.generatePrivateKey()
		derBytes := c.generateCertificate(priv)

		c.cleanUp(workPath)
		c.writeOutKey(priv)
		c.writeOutCertificate(derBytes)
	}
}
func (c *Cert) generatePrivateKey() interface{} {
	var err error
	var privateKey interface{}

	switch c.ecdsaCurve {
	// default RSA
	case "":
		privateKey, err = rsa.GenerateKey(rand.Reader, c.rsaBits)
		// Elliptic options
	case "P224":
		privateKey, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized elliptic curve: %q", c.ecdsaCurve)
		os.Exit(1)
	}
	if err != nil {
		log.Logf(log.FATAL, "failed to generate private key: %s", err)
	}

	return privateKey
}
func (c *Cert) generateCertificate(priv interface{}) []byte {
	var err error
	// Create initial certificate time as of validFrom passed in
	var notBefore time.Time
	if len(c.validFrom) == 0 {
		notBefore = time.Now()
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", c.validFrom)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse creation date: %s\n", err)
			os.Exit(1)
		}
	}

	// Create end date for certificate
	notAfter := notBefore.Add(c.validFor)

	// Generiate Serial Number
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Logf(log.FATAL, "failed to generate serial number: %s", err)
	}

	domains, parsedIPs, err := netut.ParseIPs(c.domains)
	if err != nil {
		panic(err)
	}
	// Setup a template for an x509.Certificate in order to populate
	templateX509Cert := x509.Certificate{
		DNSNames:     domains,
		IPAddresses:  parsedIPs,
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{c.org},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	// Set as Certificate Authority
	if c.isCA {
		templateX509Cert.IsCA = true
		templateX509Cert.KeyUsage |= x509.KeyUsageCertSign
	}

	// Create x509 Certificate, output to der encoded form
	derBytes, err := x509.CreateCertificate(rand.Reader, &templateX509Cert, &templateX509Cert, publicKey(priv), priv)
	if err != nil {
		log.Logf(log.FATAL, "failed to create certificate: %s", err)
	}
	return derBytes
}
func (c *Cert) cleanUp(path string) {
	// Delete directory and children
	err := os.RemoveAll(path)
	if err != nil {
		log.Logf(log.FATAL, "failed to remove cert directory %v", err)
	}

	// Recreate Directory
	err = os.MkdirAll(path, os.ModeDir|os.ModePerm)
	if err != nil {
		log.Logf(log.FATAL, "failed to create cert directory %v", err)
	}
}
func (c *Cert) writeOutCertificate(derBytes []byte) {
	// Write out Certificate
	certFile, err := os.Create(CertPath)
	if err != nil {
		log.Logf(log.FATAL, "failed to open cert.pem for writing: %v", err)
	}
	// Encode CERTIFICATE
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certFile.Close()
	log.Print("wrote cert.pem\n")
}
func (c *Cert) writeOutKey(priv interface{}) {
	// Write out Key
	keyFile, err := os.OpenFile(KeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("failed to open key.pem for writing:", err)
		return
	}
	// Encode Key
	pem.Encode(keyFile, pemBlockForKey(priv))
	keyFile.Close()
	log.Print("wrote key.pem\n")
}

func (c *Cert) readKeyCert(keyPath, certPath string) (crypto.Signer, *x509.Certificate) {
	keyContents, keyErr := ioutil.ReadFile(keyPath)
	certContents, certErr := ioutil.ReadFile(certPath)
	if os.IsNotExist(keyErr) && os.IsNotExist(certErr) {
		log.Logln(log.FATAL, "neither key nor cert exist")
	} else if keyErr != nil {
		log.Logf(log.FATAL, "%s not available but %s exists)", keyErr, certPath)
	} else if certErr != nil {
		log.Logf(log.FATAL, "%s not available but %s exists", certErr, keyPath)
	}
	key, err := readPrivateKeyContents(keyContents)
	if err != nil {
		log.Logf(log.FATAL, "reading private key %s: %s", keyPath, err)
	}
	cert, err := readCertificateContents(certContents)
	if err != nil {
		log.Logf(log.FATAL, "reading certificate %s: %s", certPath, err)
	}
	equal, err := publicKeysTestEquality(key.Public(), cert.PublicKey)
	if err != nil {
		log.Logf(log.FATAL, "comparing public keys: %s", err)
	} else if !equal {
		log.Logf(log.FATAL, "public key in CA certificate %s doesn't match private key in %s",
			certPath, keyPath)
	}
	return key, cert
}

func readPrivateKeyContents(keyContents []byte) (crypto.Signer, error) {
	block, _ := pem.Decode(keyContents)
	if block == nil {
		return nil, fmt.Errorf("no PEM found")
	} else if block.Type != "RSA PRIVATE KEY" && block.Type != "ECDSA PRIVATE KEY" {
		return nil, fmt.Errorf("incorrect PEM type %s", block.Type)
	}
	if block.Type == "RSA PRIVATE KEY" {
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	} else if block.Type == "ECDSA PRIVATE KEY" {
		return x509.ParseECPrivateKey(block.Bytes)
	}
	return nil, fmt.Errorf("read private key, other error")
}

func readCertificateContents(certContents []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certContents)
	if block == nil {
		return nil, fmt.Errorf("no PEM found")
	} else if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("incorrect PEM type %s", block.Type)
	}
	return x509.ParseCertificate(block.Bytes)
}
func publicKeysTestEquality(a, b interface{}) (bool, error) {
	aBytes, err := x509.MarshalPKIXPublicKey(a)
	if err != nil {
		return false, err
	}
	bBytes, err := x509.MarshalPKIXPublicKey(b)
	if err != nil {
		return false, err
	}
	return bytes.Compare(aBytes, bBytes) == 0, nil
}
