package test

import (
	"context"
	"encoding/json"
	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/netut/hc"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"
	"testing"
	"time"

	"github.com/colt3k/utils/netut"
)

func TestProxy(t *testing.T) {
	err := os.Setenv("http_proxy", "http://myproxy.domain.com")
	if err != nil {
		t.Errorf("issue setting env %+v", err)
	}
	err = os.Setenv("https_proxy", "http://myproxy.domain.com")
	if err != nil {
		t.Errorf("issue setting env %+v", err)
	}
}

func TestPing(t *testing.T) {
	avail, err := netut.Ping("192.168.1.1")
	if err != nil {
		t.Errorf( "issue no ping %+v", err)
	}
	if avail {
		t.Log("available")
	} else {
		t.Log("NOT available")
	}
}

func createTransport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second, // time spent establishing a TCP connection
			KeepAlive: 0,
		}).DialContext,
	}
}
func TestTrace(t *testing.T) {
	ca := log.NewConsoleAppender("*")
	log.Modify(log.LogLevel(log.DEBUG), log.ColorsOn(), log.Appenders(ca))

	endpoint := "https://google.com"
	trace, info := hc.Trace()

	tr := createTransport()
	c := &http.Client{
		Transport: tr,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 11*time.Second)
	defer cancel()

	// Prepare Request
	req, err := http.NewRequestWithContext(httptrace.WithClientTrace(ctx, trace), http.MethodGet, endpoint, nil)
	if err != nil {
		t.Fatalf("request error %v", err)
	}

	start := time.Now()
	info.GotFirstResponseByte.Start = start
	info.Start = start
	res, err := c.Do(req)
	if err != nil {
		t.Fatalf("client error %v", err)
	}
	defer res.Body.Close()

	// Read
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error %v", err)
	}
	// output details
	d, err := json.MarshalIndent(info, "", "    ")
	t.Logf("%v", string(d))
}