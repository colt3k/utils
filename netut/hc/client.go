package hc

/*
https://github.com/jsha/minica - use for self generation of certs
go get github.com/jsha/minica

Example:
# Generate a root key and cert in minica-key.pem, and minica.pem, then
# generate and sign an end-entity key and cert, storing them in ./foo.com/
$ minica --domains foo.com


Complete GUIDE to TIMEOUTS
https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
*/

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/colt3k/utils/mathut"

	"github.com/colt3k/nglog/ers/bserr"
	log "github.com/colt3k/nglog/ng"

	"github.com/colt3k/utils/netut/nettools"
)

type Client struct {
	httpClient            *http.Client
	DialTimeout           time.Duration
	DialKeepAliveTimeout  time.Duration
	MaxIdleConnections    int
	IdleConnTimeout       time.Duration
	TlsHandshakeTimeout   time.Duration
	ResponseHeaderTimeout time.Duration
	//ExpectContinueTimeout    time.Duration	// will disable HTTP2 if used
	HttpClientRequestTimeout time.Duration
	disableVerifyCert        bool
}

type Auth struct {
	Username []byte
	Password []byte
}

type ClientCert struct {
	Certificate string
	Key         string
}

func NewClient(opts ...ClientOption) *Client {
	t := new(Client)
	t.DialTimeout = 30 * time.Second
	t.DialKeepAliveTimeout = 30 * time.Second
	t.MaxIdleConnections = 100
	t.IdleConnTimeout = 90 * time.Second
	t.TlsHandshakeTimeout = 10 * time.Second
	t.ResponseHeaderTimeout = 10 * time.Second
	//t.ExpectContinueTimeout = 1 * time.Second
	t.HttpClientRequestTimeout = 30 * time.Second

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func (c *Client) Fetch(method, url string, auth *Auth, header map[string]string, body io.Reader) (*http.Response, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.disableVerifyCert,
	}
	// Test for HTTP_PROXY and HTTPS_PROXY and use appropriate one
	var netTransport = &http.Transport{
		Proxy:http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   c.DialTimeout, // time spent establishing a TCP connection
			KeepAlive: c.DialKeepAliveTimeout,
			//DualStack: true,		// now set by default and deprecated
		}).DialContext,
		MaxIdleConns:        c.MaxIdleConnections,
		IdleConnTimeout:     c.IdleConnTimeout,
		TLSHandshakeTimeout: c.TlsHandshakeTimeout, // time spent performing the TLS handshake
		//ExpectContinueTimeout: c.ExpectContinueTimeout, //time client will wait between sending request headers and receiving the go-ahead to send the body
		ResponseHeaderTimeout: c.ResponseHeaderTimeout, //time spent reading the headers of the response
		TLSClientConfig:       tlsConfig,
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout:   c.HttpClientRequestTimeout, //entire exchange, from Dial to reading the body
			Transport: netTransport,
		}
	}

	// Can be used instead of all timers to perform cancel based on time set for the client
	//https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	//ctx, cancel := context.WithCancel(context.Background())
	//timer := time.AfterFunc(5*time.Second, func() {
	//	cancel()
	//})

	req, _ := http.NewRequest(method, url, body)
	//req = req.WithContext(ctx)
	req.Close = true
	if auth != nil {
		req.SetBasicAuth(string(auth.Username), string(auth.Password))
	}
	// Add any required headers.
	for key, value := range header {
		log.Logf(log.DEBUG, "adding header setting %s=%s", key, value)
		req.Header.Add(key, value)

		if key == "Content-Length" {
			req.ContentLength = mathut.ParseInt(value)
		}
		if key == "Content-Type" {
			req.Header.Set(key, value)
		}
	}

	// Disabled due to spitting out contents of uploaded files
	//if log.IsDebug() {
	//	dump, _ := httputil.DumpRequestOut(req, true)
	//	fmt.Println(string(dump))
	//}
	// Perform said network call.
	res, err := c.httpClient.Do(req)
	if err != nil {
		//glog.Error(err.Error()) // use glog it's amazing
		return nil, err
	}

	// If response from network call is not 200, return error too.
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		return res, errors.New(res.Status)
	}
	return res, nil
}

