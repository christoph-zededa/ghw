//go:build !linux
// +build !linux

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package serial

import (
	"fmt"
	"runtime"
)

func (i *Info) load() error {
	return fmt.Errorf("serial load not implemented on %s", runtime.GOOS)
}
