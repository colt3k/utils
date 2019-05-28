package main

import (
	"fmt"
	"os"

	log "github.com/colt3k/nglog/ng"

	"github.com/colt3k/utils/netut"
)

func main() {
	avail, err := netut.Ping("192.168.2.110")
	if err != nil {
		log.Logf(log.ERROR, "issue no ping %+v", err)
	}
	if avail {
		fmt.Println("available")
	} else {
		fmt.Println("NOT available")
	}
}

func proxy() {
	err := os.Setenv("http_proxy", "http://myproxy.domain.com")
	if err != nil {
		log.Logf(log.ERROR, "issue setting env %+v", err)
	}
	err = os.Setenv("https_proxy", "http://myproxy.domain.com")
	if err != nil {
		log.Logf(log.ERROR, "issue setting env %+v", err)
	}
}
