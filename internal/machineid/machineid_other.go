//go:build !darwin && !linux && !windows

package machineid

import "errors"

func platformID() (string, error) {
	return "", errors.New("machineid: unsupported platform")
}
