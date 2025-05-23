//go:build !windows
// +build !windows

package opts // import "github.com/DevanshMathur19/docker-v23/opts"

const (
	// DefaultHTTPHost Default HTTP Host used if only port is provided to -H flag e.g. dockerd -H tcp://:8080
	DefaultHTTPHost = "localhost"

	// DefaultHost constant defines the default host string used by docker on other hosts than Windows
	DefaultHost = "unix://" + DefaultUnixSocket
)
