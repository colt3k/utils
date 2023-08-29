package hc

import (
	"github.com/colt3k/nglog/ers/bserr"
	log "github.com/colt3k/nglog/ng"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var client *Client

// HTTPClient data store
type HTTPClient struct {
	Method                  string
	URL                     string
	Header                  map[string]string
	ReturnHead              bool
	ReturnKeys              []string
	ReturnCert              bool
	ReUseClient             bool
	RedirectUseLastResponse bool
	Auth                    *Auth
	RequestTimeout          int
	ResponseHeaderTimeout   int
	DialTimeout             int
	DialKeepAliveTimeout    int
	MaxIdleConnections      int
	IdleConnectionTimeout   int
	TLSHandshakeTimeout     int
	DisableVerifyClientCert bool
	CatchAllErrStatus       int
}

// NewHTTPClient initialize HTTPClient
func NewHTTPClient(method, url string, header map[string]string, auth *Auth, settings *HTTPClientSettings) *HTTPClient {
	t := new(HTTPClient)
	t.Method = method
	t.URL = url
	t.Header = header
	t.Auth = auth
	if settings != nil {
		t.ReturnHead = settings.ReturnHeaders
		t.RequestTimeout = settings.RequestTimeout
		t.ResponseHeaderTimeout = settings.ResponseHeaderTimeout
		t.DisableVerifyClientCert = settings.DisableVerifyClientCert
		t.ReturnCert = settings.ReturnCert
		t.ReUseClient = settings.ReUseClient
		t.RedirectUseLastResponse = settings.RedirectUseLastResponse
		if settings.DialTimeout == 0 {
			t.DialTimeout = 30
		} else {
			t.DialTimeout = settings.DialTimeout
		}
		if settings.DialKeepAliveTimeout == 0 {
			t.DialKeepAliveTimeout = 30
		} else {
			t.DialKeepAliveTimeout = settings.DialKeepAliveTimeout
		}
		if settings.MaxIdleConnections == 0 {
			t.MaxIdleConnections = 100
		} else {
			t.MaxIdleConnections = settings.MaxIdleConnections
		}
		if settings.IdleConnectionTimeout == 0 {
			t.IdleConnectionTimeout = 90
		} else {
			t.IdleConnectionTimeout = settings.IdleConnectionTimeout
		}
		if settings.TLSHandshakeTimeout == 0 {
			t.TLSHandshakeTimeout = 10
		} else {
			t.TLSHandshakeTimeout = settings.TLSHandshakeTimeout
		}
		if settings.CatchAllErrStatus == 0 {
			t.CatchAllErrStatus = 1000
		}
	} else {
		t.ReturnHead = true
		t.DisableVerifyClientCert = true
		t.RequestTimeout = 3600
		t.ResponseHeaderTimeout = 3600
		t.ReturnCert = true
		t.ReUseClient = true
		t.RedirectUseLastResponse = false
		t.DialTimeout = 30
		t.DialKeepAliveTimeout = 30
		t.MaxIdleConnections = 100
		t.IdleConnectionTimeout = 90
		t.TLSHandshakeTimeout = 10
		t.CatchAllErrStatus = 1000
	}
	return t
}

// Process method on HTTPCall object, pass in reader
func (h *HTTPClient) Process(data io.Reader) (map[string]interface{}, int, error) {
	log.Logf(log.DBGL2, "-- called httpCallData for Method %s URL: %s", h.Method, h.URL)
	tmp := make(map[string]interface{})

	if client == nil || !h.ReUseClient {
		//log.Logln(log.DEBUG, "!!! Creating NEW HTTP CLIENT !!!")
		// set to timeout after a day per request, accommodates file uploads
		client = NewClient(HTTPClientRequestTimeout(h.RequestTimeout), DisableVerifyClientCert(h.DisableVerifyClientCert),
			HTTPClientRequestTimeout(h.ResponseHeaderTimeout), CheckRedirectUserLastResp(h.RedirectUseLastResponse),
			DialTimeout(h.DialTimeout), DialKeepAliveTimeout(h.DialKeepAliveTimeout), MaxIdleConnections(h.MaxIdleConnections),
			IdleConnectionTimeout(h.IdleConnectionTimeout), TLSHandshakeTimeout(h.TLSHandshakeTimeout))
	}
	t := time.Now()
	log.Logf(log.DBGL3, "Start Fetch %v", t.Format(time.RFC1123))
	var resp, err = client.Fetch(h.Method, h.URL, h.Auth, h.Header, data)
	log.Logf(log.DBGL3, "Post Fetch %v", time.Since(t))

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil && err.Error() != "201 Created" && err.Error() != "204 No Content" {
		if resp != nil && resp.Body != nil {
			body, errRA := io.ReadAll(resp.Body)
			if errRA != nil {
				return nil, resp.StatusCode, err
			}
			if resp.Header != nil {
				for name, value := range resp.Header {
					tmp[name] = value
				}
			}
			tmp["body"] = string(body)
			if h.ReturnCert && resp.TLS != nil {
				certificates := resp.TLS.PeerCertificates
				if len(certificates) > 0 {
					// you probably want certificates[0]
					cert := certificates[0]
					tmp["cert"] = cert
				}
			}
			return tmp, resp.StatusCode, err
		} else if resp != nil {
			return nil, resp.StatusCode, err
		}
		return nil, h.CatchAllErrStatus, err
	}

	if h.ReturnHead {
		for name, value := range resp.Header {
			tmp[name] = value
		}
	}
	if h.ReturnCert && resp.TLS != nil {
		certificates := resp.TLS.PeerCertificates
		if len(certificates) > 0 {
			// you probably want certificates[0]
			cert := certificates[0]
			tmp["cert"] = cert
		}
	}

	// Read body to buffer
	body, err := io.ReadAll(resp.Body)
	if bserr.Err(err, "Error reading body") {
		if resp != nil {
			return nil, resp.StatusCode, err
		}
		return nil, h.CatchAllErrStatus, err
	}

	tmp["body"] = string(body)
	return tmp, resp.StatusCode, nil
}

type HTTPClientSettings struct {
	ReturnHeaders           bool
	DisableVerifyClientCert bool
	RequestTimeout          int
	ResponseHeaderTimeout   int
	ReturnCert              bool
	ReUseClient             bool
	RedirectUseLastResponse bool
	DialTimeout             int
	DialKeepAliveTimeout    int
	MaxIdleConnections      int
	IdleConnectionTimeout   int
	TLSHandshakeTimeout     int
	CatchAllErrStatus       int
}

func NewClientSettings(returnHeaders, disableVerifyClientCert bool, requestTimeout, responseHeaderTimeout int) *HTTPClientSettings {
	return &HTTPClientSettings{
		ReturnHeaders:           returnHeaders,
		DisableVerifyClientCert: disableVerifyClientCert,
		RequestTimeout:          requestTimeout,
		ResponseHeaderTimeout:   responseHeaderTimeout,
		ReUseClient:             true,
		ReturnCert:              true,
		RedirectUseLastResponse: false,
	}
}

// MakeCall to http client
func MakeCall(method, URI string, msg io.Reader, header map[string]string, auth *Auth, settings *HTTPClientSettings) (map[string]interface{}, int, error) {
	c := NewHTTPClient(method, URI, header, auth, settings)
	mapData, status, err := c.Process(msg)
	if err != nil {
		return mapData, status, err
	}
	return mapData, status, nil
}

// HTTPStatusText lookup status code for text
func HTTPStatusText(status int) string {
	txt := http.StatusText(status)
	if len(txt) > 0 {
		txt = strings.ReplaceAll(txt, " ", "")
	} else if status == 1000 {
		txt = "ConnectRefused"
	} else {
		txt = strconv.Itoa(status)
	}
	return txt
}
