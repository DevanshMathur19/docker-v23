//go:build !windows
// +build !windows

package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

import (
	"os"
	"os/signal"

	"github.com/DevanshMathur19/docker-v23/pkg/stack"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func (daemon *Daemon) setupDumpStackTrap(root string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, unix.SIGUSR1)
	go func() {
		for range c {
			path, err := stack.DumpToFile(root)
			if err != nil {
				logrus.WithError(err).Error("failed to write goroutines dump")
			} else {
				logrus.Infof("goroutine stacks written to %s", path)
			}
		}
	}()
}
