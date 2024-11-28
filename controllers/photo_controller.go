package controllers

import (
	"os"
	"log"
	"fmt"
	"time"
	"net/http"
	"photo_gallery/config"
	"photo_gallery/models"

	"github.com/gin-gonic/gin"
)

// Get all photos
func GetPhotos(c *gin.Context) {
	var photos []models.Photo
	// Mengambil semua foto dan menggabungkan dengan kategori mereka
	if err := config.DB.Preload("Category").Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load photos"})
		return
	}
	// Mengirim data ke halaman dashboard.html
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Photos": photos,
	})
}
func AddPhoto(c *gin.Context) {
	var input struct {
		Title       string `form:"title" binding:"required"`
		Description string `form:"description"`
		Category    string `form:"category" binding:"required"` // Menggunakan nama kategori sebagai input
	}

	// Validasi input dari pengguna
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi file gambar
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Validasi jenis file (hanya JPEG dan PNG yang diizinkan)
	if file.Header.Get("Content-Type") != "image/jpeg" && file.Header.Get("Content-Type") != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG and PNG files are allowed"})
		return
	}

	// Validasi ukuran file (maksimal 5MB)
	const MaxFileSize = 5 << 20 // 5MB
	if file.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 5MB"})
		return
	}

	// Cek apakah kategori sudah ada, jika tidak buat kategori baru
	var category models.Category
	if err := config.DB.Where("name = ?", input.Category).First(&category).Error; err != nil {
		// Jika kategori tidak ditemukan, buat kategori baru
		category = models.Category{
			Name:      input.Category,
			CreatedAt: time.Now(),
		}
		if err := config.DB.Create(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
			return
		}
	}

	// Log kategori yang digunakan untuk memastikan bahwa ID benar
	log.Printf("Category used: %+v", category)

	// Simpan file ke folder uploads
	filePath := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Simpan metadata foto ke dalam tabel photos
	photo := models.Photo{
		Title:       input.Title,
		Description: input.Description,
		CategoryID:  category.ID, // Menggunakan ID kategori yang baru dibuat atau ditemukan
		ImagePath:   filePath,
		CreatedAt:   time.Now(),
	}
	if err := config.DB.Create(&photo).Error; err != nil {
		log.Printf("Error saving photo to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo metadata"})
		return
	}

	// Redirect ke dashboard setelah berhasil meng-upload foto
    c.Redirect(http.StatusFound, "/dashboard")
}

func EditPhoto(c *gin.Context) {
    photoID := c.Param("id")
    var photo models.Photo

    if err := config.DB.First(&photo, photoID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
        return
    }

    var input struct {
        Title       string `form:"title" binding:"required"`
        Description string `form:"description"`
    }

    if err := c.ShouldBind(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    photo.Title = input.Title
    photo.Description = input.Description

    config.DB.Save(&photo)

    c.Redirect(http.StatusFound, "/dashboard")
}

func DeletePhoto(c *gin.Context) {
    id := c.Param("id")

    // Cari foto berdasarkan ID
    var photo models.Photo
    if err := config.DB.First(&photo, id).Error; err != nil {
        log.Printf("Error finding photo: %v", err) // Logging error
        c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
        return
    }

    // Hapus foto dari folder uploads jika ada
    if err := os.Remove(photo.ImagePath); err != nil {
        log.Printf("Failed to remove photo from disk: %v", err)
    }

    // Hapus foto dari database
    if err := config.DB.Delete(&photo).Error; err != nil {
        log.Printf("Error deleting photo from database: %v", err) // Logging error
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo"})
        return
    }

    // Redirect kembali ke dashboard
    log.Printf("Photo with ID %s deleted successfully", id)
    c.Redirect(http.StatusFound, "/dashboard")
}
