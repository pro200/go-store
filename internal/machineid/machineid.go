// Package machineid provides a stable, machine-bound identifier
// without requiring network connectivity.
package machineid

import (
	"errors"
	"net"
	"os"
)

// ID returns a stable identifier for the current machine.
// It tries, in order: the OS machine id, the MAC address of the first
// up non-loopback interface, and the hostname. It returns an error
// only if all sources fail.
func ID() (string, error) {
	if id, err := platformID(); err == nil && id != "" {
		return id, nil
	}
	if mac, err := firstMAC(); err == nil && mac != "" {
		return mac, nil
	}
	if host, err := os.Hostname(); err == nil && host != "" {
		return host, nil
	}
	return "", errors.New("machineid: no identifier available")
}

// firstMAC returns the hardware address of the first up,
// non-loopback interface that has one.
func firstMAC() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if len(iface.HardwareAddr) > 0 {
			return iface.HardwareAddr.String(), nil
		}
	}
	return "", errors.New("machineid: no usable network interface")
}
