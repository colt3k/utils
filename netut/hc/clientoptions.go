package hc

import "time"

type ClientOption func(h *Client)

func DialTimeout(i int) ClientOption {
	return func(c *Client) {
		c.DialTimeout = time.Duration(i) * time.Second
	}
}

func DialKeepAliveTimeout(i int) ClientOption {
	return func(c *Client) {
		c.DialKeepAliveTimeout = time.Duration(i) * time.Second
	}
}

func MaxIdleConnections(i int) ClientOption {
	return func(c *Client) {
		c.MaxIdleConnections = i
	}
}

func IdleConnectionTimeout(i int) ClientOption {
	return func(c *Client) {
		c.IdleConnTimeout = time.Duration(i) * time.Second
	}
}

func TLSHandshakeTimeout(i int) ClientOption {
	return func(c *Client) {
		c.TlsHandshakeTimeout = time.Duration(i) * time.Second
	}
}

//func ExpectContinueTimeout(i int) ClientOption {
//	return func(c *Client) {
//		c.ExpectContinueTimeout = time.Duration(i) * time.Second
//	}
//}

func HttpClientRequestTimeout(i int) ClientOption {
	return func(c *Client) {
		c.HttpClientRequestTimeout = time.Duration(i) * time.Second
	}
}
func DisableVerifyClientCert(b bool) ClientOption {
	return func(c *Client) {
		c.disableVerifyCert = b
	}
}
func HttpClientResponseHeaderTimeout(i int) ClientOption {
	return func(c *Client) {
		c.ResponseHeaderTimeout = time.Duration(i) * time.Second
	}
}