package main

import (
	"photo_gallery/config"
	"photo_gallery/database"
	"photo_gallery/routes"
	"encoding/json"
	"html/template"

	"github.com/gin-gonic/gin"

)

func toJSON(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(jsonData)
}


func main() {
	// Connect to database
	config.ConnectDatabase()

	// Run database migration
	database.Migrate()

	// Setup Gin engine
	router := gin.Default()

	router.SetFuncMap(template.FuncMap{
		"toJSON": toJSON,
	})

	// Tambahkan konfigurasi folder template
	router.LoadHTMLGlob("templates/*")

	// Static files
	router.Static("/assets", "./assets")

	// Setup routes
	routes.SetupRoutes(router)

	// Jalankan server
	router.Run(":8080")
}


