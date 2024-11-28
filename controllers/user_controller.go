package controllers

import (
	"net/http"
	"photo_gallery/config"
	"photo_gallery/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var JwtSecret = []byte("secret_key")
func RegisterUser(c *gin.Context) {
	// Ambil data dari form
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Simpan user ke database
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
	}
	if err := config.DB.Create(&user).Error; err != nil {
		// Jika ada error, tetap di halaman signup dengan pesan error
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"error": "Registration failed. Please try again.",
		})
		return
	}

	// Redirect ke halaman login dengan pesan sukses
	c.Redirect(http.StatusFound, "/login?success=1")
}


func LoginUser(c *gin.Context) {
	// Ambil data dari form
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Cari user berdasarkan username
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Redirect ke dashboard
	c.Redirect(http.StatusFound, "/dashboard")
}