package logstash

import (
	"errors"
	"fmt"
	"net"
	"time"

	"log"
)

type Writer struct {
	Host       string
	Port       int
	Timeout    int
	Connection *net.TCPConn
}

// New create an instance of Writer
func New(host string, port, timeout int) *Writer {
	t := new(Writer)
	t.Host = host
	t.Port = port
	t.Timeout = timeout
	return t
}

// Show displays settings
func (s *Writer) Show() {
	log.Println("Host:", s.Host)
	log.Println("Port:", s.Port)
	log.Println("Timeout:", s.Timeout)
}

// Connect connects to host
func (s *Writer) Connect() (*net.TCPConn, error) {

	host := fmt.Sprintf("%s:%d", s.Host, s.Port)
	addr, errResolve := net.ResolveTCPAddr("tcp", host)
	if errResolve != nil {
		return nil, errResolve
	}
	con, errDial := net.DialTCP("tcp", nil, addr)
	if errDial != nil {
		return nil, errDial
	}
	if con != nil {
		s.Connection = con
		errLinger := s.Connection.SetLinger(0)
		if errLinger != nil {
			log.Printf("ERROR: issue set linger %+v\n", errLinger)
		}
		errNoDelay := s.Connection.SetNoDelay(true)
		if errNoDelay != nil {
			log.Printf("ERROR: issue set no delay %+v\n", errNoDelay)
		}
		errKeepAlive := s.Connection.SetKeepAlive(true)
		if errKeepAlive != nil {
			log.Printf("ERROR: issue set keep alive %+v\n", errKeepAlive)
		}
		errKeepAlivePeriod := s.Connection.SetKeepAlivePeriod(5 * time.Second)
		if errKeepAlivePeriod != nil {
			log.Printf("ERROR: issue set keep alive period %+v\n", errKeepAlivePeriod)
		}
		s.UpdateTimeout()
	}
	return s.Connection, nil
}

// Write sends data over connection
func (s *Writer) Write(p []byte) (n int, err error) {
	var i int
	if s.Connection != nil {
		msg := fmt.Sprintf("%s\n", string(p))
		i, err = s.Connection.Write([]byte(msg))
		if err != nil {
			var netErr net.Error
			switch {
			case errors.As(err, netErr):
				if netErr.Timeout() {
					err = s.Connection.Close()
					if err != nil {
						log.Printf("ERROR: issue closing connection %+v\n", err)
					}
					s.Connection = nil
					if err != nil {
						return i, err
					}
				}
			default:
				err = s.Connection.Close()
				if err != nil {
					log.Printf("ERROR: issue closing connection %+v\n", err)
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
		log.Printf("ERROR: issue set deadline %+v\n", err)
	}
	err = s.Connection.SetWriteDeadline(end)
	if err != nil {
		log.Printf("ERROR: issue set write deadline %+v\n", err)
	}
	err = s.Connection.SetReadDeadline(end)
	if err != nil {
		log.Printf("ERROR: issue set read deadline %+v\n", err)
	}
}
