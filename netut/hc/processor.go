package hc

import (
	"github.com/colt3k/nglog/ers/bserr"
	log "github.com/colt3k/nglog/ng"
	"io"
	"net/http"
	"strconv"
	"strings"
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
	DisableVerifyClientCert bool
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
	} else {
		t.ReturnHead = true
		t.DisableVerifyClientCert = true
		t.RequestTimeout = 120
		t.ResponseHeaderTimeout = 120
		t.ReturnCert = true
		t.ReUseClient = true
		t.RedirectUseLastResponse = false
	}
	return t
}

// Process method on HTTPCall object, pass in reader
func (h *HTTPClient) Process(data io.Reader) (map[string]interface{}, int, error) {
	log.Logf(log.DBGL2, "-- called httpCallData for Method %s URL: %s", h.Method, h.URL)
	tmp := make(map[string]interface{}, 0)

	if client == nil || !h.ReUseClient {
		//log.Logln(log.DEBUG, "!!! Creating NEW HTTP CLIENT !!!")
		// set to timeout after a day per request, accommodates file uploads
		client = NewClient(HttpClientRequestTimeout(h.RequestTimeout), DisableVerifyClientCert(h.DisableVerifyClientCert),
			HttpClientResponseHeaderTimeout(h.ResponseHeaderTimeout), CheckRedirectUserLastResp(h.RedirectUseLastResponse))
	}
	var resp, err = client.Fetch(h.Method, h.URL, h.Auth, h.Header, data)

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
		return nil, 1000, err
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
		return nil, 1000, err
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
