package controllers

import (
	"net/http"
	"time"

	"photo_gallery/config"
	"photo_gallery/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var JwtSecret = []byte("secret_key")

func RegisterUser(c *gin.Context) {
	
	username := c.PostForm("username")
	password := c.PostForm("password")

	
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
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		return
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		return
	}

	// Buat token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Token kedaluwarsa dalam 72 jam
	})

	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Simpan token di cookie
	c.SetCookie("Authorization", "Bearer "+tokenString, 3600*24*3, "/", "", false, true)

	// Redirect ke dashboard
	c.Redirect(http.StatusFound, "/dashboard")
}

// Logout user
func LogoutUser(c *gin.Context) {
	// Hapus cookie Authorization
	c.SetCookie("Authorization", "", -1, "/", "", false, true) // Expire the cookie

	// Redirect ke halaman login
	c.Redirect(http.StatusFound, "/login")
}

