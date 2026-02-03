package host

import (
	"bytes"
	log "github.com/colt3k/nglog/ng"
	"net"
)

// GetIPs looks at interfaces and returns list of comma separated ipv4 IPs
func GetIPs() string {
	var ips bytes.Buffer
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Logf(log.ERROR, "%v", err)
	}
	// handle err
	for _, i := range ifaces {
		var addrs []net.Addr
		addrs, err = i.Addrs()
		if err != nil {
			log.Logf(log.ERROR, "%v", err)
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				if !v.IP.IsLoopback() {
					if v.IP.To4() != nil {
						ip = v.IP
					}
				}
			case *net.IPAddr:
				if !v.IP.IsLoopback() {
					if v.IP.To4() != nil {
						ip = v.IP
					}
				}
			}
			if ip != nil {
				if ips.Len() > 0 {
					ips.WriteString(",")
				}
				ips.WriteString(ip.String())
			}
		}
	}

	return ips.String()
}

// GetOutboundIP find IP address that sends publicly
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Logf(log.ERROR, "%v", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// FindIP find host IP address
func FindIP(host string) string {
	ips := ""
	addr, err := net.LookupIP(host)
	if err != nil {
		ips = "unknown"
	} else {
		for l, m := range addr {
			if l > 0 {
				ips += ", "
			}
			ips += m.String()
		}
	}
	return ips
}
