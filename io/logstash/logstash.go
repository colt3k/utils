package logstash

import (
	"fmt"
	"net"
	"time"

	log "github.com/colt3k/nglog/ng"
)

type Writer struct {
	Host       string
	Port       int
	Timeout    int
	Connection *net.TCPConn
}

func New(host string, port, timeout int) *Writer {
	t := new(Writer)
	t.Host = host
	t.Port = port
	t.Timeout = timeout
	return t
}

func (s *Writer) Show() {
	log.Println("Host:", s.Host)
	log.Println("Port:", s.Port)
	log.Println("Timeout:", s.Timeout)
}
func (s *Writer) Connect() (*net.TCPConn, error) {

	host := fmt.Sprintf("%s:%d", s.Host, s.Port)
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return nil, err
	}
	con, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	if con != nil {
		s.Connection = con
		err := s.Connection.SetLinger(0)
		if err != nil {
			log.Logf(log.ERROR, "issue set linger %+v", err)
		}
		err = s.Connection.SetNoDelay(true)
		if err != nil {
			log.Logf(log.ERROR, "issue set no delay %+v", err)
		}
		err = s.Connection.SetKeepAlive(true)
		if err != nil {
			log.Logf(log.ERROR, "issue set keep alive %+v", err)
		}
		err = s.Connection.SetKeepAlivePeriod(5 * time.Second)
		if err != nil {
			log.Logf(log.ERROR, "issue set keep alive period %+v", err)
		}
		s.UpdateTimeout()
	}
	return s.Connection, nil
}

func (s *Writer) Write(p []byte) (n int, err error) {
	var i int
	if s.Connection != nil {
		msg := fmt.Sprintf("%s\n", string(p))
		i, err = s.Connection.Write([]byte(msg))
		if err != nil {
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				err = s.Connection.Close()
				if err != nil {
					log.Logf(log.ERROR, "issue closing connection %+v", err)
				}
				s.Connection = nil
				if err != nil {
					return i, err
				}
			} else {
				err = s.Connection.Close()
				if err != nil {
					log.Logf(log.ERROR, "issue closing connection %+v", err)
				}
				s.Connection = nil
				return i, err
			}
		}
		s.UpdateTimeout()
		return i, nil
	}
	return i, fmt.Errorf("tcp connection nil")
}

func (s *Writer) UpdateTimeout() {
	end := time.Now().Add(time.Duration(s.Timeout) * time.Second)
	err := s.Connection.SetDeadline(end)
	if err != nil {
		log.Logf(log.ERROR, "issue set deadline %+v", err)
	}
	err = s.Connection.SetWriteDeadline(end)
	if err != nil {
		log.Logf(log.ERROR, "issue set write deadline %+v", err)
	}
	err = s.Connection.SetReadDeadline(end)
	if err != nil {
		log.Logf(log.ERROR, "issue set read deadline %+v", err)
	}
}
