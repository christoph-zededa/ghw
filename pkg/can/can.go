package can

import (
	"fmt"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/option"
)

type PCIAddress struct {
	Domain   string `json:"domain"`
	Bus      string `json:"bus"`
	Device   string `json:"device"`
	Function string `json:"function"`
}

type USBAddress struct {
	Bus    string `json:"bus"`
	Devnum string `json:"devnum"`
}

type BusParent struct {
	PCI *PCIAddress `json:"pci,omitempty"`
	USB *USBAddress `json:"usb,omitempty"`
}

type Device struct {
	Name   string     `json:"name"`
	Parent *BusParent `json:"parent,omitempty"`
}

type Info struct {
	ctx     *context.Context
	Devices []*Device `json:"devices"`
}

func (i *Info) String() string {
	return fmt.Sprintf("can (%d devices)", len(i.Devices))
}

func New(opts ...*option.Option) (*Info, error) {
	ctx := context.New(opts...)
	info := &Info{ctx: ctx}
	if err := info.load(ctx); err != nil {
		return nil, err
	}
	return info, nil
}

func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(i.ctx, i, indent)
}

func (i *Info) YAMLString() string {
	return marshal.SafeYAML(i.ctx, i)
}
