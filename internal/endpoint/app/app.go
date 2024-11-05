package app

import (
	config "db-worker/internal/config/app"
	"log"

	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	router    *gin.Engine
	operation *gin.RouterGroup
}

func New() *Endpoint {
	router := gin.Default()
	operationRouter := router.Group("/operation")
	return &Endpoint{
		router:    router,
		operation: operationRouter,
	}
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
	log.Fatal(e.router.Run(":" + config.PORT))
}
