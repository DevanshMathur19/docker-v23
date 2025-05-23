package container // import "github.com/DevanshMathur19/docker-v23/integration/container"

import (
	"context"
	"strings"
	"testing"

	"github.com/DevanshMathur19/docker-v23/api/types"
	"github.com/DevanshMathur19/docker-v23/api/types/versions"
	"github.com/DevanshMathur19/docker-v23/integration/internal/container"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/skip"
)

func TestExecConsoleSize(t *testing.T) {
	skip.If(t, testEnv.DaemonInfo.OSType != "linux")
	skip.If(t, versions.LessThan(testEnv.DaemonAPIVersion(), "1.42"), "skip test from new feature")

	defer setupTest(t)()
	client := testEnv.APIClient()
	ctx := context.Background()

	cID := container.Run(ctx, t, client, container.WithImage("busybox"))

	result, err := container.Exec(ctx, client, cID, []string{"stty", "size"},
		func(ec *types.ExecConfig) {
			ec.Tty = true
			ec.ConsoleSize = &[2]uint{57, 123}
		},
	)

	assert.NilError(t, err)
	assert.Equal(t, strings.TrimSpace(result.Stdout()), "57 123")
}
