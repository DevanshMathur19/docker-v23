package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

import (
	"strings"

	"github.com/DevanshMathur19/docker-v23/container"
)

// excludeByIsolation is a platform specific helper function to support PS
// filtering by Isolation. This is a Windows-only concept, so is a no-op on Unix.
func excludeByIsolation(container *container.Snapshot, ctx *listContext) iterationAction {
	i := strings.ToLower(string(container.HostConfig.Isolation))
	if i == "" {
		i = "default"
	}
	if !ctx.filters.Match("isolation", i) {
		return excludeContainer
	}
	return includeContainer
}
