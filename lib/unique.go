package lib

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
)

func outboundInterface() (net.IP, error) {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

func findInterfaceByIP(target net.IP) (*net.Interface, error) {
	ifaces, _ := net.Interfaces()

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && ipNet.IP.Equal(target) {
				return &iface, nil
			}
		}
	}
	return nil, fmt.Errorf("not found interfaces")
}

// 컴퓨터를 구분할수 있는 ID로 맥어드레스를 md5로 암호화한 문자열
func CUID() (string, error) {
	ip, err := outboundInterface()
	if err != nil {
		return "", errors.New("not connected internet")
	}

	iface, err := findInterfaceByIP(ip)
	if err != nil {
		return "", errors.New("not found interfaces")
	}

	hAddr := iface.HardwareAddr.String()
	uid := md5.Sum([]byte(hAddr))
	return hex.EncodeToString(uid[:]), nil
}
