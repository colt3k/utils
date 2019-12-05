package logstash

import (
	"encoding/json"
	"testing"

	log "github.com/colt3k/nglog/ng"
)

type Msg struct {
	Hostname    string
	Application string
	Text        string
}

func TestNonTLS(t *testing.T) {
	server := New("127.0.0.1", 5000, 60)
	_, err := server.Connect()
	if err != nil {
		log.Logf(log.FATAL, "issue opening tcp server %+v", err)
	}
	m := &Msg{Hostname:"server1", Application:"logstash", Text:"non TLS message"}
	b, err := json.Marshal(m)
	if err != nil {
		log.Logf(log.FATAL, "issue marshalling message %+v", err)
	}
	_, err = server.Write(b)
	server.Close()
}

func TestTLS(t *testing.T) {

	server := New("127.0.0.1", 5000, 60)
	_, err := server.ConnectTLS()
	if err != nil {
		log.Logf(log.FATAL, "issue opening tcp server %+v", err)
	}
	m := &Msg{Hostname:"server2", Application:"logstash", Text:"TLS message"}
	b, err := json.Marshal(m)
	if err != nil {
		log.Logf(log.FATAL, "issue marshalling message %+v", err)
	}
	_, err = server.WriteTLS(b)
	server.CloseTLS()
}
