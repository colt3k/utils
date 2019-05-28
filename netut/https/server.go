package https

/*
https://github.com/jsha/minica - use for self generation of certs
go get github.com/jsha/minica

Example:
# Generate a root key and cert in minica-key.pem, and minica.pem, then
# generate and sign an end-entity key and cert, storing them in ./foo.com/
$ minica --domains foo.com


COMPLETE GUIDE TO TIMEOUTS
https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
*/
import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/colt3k/nglog/ng"
	"github.com/gorilla/mux"

	"github.com/colt3k/utils/netut/nettools"
)

type ServerCert struct {
	Certificate string
	Key         string
}

func testGET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "good")
}

//HandleHTTPRequests handles requests over http
func HandleHTTPRequests(host, port string) {

	if len(host) > 0 && len(port) > 0 {
		port := ":" + port
		myRouter := mux.NewRouter().StrictSlash(true)
		myRouter.HandleFunc("/", testGET).Methods("GET")
		log.Logln(log.INFO, "HTTP Server Listening on "+host+":"+port)
		log.Logln(log.FATAL, http.ListenAndServe(port, myRouter))
	}
}

//HandleHTTPRequest handle http request
func HandleHTTPRequest(host string) {

	hostdata := strings.Split(host, ":")
	if len(hostdata) > 1 {
		port := ":" + hostdata[1]
		myRouter := mux.NewRouter().StrictSlash(true)
		myRouter.HandleFunc("/", testGET).Methods("GET")
		log.Logln(log.INFO, "HTTP Server Listening on "+hostdata[0]+":"+hostdata[1])
		log.Logln(log.FATAL, http.ListenAndServe(port, myRouter))
	}
}

// https://github.com/dlsniper/gopherconuk
func New(mux http.Handler, serverAddress string, skipVerify bool) *http.Server {
	// See https://blog.cloudflare.com/exposing-go-on-the-internet/ for details
	// about these settings
	tlsConfig := &tls.Config{
		InsecureSkipVerify: skipVerify,
		// Causes servers to use Go's default cipher suite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},

		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
	var srv *http.Server
	if mux != nil {
		srv = &http.Server{
			Addr:         serverAddress,
			ReadTimeout:  5 * time.Second,  //when the connection is accepted to when the request body is fully read
			WriteTimeout: 10 * time.Second, //time from the end of the request header read to the end of the response write
			IdleTimeout:  120 * time.Second,
			TLSConfig:    tlsConfig,
			Handler:      mux,
		}
	} else {
		// setup by default to send http to https if no mux
		srv = &http.Server{
			Addr:         serverAddress,
			ReadTimeout:  5 * time.Second,  //when the connection is accepted to when the request body is fully read
			WriteTimeout: 10 * time.Second, //time from the end of the request header read to the end of the response write
			IdleTimeout:  120 * time.Second,
			TLSConfig:    tlsConfig,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Connection", "close")
				url := "https://" + req.Host + req.URL.String()
				http.Redirect(w, req, url, http.StatusMovedPermanently)
			}),
		}
	}
	return srv
}

func NewWithCustCA(mux http.Handler, serverAddress string, clientCAPath string, serverCert ServerCert) *http.Server {
	cp := x509.NewCertPool()
	data, _ := ioutil.ReadFile(clientCAPath)
	cp.AppendCertsFromPEM(data)

	// See https://blog.cloudflare.com/exposing-go-on-the-internet/ for details
	// about these settings
	tlsConfig := &tls.Config{
		// Causes servers to use Go's default cipher suite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},

		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		ClientCAs:             cp,
		ClientAuth:            tls.RequireAndVerifyClientCert,
		GetCertificate:        nettools.CertReqFunc(serverCert.Certificate, serverCert.Key),
		VerifyPeerCertificate: nettools.CertificateChains,
	}
	srv := &http.Server{
		Addr:         serverAddress,
		ReadTimeout:  5 * time.Second,  //when the connection is accepted to when the request body is fully read
		WriteTimeout: 10 * time.Second, //time from the end of the request header read to the end of the response write
		IdleTimeout:  120 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      mux,
	}
	return srv
}

//func LimitAmount(limit int64) func(next http.Handler) http.Handler {
//	return func(next http.Handler) {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			defer r.Body.Close()
//			r.Body = http.MaxBytesReader(w, r.Body, limit)
//			next.ServeHTTP(w, r)
//		})
//	}
//}
