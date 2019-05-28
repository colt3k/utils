package cert

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/acme/autocert"

	log "github.com/colt3k/nglog/ng"
)

func LECerts(host, email string) autocert.Manager {

	pwd, err := os.Getwd()
	if err != nil {
		log.Logln(log.ERROR, "unable to find current directory")
	}

	cache := autocert.DirCache(filepath.Join(pwd, "system", "tls", "certs"))
	if _, err := os.Stat(string(cache)); os.IsNotExist(err) {
		err := os.MkdirAll(string(cache), os.ModePerm|os.ModeDir)
		if err != nil {
			log.Logln(log.FATAL, "Unable to create cert directory at", cache)
		}
	}

	if len(host) <= 0 {
		log.Logln(log.FATAL, "no host set, set before obtaining certs")
	}

	log.Println("Using", host, "as host/domain for certificate...")
	log.Println("if the host is not configured properly or is unreachable, set-up will fail")

	if len(email) <= 0 {
		log.Logln(log.FATAL, "email not sent in during make certificates.")
	}
	fmt.Println(email, "as contact email for certificate")

	return autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		Cache:       cache,
		HostPolicy:  autocert.HostWhitelist(host),
		RenewBefore: time.Hour * 24 * 30,
		Email:       email,
	}
}

func ExampleUsage() {
	m := LECerts("localhost", "my@email.com")

	server := &http.Server{
		Addr:      fmt.Sprintf(":%s", "8443"),
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
	}

	// launch http listener for "http-01" ACME challenge
	go http.ListenAndServe(":http", m.HTTPHandler(nil))

	log.Logln(log.FATAL, server.ListenAndServeTLS("", ""))
}
