package plugin_en

import (
	"context"
	"fmt"

	"github.com/snassr/blog-348OEjOG-goplugins/external/plugin/v1/plugin"
)

type pn struct{}

func New() plugin.Plugin {
	return &pn{}
}

func (p *pn) Greet(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Hello, %s", name), nil
}
