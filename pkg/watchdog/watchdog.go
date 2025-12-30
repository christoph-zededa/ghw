package watchdog

import (
	"fmt"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/marshal"
	"github.com/jaypipes/ghw/pkg/option"
)

type Info struct {
	ctx     *context.Context
	Present bool `json:"present"`
}

func (i *Info) String() string {
	return fmt.Sprintf("Watchdog present: %v", i.Present)
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
