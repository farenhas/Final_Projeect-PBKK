package routes

import (
	"log"
	"encoding/json"
	"net/http"
	"photo_gallery/config"
	"photo_gallery/middlewares"
	"photo_gallery/controllers"
	"photo_gallery/models"
	

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Load templates
	router.LoadHTMLGlob("templates/*")

	// Static files (CSS, JS, Images)
	router.Static("/uploads", "./uploads")

	// Route untuk halaman login
	router.GET("/login", func(c *gin.Context) {
		successMessage := ""
		if c.Query("success") == "1" {
			successMessage = "Successfully registered. Please login."
		}
		c.HTML(http.StatusOK, "login.html", gin.H{
			"success": successMessage,
		})
	})

	// Route untuk halaman signup
	router.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})

	// Route untuk proses registrasi
	router.POST("/api/register", controllers.RegisterUser)

	
	router.POST("/api/login", controllers.LoginUser)

	// Route untuk logout
	router.GET("/logout", func(c *gin.Context) {
	// Hapus cookie Authorization
	c.SetCookie("Authorization", "", -1, "/", "", false, true) 

	// Redirect ke halaman login
	c.Redirect(http.StatusFound, "/login")
})

	// Grup rute yang membutuhkan autentikasi
	authRoutes := router.Group("/")
	authRoutes.Use(middlewares.AuthMiddleware()) // Middleware autentikasi diterapkan
	{
		authRoutes.GET("/api/activity-summary", func(c *gin.Context) {
			var summary []struct {
				Entity string `json:"entity"`
				Action string `json:"action"`
				Count  int    `json:"count"`
			}
		
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
		})
		
		// Route untuk halaman dashboard
		authRoutes.GET("/dashboard", func(c *gin.Context) {
			var photos []models.Photo
			query := c.Query("query")
			dbQuery := config.DB.Preload("Category")
			if query != "" {
				dbQuery = dbQuery.Joins("JOIN categories ON categories.id = photos.category_id").
					Where("photos.title LIKE ? OR photos.description LIKE ? OR categories.name LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
			}
			if err := dbQuery.Find(&photos).Error; err != nil {
				c.HTML(http.StatusInternalServerError, "dashboard.html", gin.H{
					"error":  "Unable to load photos from database",
					"Photos": []models.Photo{},
				})
				return
			}
			c.HTML(http.StatusOK, "dashboard.html", gin.H{
				"Photos":      photos,
				"SearchQuery": query,
			})
		})

		// Route untuk halaman upload
		authRoutes.GET("/upload", func(c *gin.Context) {
			c.HTML(http.StatusOK, "upload.html", nil)
		})

		// Route untuk proses upload foto
		authRoutes.POST("/api/upload", controllers.AddPhoto)

		// Route untuk halaman edit photo
		authRoutes.GET("/edit/:id", func(c *gin.Context) {
			photoID := c.Param("id")
			var photo models.Photo
			if err := config.DB.First(&photo, photoID).Error; err != nil {
				c.HTML(http.StatusNotFound, "dashboard.html", gin.H{
					"error": "Photo not found",
				})
				return
			}
			c.HTML(http.StatusOK, "edit.html", gin.H{
				"Photo": photo,
			})
		})

		// Route untuk proses edit photo
		authRoutes.POST("/api/edit/:id", controllers.EditPhoto)

		// Route untuk menghapus foto
		authRoutes.POST("/api/delete/:id", controllers.DeletePhoto)

		// Route untuk statistik
		authRoutes.GET("/statistics", func(c *gin.Context) {
			var logs []models.ActivityLog
		
			// Ambil semua logs dari database
			if err := config.DB.Find(&logs).Error; err != nil {
				log.Printf("Error fetching logs: %v", err)
				c.HTML(http.StatusInternalServerError, "statistics.html", gin.H{
					"error": "Failed to fetch activity logs",
					"Logs":  nil,
				})
				return
			}
		
			// Hitung jumlah aktivitas berdasarkan jenis tindakan
			activityCounts := map[string]int{
				"Create": 0,
				"Read":   0,
				"Update": 0,
				"Delete": 0,
			}
		
			for _, log := range logs {
				if _, ok := activityCounts[log.Action]; ok {
					activityCounts[log.Action]++
				}
			}
		
			// Encode activityCounts menjadi JSON string
			activityCountsJSON, err := json.Marshal(activityCounts)
			if err != nil {
				log.Printf("Error encoding activityCounts: %v", err)
				activityCountsJSON = []byte("{}")
			}
		
			// Kirim data ke template
			c.HTML(http.StatusOK, "statistics.html", gin.H{
				"Logs":           logs,
				"ActivityCounts": string(activityCountsJSON), // Kirim sebagai string JSON
			})
		})
		
		
	}
}