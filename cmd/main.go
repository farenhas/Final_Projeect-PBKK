package main

import (
	"photo_gallery/config"
	"photo_gallery/database"
	"photo_gallery/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	config.ConnectDatabase()

	// Run database migration
	database.Migrate()

	// Setup Gin engine
	router := gin.Default()

	// Tambahkan konfigurasi folder template
	router.LoadHTMLGlob("templates/*")

	// Static files
	router.Static("/assets", "./assets")

	// Setup routes
	routes.SetupRoutes(router)

	// Jalankan server
	router.Run(":8080")
}
