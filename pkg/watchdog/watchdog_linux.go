package watchdog

import (
	"os"

	"github.com/jaypipes/ghw/pkg/context"
)

func (i *Info) load(ctx *context.Context) error {
	// Check /dev/watchdog
	if _, err := os.Stat("/dev/watchdog"); err == nil {
		i.Present = true
		return nil
	}

	// Check /sys/class/watchdog/
	entries, err := os.ReadDir("/sys/class/watchdog")
	if err == nil && len(entries) > 0 {
		i.Present = true
	}

	return nil
}
