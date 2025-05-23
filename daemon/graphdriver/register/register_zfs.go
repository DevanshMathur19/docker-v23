//go:build (!exclude_graphdriver_zfs && linux) || (!exclude_graphdriver_zfs && freebsd)
// +build !exclude_graphdriver_zfs,linux !exclude_graphdriver_zfs,freebsd

package register // import "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/register"

import (
	// register the zfs driver
	_ "github.com/DevanshMathur19/docker-v23/daemon/graphdriver/zfs"
)
