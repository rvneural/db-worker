package app

import (
	config "db-worker/internal/config/app"
	"os"

	"log/slog"

	sloggin "github.com/samber/slog-gin"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/dvwright/xss-mw"
)

type Endpoint struct {
	router    *gin.Engine
	operation *gin.RouterGroup
	logger    *slog.Logger
}

func New() *Endpoint {

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(gin.Recovery())
	operationRouter := router.Group("/operation")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	router.Use(sloggin.New(logger))

	var xssMdlwr xss.XssMw
	router.Use(xssMdlwr.RemoveXss())

	return &Endpoint{
		router:    router,
		operation: operationRouter,
		logger:    logger,
	}
}

func (e *Endpoint) GetLogger() *slog.Logger {
	return e.logger
}

func (e *Endpoint) AddPostHandler(pattern string, handler gin.HandlerFunc) {
	e.router.POST(pattern, handler)
}

func (e *Endpoint) AddGetHandler(pattern string, handler gin.HandlerFunc) {
	e.router.GET(pattern, handler)
}

func (e *Endpoint) AddOperationPostHandler(pattern string, handler gin.HandlerFunc) {
	e.operation.POST(pattern, handler)
}

func (e *Endpoint) AddOperationGetHandler(pattern string, handler gin.HandlerFunc) {
	e.operation.GET(pattern, handler)
}

func (e *Endpoint) Start() {
	e.logger.Error(e.router.Run(":" + config.PORT).Error())
}
