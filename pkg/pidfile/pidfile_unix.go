//go:build !windows && !darwin
// +build !windows,!darwin

package pidfile // import "github.com/DevanshMathur19/docker-v23/pkg/pidfile"

import (
	"os"
	"path/filepath"
	"strconv"
)

func processExists(pid int) bool {
	if _, err := os.Stat(filepath.Join("/proc", strconv.Itoa(pid))); err == nil {
		return true
	}
	return false
}
