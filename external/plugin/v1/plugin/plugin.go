package plugin

import "context"

type Plugin interface {
	Greet(ctx context.Context, name string) (string, error)
	StreamGreet(ctx context.Context, name string, send func(msg string) error) error
}
