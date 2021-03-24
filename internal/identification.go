package internal

import (
	"fmt"
	"net"
	"os"
)

func IPs() ([]string, error) {
	var ips []string
	interfaces, err := net.Interfaces()
	if err != nil {
		return ips, fmt.Errorf("error while getting interfaces: %w", err)
	}
	for _, i := range interfaces {
		addresses, err := i.Addrs()
		if err != nil {
			return ips, fmt.Errorf("error while getting adresses for interface %v: %w", i.Name, err)
		}
		for _, address := range addresses {
			var ip net.IP
			switch v := address.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if !ip.IsLoopback() {
				ips = append(ips, ip.String())
			}
		}
	}
	return ips, nil
}

func Hostname() (string, error) {
	return os.Hostname()
}
