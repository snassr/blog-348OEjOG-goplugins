package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"

	adminv1 "github.com/snassr/blog-348OEjOG-goplugins/api/proto/gen/go/admin/v1"
	"github.com/snassr/blog-348OEjOG-goplugins/api/proto/gen/go/admin/v1/adminv1connect"
	pluginv1 "github.com/snassr/blog-348OEjOG-goplugins/external/plugin/v1/api/proto/gen/go/plugin/v1"
	"github.com/snassr/blog-348OEjOG-goplugins/external/plugin/v1/api/proto/gen/go/plugin/v1/pluginv1connect"
)

const (
	addr     string = "localhost:50101"
	hostAddr string = "http://localhost:8087"

	id string = "EXTERNAL_PLUGIN_GO"
)

func main() {
	client := adminv1connect.NewAdminServiceClient(http.DefaultClient, hostAddr)

	// register plugin
	pluginID := id
	addr := addr

	req := connect.NewRequest(&adminv1.RegisterPluginRequest{
		Id:      &pluginID,
		Address: &addr,
	})

	resp, err := client.RegisterPlugin(context.Background(), req)
	if err != nil {
		slog.Error(
			"failed to register plugin",
			slog.String("plugin_id", pluginID),
			slog.String("err", err.Error()),
		)
		return
	}

	status := *resp.Msg.Status
	slog.Info(
		"successfully registered plugin",
		slog.String("plugin_id", pluginID),
		slog.String("plugin_status", status),
	)

	// run server
	mux := http.NewServeMux()
	path, handler := pluginv1connect.NewPluginServiceHandler(&GreeterPlugin{})
	mux.Handle(path, handler)

	slog.Info(
		"starting %s plugin server",
		slog.String("plugin_id", pluginID),
		slog.String("addr", addr),
	)

	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error(
			"failed to start plugin server",
			slog.String("plugin_id", pluginID),
			slog.String("err", err.Error()),
		)
		return
	}
}

type GreeterPlugin struct{}

func (g *GreeterPlugin) Greet(
	ctx context.Context,
	req *connect.Request[pluginv1.GreetRequest],
) (*connect.Response[pluginv1.GreetResponse], error) {
	name := req.Msg.GetName()
	msg := fmt.Sprintf("你好, %s", name)
	return connect.NewResponse(&pluginv1.GreetResponse{Message: &msg}), nil
}
