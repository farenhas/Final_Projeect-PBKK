package controllers

import (
	"fmt"
	"net/http"
	"photo_gallery/config"
	"photo_gallery/models"

	"github.com/gin-gonic/gin"
)

type ActivitySummary struct {
	Entity string `json:"entity"`  
	Action string `json:"action"`  
	Count  int    `json:"count"`   
}

func TestImport(c *gin.Context) {
	fmt.Println("Import berhasil:", models.ActivityLog{})
	c.JSON(http.StatusOK, gin.H{"message": "Import berhasil"})
}

func GetActivityLogSummary(c *gin.Context) {
	var summary []ActivitySummary

	query := `
		SELECT entity, action, COUNT(*) AS count
		FROM activity_logs
		GROUP BY entity, action
	`

	if err := config.DB.Raw(query).Scan(&summary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
