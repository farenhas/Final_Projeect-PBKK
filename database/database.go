package database

import (
	"fmt"
	"photo_gallery/config"
	"photo_gallery/models"
)

func Migrate() {
	err := config.DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Photo{}, &models.ActivityLog{})
	if err != nil {
		panic("Failed to migrate database!")
	}
	fmt.Println("Database migration completed!")
}
