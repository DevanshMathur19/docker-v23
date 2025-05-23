package system // import "github.com/DevanshMathur19/docker-v23/pkg/system"

// Mknod is not implemented on Windows.
func Mknod(path string, mode uint32, dev int) error {
	return ErrNotSupportedPlatform
}

// Mkdev is not implemented on Windows.
func Mkdev(major int64, minor int64) uint32 {
	panic("Mkdev not implemented on Windows.")
}