func (c *Client) FetchTLS(method, url string, auth Auth, header map[string]string, body io.Reader, serverCAPath string, cert ClientCert) (*http.Response, error) {

	cp, _ := x509.SystemCertPool()
	data, _ := ioutil.ReadFile(serverCAPath)
	cp.AppendCertsFromPEM(data)

	config := &tls.Config{
		// Certificates: []tls.Certificate{c},
		RootCAs:               cp,
		GetClientCertificate:  nettools.ClientCertReqFunc(cert.Certificate, cert.Key),
		VerifyPeerCertificate: nettools.CertificateChains,
		InsecureSkipVerify:    c.disableVerifyCert,
	}

	var netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   c.DialTimeout,
			KeepAlive: c.DialKeepAliveTimeout,
			//DualStack: true,		// now set by default and deprecated
		}).DialContext,
		MaxIdleConns:        c.MaxIdleConnections,
		IdleConnTimeout:     c.IdleConnTimeout,
		TLSHandshakeTimeout: c.TlsHandshakeTimeout,
		//ExpectContinueTimeout: c.ExpectContinueTimeout,
		TLSClientConfig: config,
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout:   c.HttpClientRequestTimeout,
			Transport: netTransport,
		}
	}

	req, _ := http.NewRequest(method, url, body)
	req.Close = true
	if len(strings.TrimSpace(string(auth.Username))) > 0 {
		req.SetBasicAuth(string(auth.Username), string(auth.Password))
	}
	// Add any required headers.
	for key, value := range header {
		req.Header.Add(key, value)
	}
	// Perform said network call.
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	// If response from network call is not 200, return error too.
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		return res, errors.New(res.Status)
	}
	return res, nil
}

// ProxiedClient Proxy should be set with os.Setenv
func (c *Client) ProxiedClient() *http.Client {
	proxy := os.Getenv("http_proxy")
	proxys := os.Getenv("https_proxy")
	if len(strings.TrimSpace(proxy)) == 0 {
		proxy = proxys
	}
	if len(proxy) > 0 && !strings.HasPrefix(proxy, "http") {
		log.Logf(log.WARN, "proxy should have a prefix of http:// or https:// if utilized.\n Your Proxy: '%s'", proxy)
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.disableVerifyCert,
	}
	netTransport := &http.Transport{}
	if len(proxy) > 0 {
		proxyURL, _ := url.Parse(proxy)
		netTransport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   c.DialTimeout,
				KeepAlive: c.DialKeepAliveTimeout,
				//DualStack: true,		// now set by default and deprecated
			}).DialContext,
			MaxIdleConns:        c.MaxIdleConnections,
			IdleConnTimeout:     c.IdleConnTimeout,
			TLSHandshakeTimeout: c.TlsHandshakeTimeout,
			//ExpectContinueTimeout: c.ExpectContinueTimeout,
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: tlsConfig,
		}
	} else {
		netTransport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   c.DialTimeout,
				KeepAlive: c.DialKeepAliveTimeout,
				//DualStack: true,		// now set by default and deprecated
			}).DialContext,
			MaxIdleConns:        c.MaxIdleConnections,
			IdleConnTimeout:     c.IdleConnTimeout,
			TLSHandshakeTimeout: c.TlsHandshakeTimeout,
			//ExpectContinueTimeout: c.ExpectContinueTimeout,
			TLSClientConfig: tlsConfig,
		}
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout:   c.HttpClientRequestTimeout,
			Transport: netTransport,
		}
	}
	return c.httpClient
}

// Reuse if possible
var httpClient *Client
var responseTimeout int

// Reachable is the url reachable
func Reachable(host, name string, timeout int, disableVerifyCert bool) (bool, error) {
	if httpClient == nil || responseTimeout != timeout {
		responseTimeout = timeout
		httpClient = NewClient(HttpClientRequestTimeout(responseTimeout), DisableVerifyClientCert(disableVerifyCert))
	}
	httpClient.disableVerifyCert = disableVerifyCert
	resp, err := httpClient.Fetch("GET", host, nil, nil, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	// 202 occurs when a http.DELETE is ran
	if err != nil {
		if strings.Index(err.Error(), "Client.Timeout ") > -1 {
			return false, errors.New("site unreachable: " + name)
		}
		return false, fmt.Errorf("site unreachable\n%+v", err.Error())
	}
	// Read body to buffer
	body, err := ioutil.ReadAll(resp.Body)
	if bserr.Err(err, "Error reading body") {
		return false, errors.New("unable to read response")
	}

	if body != nil && len(body) > 0 {
		//log.Println(string(body))
		return true, nil
	} else if body != nil && len(body) == 0 {
		// no body but reachable
		return true, nil
	}
	return false, nil
}
