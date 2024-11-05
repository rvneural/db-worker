package rest

import (
	"db-worker/internal/models/db"
	"db-worker/internal/service/keygen"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RestAPI struct {
	generator *keygen.Generator
	worker    db.DBWorker
}

func New(worker db.DBWorker, generator *keygen.Generator) *RestAPI {
	return &RestAPI{
		worker:    worker,
		generator: generator,
	}
}

func (r *RestAPI) GetID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": r.generator.Generate()})
}

func (r *RestAPI) RegisterOperation(c *gin.Context) {
	type Resuest struct {
		UniqID         string `json:"id"`
		Operation_type string `json:"type"`
	}
	request := Resuest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = r.worker.RegisterOperation(request.UniqID, request.Operation_type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
