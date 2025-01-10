package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/quickube/piper/pkg/clients"
	"github.com/quickube/piper/pkg/conf"
	"github.com/quickube/piper/pkg/server/routes"
	"github.com/quickube/piper/pkg/webhook_creator"
	"log"
	"net/http"
)

func NewServer(config *conf.GlobalConfig, clients *clients.Clients) *Server {
	srv := &Server{
		router:         gin.New(),
		config:         config,
		clients:        clients,
		webhookCreator: webhook_creator.NewWebhookCreator(config, clients),
	}

	return srv
}

func (s *Server) startServer() *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}

	go func() {
		log.Printf("Server is listening on %s", s.httpServer.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return srv
}

func (s *Server) registerMiddlewares() {
	s.router.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/healthz", "/readyz"},
		}),
		gin.Recovery(),
	)

}

func (s *Server) getRoutes() {
	v1 := s.router.Group("/")
	routes.AddReadyRoutes(v1)
	routes.AddHealthRoutes(v1, s.webhookCreator, s.config)
	routes.AddWebhookRoutes(s.config, s.clients, v1, s.webhookCreator)
}

func (s *Server) startServices(ctx context.Context) {
	s.webhookCreator.Start(ctx)
}

func (s *Server) Start(ctx context.Context) {

	s.registerMiddlewares()

	s.getRoutes()

	s.httpServer = s.startServer()

	s.startServices(ctx)

}
