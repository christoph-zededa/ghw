// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package serial

import (
	"fmt"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/option"
)

type Device struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	IO      string `json:"io"`
	IRQ     string `json:"irq"`
}

func (d Device) String() string {
	return fmt.Sprintf("%s (%s)", d.Name, d.Address)
}

type Info struct {
	ctx     *context.Context
	Devices []*Device `json:"devices"`
}

func (i *Info) String() string {
	return fmt.Sprintf(
		"Serial (%d devices)",
		len(i.Devices),
	)
}

func New(opts ...*option.Option) (*Info, error) {
	ctx := context.New(opts...)
	info := &Info{ctx: ctx}
	if err := ctx.Do(info.load); err != nil {
		return nil, err
	}

	return info, nil
}

func (i *Info) YAMLString() string {
	return marshal.SafeYAML(i.ctx, i)
}

func (i *Info) JSONString(indent bool) string {
	return marshal.SafeJSON(i.ctx, i, indent)
}
