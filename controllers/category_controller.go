package controllers

import (
	"net/http"
	"photo_gallery/config"
	"photo_gallery/models"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	var categories []models.Category
	config.DB.Find(&categories)
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func AddCategory(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{Name: input.Name}
	config.DB.Create(&category)
	c.JSON(http.StatusCreated, gin.H{"data": category})
}
