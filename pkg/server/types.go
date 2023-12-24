package server

import (
	"github.com/gin-gonic/gin"
	"github.com/quickube/piper/pkg/clients"
	"github.com/quickube/piper/pkg/conf"
	"github.com/quickube/piper/pkg/webhook_creator"
	"net/http"
)

type Server struct {
	router         *gin.Engine
	config         *conf.GlobalConfig
	clients        *clients.Clients
	webhookCreator *webhook_creator.WebhookCreatorImpl
	httpServer     *http.Server
}

type Interface interface {
	startServer() *http.Server
	registerMiddlewares()
	getRoutes()
	Start() *http.Server
}
