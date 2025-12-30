//go:build !linux
// +build !linux

package tpm

import (
	"github.com/jaypipes/ghw/pkg/context"
)

func (i *Info) load(ctx *context.Context) error {
	return nil
}
