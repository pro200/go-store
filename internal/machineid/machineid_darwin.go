//go:build darwin

package machineid

import (
	"errors"
	"os/exec"
	"strings"
)

// platformID returns the IOPlatformUUID, which is stable for the
// lifetime of the machine.
func platformID() (string, error) {
	out, err := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice").Output()
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, "IOPlatformUUID") {
			continue
		}
		// "IOPlatformUUID" = "XXXXXXXX-XXXX-..."
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		return strings.Trim(strings.TrimSpace(parts[1]), `"`), nil
	}
	return "", errors.New("machineid: IOPlatformUUID not found")
}
