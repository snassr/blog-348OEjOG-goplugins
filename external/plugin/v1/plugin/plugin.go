package plugin

import "context"

type Plugin interface {
	Greet(ctx context.Context, name string) (string, error)
}
