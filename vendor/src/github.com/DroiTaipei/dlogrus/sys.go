package dlogrus

/*
URL: https://github.com/mccoyst/myip/blob/master/myip.go
URL: http://changsijay.com/2013/07/28/golang-get-ip-address/
*/

import (
	"errors"
	"net"
	"os"
)

func GetPodname() (string, error) {
	return os.Hostname()
}

func GetIP() (ret string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ret = ipnet.IP.String()
				return
			}
		}
	}
	err = errors.New("Got IP Failed")
	return
}
