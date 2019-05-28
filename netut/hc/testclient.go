package hc

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/colt3k/nglog/ng"

	"github.com/colt3k/utils/netut"
)

type TestClient struct {
}

//Client create an HTTP Client
func (c *TestClient) Client(contype string, hostAr []string, proxy string) *[]netut.Host {

	hosts := make([]netut.Host, len(hostAr))

	for i, d := range hostAr {

		hosts[i] = netut.Host{URL: d}

		hosts[i].IP = *getHost(d)
		hosts[i].Port = *getPort(d)

		log.Logln(log.DEBUG, "")
		log.Logln(log.DEBUG, "*********************** TEST  ***********************")

		tr := &http.Transport{}
		if len(strings.TrimSpace(proxy)) > 0 {
			proxyURL, _ := url.Parse(proxy)

			tr = &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    15 * time.Second,
				DisableCompression: true,
				Proxy:              http.ProxyURL(proxyURL),
			}

		} else {

			tr = &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    15 * time.Second,
				DisableCompression: true,
			}

		}

		client := &http.Client{Transport: tr}

		resp, err := client.Get(hosts[i].URL)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			log.Logln(log.DEBUG, "Failed to resolve address: ", d, " on ", contype)
			hosts[i].Pass = false
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)

		log.Logln(log.DEBUG, "Reply from Server: ", d, " Response: ", string(body))
		hosts[i].Pass = true
	}
	return &hosts
}

func getHost(url string) *string {
	var host string
	pfx := 2
	idx := strings.Index(url, "//")
	if idx == -1 {
		pfx = 0
		idx = 0
	}
	lastIdx := strings.LastIndex(url, ":")
	runes := []rune(url)
	if lastIdx < idx {
		//Go by next slash instead, there is no next :
		tmp := string(runes[idx+pfx : len(url)])
		nextSlash := strings.Index(tmp, "/")
		runes2 := []rune(tmp)
		host = string(runes2[:nextSlash])
	} else {
		host = string(runes[idx+pfx : lastIdx])
	}
	return &host
}

func getPort(url string) *string {
	lastIdx := strings.LastIndex(url, ":")
	//Not the colon in beginning, get port
	var port string
	if lastIdx > 6 {
		runes := []rune(url)
		postPart1 := string(runes[lastIdx:len(url)])
		nxtSlash := strings.Index(postPart1, "/")
		if nxtSlash == -1 {
			nxtSlash = len(url)
		}
		port = string(runes[lastIdx+1 : nxtSlash])
	} else {
		runes := []rune(url)
		pfx := string(runes[:lastIdx])
		if pfx == "http" {
			port = "80"
		} else if pfx == "https" {
			port = "443"
		}
	}
	return &port
}
