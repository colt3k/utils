package nettools

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// CertReqFunc returns a function for tlsConfig.GetCertificate
func CertReqFunc(certfile, keyfile string) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	c, err := getCert(certfile, keyfile)

	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		fmt.Printf("Received TLS Hello asking for %s: sending certificate\n", hello.ServerName)
		if err != nil || certfile == "" {
			fmt.Println("I have no certificate")
		} else {
			err := OutputPEMFile(certfile)
			if err != nil {
				log.Printf("\n%+v", err)
			}
		}
		Wait()
		return &c, nil
	}
}

// ClientCertReqFunc returns a function for tlsConfig.GetClientCertificate
func ClientCertReqFunc(certfile, keyfile string) func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
	c, err := getCert(certfile, keyfile)

	return func(certReq *tls.CertificateRequestInfo) (*tls.Certificate, error) {
		fmt.Println("Received certificate request: sending certificate")
		if err != nil || certfile == "" {
			fmt.Println("I have no certificate")
		} else {
			err := OutputPEMFile(certfile)
			if err != nil {
				log.Printf("\n%+v", err)
			}
		}
		Wait()
		return &c, nil
	}
}

func getCert(certfile, keyfile string) (c tls.Certificate, err error) {
	if certfile != "" && keyfile != "" {
		c, err = tls.LoadX509KeyPair(certfile, keyfile)
		if err != nil {
			fmt.Printf("Error loading key pair: %v\n", err)
		}
	} else {
		err = fmt.Errorf("I have no certificate")
	}
	return
}

// OutputPEMFile reads info from a PEM file and displays it
func OutputPEMFile(filename string) error {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	for len(data) > 0 {
		var block *pem.Block
		block, data = pem.Decode(data)
		fmt.Printf("Type: %#v\n", block.Type)
		switch block.Type {
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return err
			}
			fmt.Printf(CertificateInfo(cert))
		default:
			fmt.Println(block.Type)
		}
	}
	return nil
}

// CertificateInfo returns a string describing the certificate
func CertificateInfo(cert *x509.Certificate) string {
	if cert.Subject.CommonName == cert.Issuer.CommonName {
		return fmt.Sprintf("    Self-signed certificate %v\n", cert.Issuer.CommonName)
	}

	s := fmt.Sprintf("    Subject %v\n", cert.DNSNames)
	s += fmt.Sprintf("    Issued by %s\n", cert.Issuer.CommonName)
	return s
}

// CertificateChains prints information about verified certificate chains
func CertificateChains(rawCerts [][]byte, chains [][]*x509.Certificate) error {
	if len(chains) > 0 {
		fmt.Println("Verified certificate chain from peer:")

		for _, v := range chains {
			// fmt.Printf("Chain %d:\n", j)
			for i, cert := range v {
				fmt.Printf("  Cert %d:\n", i)
				fmt.Printf(CertificateInfo(cert))
			}
		}
		Wait()
	}

	return nil
}

// Wait holds up proceedings until the user presses return
func Wait() {
	fmt.Printf("[Press enter to proceed]")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	fmt.Println()
}