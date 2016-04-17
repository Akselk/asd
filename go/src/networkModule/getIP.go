package networkModule

import (
	"net"
	"strings"
)

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetBroadcastIP() string {
	fields := strings.SplitAfterN(GetLocalIP(), ".", 4)
	return strings.Join(fields[0:3], "") + "255"
}
