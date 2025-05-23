//go:build !windows
// +build !windows

package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

import (
	"github.com/DevanshMathur19/docker-v23/container"
	volumemounts "github.com/DevanshMathur19/docker-v23/volume/mounts"
)

// checkIfPathIsInAVolume checks if the path is in a volume. If it is, it
// cannot be in a read-only volume. If it  is not in a volume, the container
// cannot be configured with a read-only rootfs.
func checkIfPathIsInAVolume(container *container.Container, absPath string) (bool, error) {
	var toVolume bool
	parser := volumemounts.NewParser()
	for _, mnt := range container.MountPoints {
		if toVolume = parser.HasResource(mnt, absPath); toVolume {
			if mnt.RW {
				break
			}
			return false, ErrVolumeReadonly
		}
	}
	return toVolume, nil
}

// isOnlineFSOperationPermitted returns an error if an online filesystem operation
// is not permitted.
func (daemon *Daemon) isOnlineFSOperationPermitted(container *container.Container) error {
	return nil
}
