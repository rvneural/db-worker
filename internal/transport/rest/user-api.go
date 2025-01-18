package rest

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (r *RestAPI) RegisterNewUser(c *gin.Context) {
	type DBUser struct {
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
	}
	model := DBUser{}
	err := c.BindJSON(&model)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	exists, err := r.worker.CheckEmail(model.Email)
	if err != nil {
		log.Println("Error checking user exists:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(400, gin.H{"error": "User already exists"})
		return
	}
	id, err := r.worker.RegisterNewUser(model.Email, model.Password, model.FirstName, model.LastName)
	if err != nil {
		log.Println("Error registering new user:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"id": id})
}

func (r *RestAPI) ComparePassword(c *gin.Context) {
	type DBUser struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	model := DBUser{}
	err := c.BindJSON(&model)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ok, id, err := r.worker.CheckCorrectPassword(model.Email, model.Password)
	if err != nil {
		log.Println("Error comparing password:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": ok,
		"id":     id,
	})
}

func (r *RestAPI) CheckExists(c *gin.Context) {
	type Request struct {
		Email string `json:"email" binding:"required"`
	}
	model := Request{}
	err := c.BindJSON(&model)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	email := model.Email
	if !strings.HasSuffix(email, "@realnoevremya.ru") {
		email += "@realnoevremya.ru"
	}
	log.Println("Checking user exists:", email)
	if !strings.HasSuffix(email, "@realnoevremya.ru") {
		email += "@realnoevremya.ru"
	}
	ok, err := r.worker.CheckEmail(email)
	if err != nil {
		log.Println("Error checking user exists:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": ok})
}

func (r *RestAPI) GetUser(c *gin.Context) {
	type Request struct {
		Email string `json:"email" binding:"required"`
	}
	model := Request{}
	err := c.BindJSON(&model)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	email := model.Email
	if !strings.HasSuffix(email, "@realnoevremya.ru") {
		email += "@realnoevremya.ru"
	}
	user, err := r.worker.GetUserByEmail(email)
	if err != nil {
		log.Println("Error getting user:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func (r *RestAPI) GetUserByID(c *gin.Context) {
	str_id := c.Param("id")
	if str_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	user_id, err := strconv.Atoi(str_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := r.worker.GetUserByID(user_id)
	if err != nil {
		log.Println("Error getting user:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func (r *RestAPI) GetAllUsers(c *gin.Context) {
	users, err := r.worker.GetAllUsers()
	if err != nil {
		log.Println("Error getting all users:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"users": users})
}
