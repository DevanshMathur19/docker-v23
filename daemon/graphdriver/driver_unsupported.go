//go:build !linux && !windows && !freebsd
// +build !linux,!windows,!freebsd

package graphdriver // import "github.com/DevanshMathur19/docker-v23/daemon/graphdriver"

var (
	// List of drivers that should be used in an order
	priority = "unsupported"
)

// GetFSMagic returns the filesystem id given the path.
func GetFSMagic(rootpath string) (FsMagic, error) {
	return FsMagicUnsupported, nil
}
