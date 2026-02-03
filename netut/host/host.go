package host

import (
	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/osut"
	"net"
	"os"
	"runtime"
	"strings"
)

// GetFQDN retrieve the fully qualified domain name
func GetFQDN() string {
	hostname, _ := os.Hostname()
	// ensure it has an fqdn with parts separated by periods
	if strings.Contains(hostname, ".") {
		return hostname
	}
	// if fqdn doesn't have parts use hostname -f command
	if runtime.GOOS != "windows" {
		var hostname2 string
		var err error
		if runtime.GOOS != "darwin" {
			hostname2, err = osut.CallCmd("hostname --fqdn")
			if err != nil {
				log.Logf(log.ERROR, "-- 'hostname --fqdn' %v", err)
			}
		} else if runtime.GOOS == "darwin" {
			hostname2, err = osut.CallCmd("hostname -f")
			if err != nil {
				log.Logf(log.ERROR, "-- 'hostname -f' %v", err)
			}
		}
		hostname2 = strings.TrimSpace(hostname2)
		log.Logf(log.DEBUG, "-- via 'hostname --fqdn' %v", hostname2)
		// is this an FQDN ???
		if strings.Contains(hostname2, ".") {
			return hostname2
		}
		hostname = hostname2
	}

	addrs, err := net.LookupIP(hostname)
	if err != nil {
		return hostname
	}
	log.Logf(log.DEBUG, "-- LookupIP addresses %v", addrs)

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			var ip []byte
			ip, err = ipv4.MarshalText()
			if err != nil {
				return hostname
			}
			log.Logf(log.DEBUG, "-- ipv4 marshal text %v", string(ip))
			var hosts []string
			hosts, err = net.LookupAddr(string(ip))
			if err != nil || len(hosts) == 0 {
				return hostname
			}
			log.Logf(log.DEBUG, "-- hosts %v", hosts)
			fqdn := hosts[0]
			t := strings.TrimSuffix(fqdn, ".") // return fqdn without trailing dot
			return t
		}
	}
	return hostname
}

// GetFQDNLower retrieve the fully qualified domain name in lowercase
func GetFQDNLower() string {
	fqdn, _ := os.Hostname()
	return strings.ToLower(fqdn)
}

// GetHost retrieve the first part of the FQDN
func GetHost() string {
	pts := strings.Split(GetFQDN(), ".")
	output := pts[0]
	return output
}

// GetHostLower retrieve the first part of the FQDN in lowercase
func GetHostLower() string {
	h := GetHost()
	return strings.ToLower(h)
}
