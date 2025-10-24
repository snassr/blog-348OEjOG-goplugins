package pluginruntime

import (
	"context"

	"connectrpc.com/connect"

	pluginv1 "github.com/snassr/blog-348OEjOG-goplugins/external/gen/plugin-proto-go/plugin/v1"
	"github.com/snassr/blog-348OEjOG-goplugins/external/gen/plugin-proto-go/plugin/v1/pluginv1connect"
)

type PluginAdapter struct {
	client pluginv1connect.PluginServiceClient
}

func (a *PluginAdapter) Greet(ctx context.Context, name string) (string, error) {
	n := name
	req := pluginv1.GreetRequest{
		Name: &n,
	}

	resp, err := a.client.Greet(ctx, connect.NewRequest(&req))
	if err != nil {
		return "", err
	}

	return *resp.Msg.Message, nil
}

func (a *PluginAdapter) StreamGreet(ctx context.Context, name string, send func(msg string) error) error {
	req := connect.NewRequest(&pluginv1.StreamGreetRequest{
		Name: &name,
	})

	stream, err := a.client.StreamGreet(ctx, req)
	if err != nil {
		return err
	}

	for stream.Receive() {
		msg := stream.Msg().GetMessage()
		if err := send(msg); err != nil {
			return err
		}
	}

	if err := stream.Err(); err != nil {
		return err
	}

	return nil
}
