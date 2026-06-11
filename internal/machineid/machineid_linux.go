//go:build linux

package machineid

import (
	"errors"
	"os"
	"strings"
)

// platformID reads the systemd/dbus machine id.
func platformID() (string, error) {
	for _, path := range []string{"/etc/machine-id", "/var/lib/dbus/machine-id"} {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if id := strings.TrimSpace(string(data)); id != "" {
			return id, nil
		}
	}
	return "", errors.New("machineid: machine-id file not found")
}
