package system // import "github.com/DevanshMathur19/docker-v23/pkg/system"

import (
	"errors"
)

var (
	// ErrNotSupportedPlatform means the platform is not supported.
	ErrNotSupportedPlatform = errors.New("platform and architecture is not supported")

	// ErrNotSupportedOperatingSystem means the operating system is not supported.
	ErrNotSupportedOperatingSystem = errors.New("operating system is not supported")
)
