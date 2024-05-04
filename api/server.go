package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/online-tryout/parsing-sheets-api/db/sqlc"
	// "github.com/online-tryout/parsing-sheets-api/docs"
	"github.com/online-tryout/parsing-sheets-api/broker"
	"github.com/online-tryout/parsing-sheets-api/util"
	// swaggerfiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	rabbitmq   broker.RabbitMq
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewServer(config *util.Config, store db.Store, rmq *broker.RabbitMq) (*Server, error) {
	server := &Server{config: *config, store: store, rabbitmq: *rmq}
	server.setupRouter()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/api/tryout", server.createTryout)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
