package netut

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"

	log "github.com/colt3k/nglog/ng"
)

//RetrieveIP retrieve system ipv4 address
func RetrieveIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network")
}

// Parse and set Domain/IP's passed in
func ParseIPs(ipAddresses []string) ([]string, []net.IP, error) {
	var parsed []net.IP
	var domains []string
	for _, s := range ipAddresses {
		domainIp := strings.Split(s, "/")
		for _, j := range domainIp {
			if ip := net.ParseIP(j); ip != nil {
				parsed = append(parsed, ip)
			} else {
				domains = append(domains, j)
			}
		}
	}
	return domains, parsed, nil
}
func GetLocalIP(doesNotStartWith []string) string {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Logf(log.FATAL, "issue retrieving network interfaces\n%+v", err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Logf(log.FATAL, "issue getting address\n%+v", err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// process IP address

			// Compare against all prefixes
			var found bool
			ipStr := ip.String()
			for _, d := range doesNotStartWith {
				if strings.HasPrefix(ipStr, d) {
					found = true
				}
			}
			if !found && strings.Index(ipStr, ":") == -1 {
				return ip.String()
			}
		}
	}
	return ""
}

func Ping(server string) (bool, error) {

	var valid bool
	var err error
	log.Logln(log.DEBUG, "******************** PING - resolution ****************************")
	ip, err := net.ResolveIPAddr("ip4:icmp", server)
	if err != nil {
		log.Logln(log.DEBUG, "******************** END PING - resolution *************************")
		return valid, err
	}
	log.Logf(log.DEBUG, "%s resolved to IP %s", server, ip.IP.String())
	addr := net.ParseIP(ip.IP.String())
	if addr == nil {
		valid = false
		err = fmt.Errorf("invalid ip address")
	}
	if isIPV4(addr) || isIPv6(addr) {
		valid = call(addr.String())
		err = nil
	} else {
		valid = false
		err = fmt.Errorf("invalid format for ipv4 or ipv6")
	}
	log.Logln(log.DEBUG, "******************** END PING - resolution *************************\n")
	return valid, err
}

func buildPingCmd(addr string) []string {
	if runtime.GOOS == "darwin" {
		return []string{addr, "-c 1", "-t 3"}
	} else {
		return []string{addr, "-c 1", "-W 3"}
	}
}

func call(addr string) bool {
	var available bool
	parms := buildPingCmd(addr)
	log.Logln(log.DEBUG, "******************** START PING ****************************")
	log.Logln(log.DEBUG, "executing ping ", parms)
	out, _ := exec.Command("ping", parms...).Output()
	log.Logf(log.DEBUG, "ping output |%s|", string(out))
	if strings.Contains(string(out), "Destination Host Unreachable") {
		log.Logln(log.DEBUG, "destination host unreachable")
		available = false
	} else if strings.Contains(string(out), "100.0% packet loss") {
		log.Logln(log.DEBUG, "100.0% packet lost")
		available = false
	} else {
		available = true
	}
	log.Logln(log.DEBUG, "******************** END PING ****************************")
	return available
}

func isIPV4(ip net.IP) bool {
	return len(ip.To4()) == net.IPv4len
}
func isIPv6(ip net.IP) bool {
	return len(ip) == net.IPv6len
}
