//go:build !linux && !freebsd
// +build !linux,!freebsd

package logger // import "github.com/DevanshMathur19/docker-v23/daemon/logger"

import (
	"errors"
	"io"
)

func openPluginStream(a *pluginAdapter) (io.WriteCloser, error) {
	return nil, errors.New("log plugin not supported")
}
