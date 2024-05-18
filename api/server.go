package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/online-tryout/parsing-sheets-api/broker"
	db "github.com/online-tryout/parsing-sheets-api/db/sqlc"
	"github.com/online-tryout/parsing-sheets-api/docs"
	"github.com/online-tryout/parsing-sheets-api/util"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

func NewServer(config *util.Config, rmq *broker.RabbitMq) (*Server, error) {
	server := &Server{config: *config, rabbitmq: *rmq}
	server.setupRouter()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// configure swagger docs
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = server.config.BackendSwaggerHost
	router.GET("/api/parsing-sheets/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// health check api
	router.GET("/api/parsing-sheets/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "server is running"})
	})

	router.POST("/api/parsing-sheets/parse", server.parsingSheets)
	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
