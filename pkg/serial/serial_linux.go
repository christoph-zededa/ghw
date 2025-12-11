// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package serial

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/linuxpath"
)

func (i *Info) load() error {
	var errs []error

	i.Devices, errs = serials(i.ctx)

	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("error(s) happened during reading serial info: %+v", errs)
}

func serials(ctx *context.Context) ([]*Device, []error) {
	devs := make([]*Device, 0)
	errs := []error{}

	paths := linuxpath.New(ctx)
	ttyDirs, err := os.ReadDir(paths.SysClassTty)
	if err != nil {
		// If the directory doesn't exist, it's not an error, just no devices
		if os.IsNotExist(err) {
			return devs, nil
		}
		return devs, []error{err}
	}

	id := 1
	for _, dir := range ttyDirs {
		ttyName := dir.Name()
		ttyPath := filepath.Join(paths.SysClassTty, ttyName)

		devicePath := filepath.Join(ttyPath, "device")
		resourcesPath := filepath.Join(devicePath, "resources")
		irqPath := filepath.Join(ttyPath, "irq")

		var io, irq string

		if _, err := os.Stat(resourcesPath); err == nil {
			io, irq = parseResources(resourcesPath)
		} else if runtime.GOARCH == "arm64" {
			if _, err := os.Stat(irqPath); err == nil {
				content, _ := os.ReadFile(irqPath)
				irqVal := strings.TrimSpace(string(content))
				if irqVal != "" && irqVal != "0" {
					irq = irqVal
				}
			}
		}

		if io != "" || irq != "" {
			dev := &Device{
				Name:    fmt.Sprintf("COM%d", id),
				Address: fmt.Sprintf("/dev/%s", ttyName),
				IO:      io,
				IRQ:     irq,
			}
			devs = append(devs, dev)
			id++
		}
	}

	return devs, errs
}

func parseResources(path string) (string, string) {
	f, err := os.Open(path)
	if err != nil {
		return "", ""
	}
	defer f.Close()

	var io, irq string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "io ") {
			// Replicate: sed -e 's#io 0x##' -e 's#0x##'
			line = strings.TrimPrefix(line, "io 0x")
			line = strings.Replace(line, "0x", "", -1)
			io = line
		} else if strings.HasPrefix(line, "irq ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				irq = parts[1]
			}
		}
	}
	return io, irq
}
