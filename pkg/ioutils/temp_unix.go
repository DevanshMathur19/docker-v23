//go:build !windows
// +build !windows

package ioutils // import "github.com/DevanshMathur19/docker-v23/pkg/ioutils"

import "os"

// TempDir on Unix systems is equivalent to os.MkdirTemp.
func TempDir(dir, prefix string) (string, error) {
	return os.MkdirTemp(dir, prefix)
}
