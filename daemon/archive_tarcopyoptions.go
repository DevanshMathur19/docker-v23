package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

import (
	"github.com/DevanshMathur19/docker-v23/pkg/archive"
)

// defaultTarCopyOptions is the setting that is used when unpacking an archive
// for a copy API event.
func (daemon *Daemon) defaultTarCopyOptions(noOverwriteDirNonDir bool) *archive.TarOptions {
	return &archive.TarOptions{
		NoOverwriteDirNonDir: noOverwriteDirNonDir,
		IDMap:                daemon.idMapping,
	}
}
