package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	adminv1 "github.com/snassr/blog-348OEjOG-goplugins/external/gen/goplugins-api-proto-go/admin/v1"
	"github.com/snassr/blog-348OEjOG-goplugins/internal/pluginruntime"
)

type AdminHandler struct {
	pm *pluginruntime.Manager
}

func NewAdminHandler(pm *pluginruntime.Manager) *AdminHandler {
	return &AdminHandler{
		pm: pm,
	}
}

func (h *AdminHandler) RegisterPlugin(ctx context.Context, req *connect.Request[adminv1.RegisterPluginRequest]) (*connect.Response[adminv1.RegisterPluginResponse], error) {
	err := h.pm.RegisterRemote(*req.Msg.Id, *req.Msg.Address)
	if err != nil {
		return nil, err
	}

	status := "REGISTERED"

	return &connect.Response[adminv1.RegisterPluginResponse]{
		Msg: &adminv1.RegisterPluginResponse{
			Status: &status,
		},
	}, nil

}

func (h *AdminHandler) AllGreetings(ctx context.Context, req *connect.Request[adminv1.AllGreetingsRequest]) (*connect.Response[adminv1.AllGreetingsResponse], error) {
	list := h.pm.List()
	greetings := []string{}

	for _, id := range list {
		p, ok := h.pm.Get(id)
		if !ok {
			return nil, fmt.Errorf("not found")
		}

		name := *req.Msg.Name

		g, err := p.Greet(ctx, name)
		if err != nil {
			return nil, err
		}

		greetings = append(greetings, g)
	}

	return &connect.Response[adminv1.AllGreetingsResponse]{
		Msg: &adminv1.AllGreetingsResponse{
			Messages: greetings,
		},
	}, nil
}
