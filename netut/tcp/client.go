package tcp

import (
	"fmt"
	"net"
	"strings"

	log "github.com/colt3k/nglog/ng"

	"github.com/colt3k/utils/netut"
)

type NetworkClient struct {
}

func NewClient() {
}

/*
Start as a client to test ports on passed server:port entries and output a report
*/
func (c *NetworkClient) Client(contype string, hostAr []string, proxy string) *[]netut.Host {
	ip, err := netut.RetrieveIP()
	if err != nil {
		log.Logf(log.FATAL, "issue retrieving ip\n%+v", err)
	}

	hosts := make([]netut.Host, len(hostAr))

	for i, d := range hostAr {
		hostdata := strings.Split(d, ":")
		if len(hostdata) > 1 {
			hosts[i] = netut.Host{IP: hostdata[0], Port: hostdata[1]}
		}
		log.Logln(log.DBGL3, "")
		log.Logln(log.DBGL3, "*********************** TEST  ***********************")

		log.Logln(log.DBGL3, "Resovling: ", d)
		tcpAdr, err := net.ResolveTCPAddr(contype, d)
		log.Logln(log.DBGL3, "Resolved: ", d)
		if err != nil {
			log.Logf(log.ERROR, "failed to resolve address: %s on %s\n%+v", d, contype, err)
			hosts[i].Pass = false
			continue
		}

		con, err := net.DialTCP(contype, nil, tcpAdr)
		if err != nil {
			log.Logf(log.ERROR, "dial failed on: %s\n%+v", d, err)
			hosts[i].Pass = false
			continue
		}

		if con != nil {
			_, err = con.Write([]byte(fmt.Sprintf("Connection Test from:%s", ip)))
			if err != nil {
				log.Logf(log.ERROR, "write to server %s failed\n%+v", d, err)
				hosts[i].Pass = false
				if con != nil {
					err = con.Close()
					if err != nil {
						log.Logf(log.ERROR, "issue closing %+v", err)
					}
				}
				continue
			}
			if con != nil {
				err = con.Close()
				if err != nil {
					log.Logf(log.ERROR, "issue closing %+v", err)
				}
			}
		} else {
			log.Logln(log.ERROR, "TCP Write to server: ", d, " failed.")
			hosts[i].Pass = false

			continue
		}

		reply := make([]byte, 1024)

		repLen, err := con.Read(reply)
		if err != nil {
			log.Logf(log.ERROR, "TCP read from server failed\n%+v", err)
			hosts[i].Pass = false
			if con != nil {
				err = con.Close()
				if err != nil {
					log.Logf(log.ERROR, "issue closing %+v", err)
				}
			}
			continue
		}

		log.Logln(log.DBGL3, "TCP  Reply from Server: ", d, " Response: ", string(reply[:repLen]))
		hosts[i].Pass = true
	}

	return &hosts
}
