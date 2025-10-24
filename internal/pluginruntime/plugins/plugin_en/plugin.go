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

func (p *pn) StreamGreet(ctx context.Context, name string, send func(msg string) error) error {
	for i := 1; i <= 5; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg := fmt.Sprintf("Hello #%d to %s from internal plugin", i, name)
			if err := send(msg); err != nil {
				return err
			}
		}
	}

	return nil
}
