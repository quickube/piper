package server

import (
	"github.com/quickube/piper/pkg/clients"
	"github.com/quickube/piper/pkg/conf"
	"golang.org/x/net/context"
	"log"
)

func Start(ctx context.Context, stop context.CancelFunc, cfg *conf.GlobalConfig, clients *clients.Clients) {

	srv := NewServer(cfg, clients)
	gracefulShutdownHandler := NewGracefulShutdown(ctx, stop)
	srv.Start(ctx)

	gracefulShutdownHandler.Shutdown(srv)

	log.Println("Server exiting")
}
