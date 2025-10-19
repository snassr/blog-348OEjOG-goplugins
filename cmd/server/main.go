package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/snassr/blog-348OEjOG-goplugins/api/proto/gen/admin/v1/adminv1connect"
	"github.com/snassr/blog-348OEjOG-goplugins/internal/pluginruntime"
	"github.com/snassr/blog-348OEjOG-goplugins/internal/server"
)

const host string = ":8087"

func main() {
	fmt.Println("Hello World!")

	pm := pluginruntime.NewManager()

	adminHandler := server.NewAdminHandler(pm)

	mux := http.NewServeMux()

	path, handler := adminv1connect.NewAdminServiceHandler(adminHandler)
	mux.Handle(path, handler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "Accept", "X-Requested-With", "Connect-Protocol-Version", "connect-src"},
		Debug:          true,
	}).Handler(mux)

	srv := &http.Server{
		Addr:    host,
		Handler: h2c.NewHandler(corsHandler, &http2.Server{}),

		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to start admin server", slog.String("err", err.Error()))
		return
	}
}
