package cert

import (
	"testing"

	log "github.com/colt3k/nglog/ng"
)

func TestNew(t *testing.T) {

	cert := New(".", "Roadrunner", "localhost/127.0.0.1", "", "P256", nil, true, 2048)
	cert.CreateCert(false)

	_, certif := cert.readKeyCert(KeyPath, CertPath)
	log.Logln(log.INFO, "DNSNames:", certif.DNSNames)
	log.Logln(log.INFO, "Email:", certif.EmailAddresses)
	log.Logln(log.INFO, "IPs:", certif.IPAddresses)
	log.Logln(log.INFO, "IsCA:", certif.IsCA)
	log.Logln(log.INFO, "NotBefore:", certif.NotBefore)
	log.Logln(log.INFO, "NotAfter:", certif.NotAfter)
	log.Logln(log.INFO, "Version:", certif.Version)
	log.Logln(log.INFO, "Issuer:", certif.Issuer)
	log.Logln(log.INFO, "PublicKeyAlgorithm:", certif.PublicKeyAlgorithm)
	log.Logln(log.INFO, "SerialNumber:", certif.SerialNumber)
	log.Logln(log.INFO, "SignatureAlgorithm:", certif.SignatureAlgorithm)
	log.Logln(log.INFO, "Subject:", certif.Subject)
}
