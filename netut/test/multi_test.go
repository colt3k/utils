package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/netut/hc"
	"github.com/colt3k/utils/netut/https"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"

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
		t.Errorf("issue no ping %+v", err)
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

func FileExistsAndIsADir(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
func certPathDir() string {
	homeDir, _ := os.UserHomeDir()
	certPth := filepath.Join(homeDir, "dev", "dev-setup", "certs")
	if !FileExistsAndIsADir(certPth) {
		log.Logf(log.ERROR, "Certs NOT Found: %v", certPth)
	}
	return certPth
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

func TestLoadCert(t *testing.T) {
	IP := ""
	testMode := false
	sslDisable := false
	certPath := certPathDir()
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("SSL Intermediate Chain Working!"))
	})
	var port int64
	port = 8088
	log.Logf(log.INFO, "Listening on ...%s:%s", IP, strconv.FormatInt(port, 10))
	ctx := context.Background()
	server := https.NewWithContext(ctx, muxRouter, fmt.Sprintf("%s:%s", IP, strconv.FormatInt(port, 10)), false)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {

		if testMode && !sslDisable {
			log.Logf(log.INFO, "-- running in SSL mode testMode")
			return server.ListenAndServeTLS(filepath.Join(certPath, "fullchain-local.pem"),
				filepath.Join(certPath, "privkey-local.pem"))
		} else if !sslDisable {
			log.Logf(log.INFO, "-- running in SSL mode NON testMode\n\tcert: %v\n\tkey %v",
				filepath.Join(certPath, "fullchain.pem"),
				filepath.Join(certPath, "privkey.pem"))
			return server.ListenAndServeTLS(filepath.Join(certPath, "fullchain.pem"),
				filepath.Join(certPath, "privkey.pem"))
		}
		log.Logf(log.INFO, "-- running in NON SSL mode")
		return server.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return server.Shutdown(context.Background())
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}
