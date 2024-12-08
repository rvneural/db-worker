package rest

import (
	"db-worker/internal/models/db"
	"db-worker/internal/service/keygen"
	"io"
	"net/http"
	"strconv"

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
		UserID         int    `json:"user_id"`
	}
	request := Resuest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = r.worker.RegisterOperation(request.UniqID, request.Operation_type, request.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (r *RestAPI) GetOperation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is empty"})
		return
	}
	operation, err := r.worker.GetOperation(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, operation)
}

func (r *RestAPI) GetAllOperations(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "0")
	operation_type := c.DefaultQuery("type", "")
	str_user_id := c.DefaultQuery("userid", "")
	operation_id := c.DefaultQuery("id", "")
	if operation_id == "" {
		var limit int = 0
		var errL error
		limit, errL = strconv.Atoi(limitStr)
		if errL != nil {
			limit = 0
		}
		user_id, err := strconv.Atoi(str_user_id)
		if err != nil {
			user_id = -1
		}
		operations, err := r.worker.GetAllOperations(limit, operation_type, user_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"operations": operations})
	} else {
		operation, err := r.worker.GetOperation(operation_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"operations": []db.DBResult{operation}})
	}
}

func (r *RestAPI) GetVersion(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is empty"})
		return
	}
	version, err := r.worker.GetVersion(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"version": version})
}

func (r *RestAPI) SetResult(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is empty"})
		return
	}
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = r.worker.SetResult(id, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
