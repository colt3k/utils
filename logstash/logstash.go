package logstash

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"time"

	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/file"
)

const (
	// TCPPROTOCOL variable constant for tcp
	TCPPROTOCOL = "tcp"
)

// Server struct to hold building of a server
type Server struct {
	Host          string
	Port          int
	Connection    *net.TCPConn
	TLSConnection *tls.Conn
	Timeout       int
}

// New host name, port and timeout in seconds
func New(host string, port, timeout int) *Server {
	t := new(Server)
	t.Host = host
	t.Port = port
	t.Timeout = timeout
	return t
}

// Show host port and timeout
func (s *Server) Show() {
	log.Logln(log.DEBUG, "Host:", s.Host)
	log.Logln(log.DEBUG, "Port:", s.Port)
	log.Logln(log.DEBUG, "Timeout:", s.Timeout)
}

// Connect create server connection
func (s *Server) Connect() (*net.TCPConn, error) {

	host := fmt.Sprintf("%s:%d", s.Host, s.Port)
	addr, err := net.ResolveTCPAddr(TCPPROTOCOL, host)
	if err != nil {
		return nil, err
	}
	con, err := net.DialTCP(TCPPROTOCOL, nil, addr)
	if err != nil {
		return nil, err
	}
	if con != nil {
		s.Connection = con
		s.Connection.SetLinger(0)
		s.Connection.SetNoDelay(true)
		s.Connection.SetKeepAlive(true)
		s.Connection.SetKeepAlivePeriod(5 * time.Second)
		s.UpdateTimeout()
	}
	return s.Connection, nil
}

// ConnectTLS connect using TLS
func (s *Server) ConnectTLS() (*tls.Conn, error) {
	cp := x509.NewCertPool()
	home := file.HomeFolder()
	data, err := ioutil.ReadFile(filepath.Join(home, "mycerts", "lsCA.pem"))
	if err != nil {
		return nil, err
	}
	cp.AppendCertsFromPEM(data)

	var cfg tls.Config
	cfg.RootCAs = cp
	cfg.ServerName = s.Host
	cfg.BuildNameToCertificate()
	cfg.ClientAuth = tls.RequireAndVerifyClientCert
	//cfg.CipherSuites = []uint16 {
	//	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	//}
	//cfg.InsecureSkipVerify = true		// only needed if you don't have the CA that created the servers key/cert

	host := fmt.Sprintf("%s:%d", s.Host, s.Port)
	log.Logf(log.DEBUG, "connecting to %s", host)
	con, err := tls.Dial(TCPPROTOCOL, host, &cfg)
	if err != nil {
		if con != nil {
			con.Close()
		}
		return nil, err
	}
	s.TLSConnection = con
	return s.TLSConnection, nil
}

// Write out data on connection
func (s *Server) Write(p []byte) (n int, err error) {
	var i int
	if s.Connection != nil {
		msg := fmt.Sprintf("%s\n", string(p))
		i, err = s.Connection.Write([]byte(msg))
		if err != nil {
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				s.Connection.Close()
				s.Connection = nil
				if err != nil {
					return i, err
				}
			} else {
				s.Connection.Close()
				s.Connection = nil
				return i, err
			}
		}
		s.UpdateTimeout()
		return i, nil
	}
	return i, fmt.Errorf("tcp connection nil")
}

// WriteTLS write via TLS connection
func (s *Server) WriteTLS(p []byte) (n int, err error) {
	var i int
	if s.TLSConnection != nil {
		msg := fmt.Sprintf("%s\n", string(p))
		i, err = s.TLSConnection.Write([]byte(msg))
		if err != nil {
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				s.TLSConnection.Close()
				s.TLSConnection = nil
				if err != nil {
					return i, err
				}
			} else {
				s.TLSConnection.Close()
				s.TLSConnection = nil
				return i, err
			}
		}
		s.UpdateTLSTimeout()
		return i, nil
	}
	return i, fmt.Errorf("tcp connection nil")
}

// UpdateTimeout increase time to continue work
func (s *Server) UpdateTimeout() {
	end := time.Now().Add(time.Duration(s.Timeout) * time.Second)
	s.Connection.SetDeadline(end)
	s.Connection.SetWriteDeadline(end)
	s.Connection.SetReadDeadline(end)
}

// UpdateTLSTimeout increase time to continue work
func (s *Server) UpdateTLSTimeout() {
	end := time.Now().Add(time.Duration(s.Timeout) * time.Second)
	s.TLSConnection.SetDeadline(end)
	s.TLSConnection.SetWriteDeadline(end)
	s.TLSConnection.SetReadDeadline(end)
}

// Close close server connection
func (s *Server) Close() {
	err := s.Connection.Close()
	if err != nil {
		log.Logln(log.ERROR, "error closing connection ", err)
	}
}

// CloseTLS close TLS server connection
func (s *Server) CloseTLS() {
	err := s.TLSConnection.Close()
	if err != nil {
		log.Logln(log.ERROR, "error closing connection ", err)
	}
}
