package routes

import (
	"net/http"
	"photo_gallery/controllers"
	"photo_gallery/models"
	"photo_gallery/config"

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


	// Route untuk halaman dashboard
	router.GET("/dashboard", func(c *gin.Context) {
    var photos []models.Photo
    // Mencoba memuat data foto dari database, menampilkan pesan error jika terjadi masalah
    if err := config.DB.Preload("Category").Find(&photos).Error; err != nil {
        c.HTML(http.StatusInternalServerError, "dashboard.html", gin.H{
            "error": "Unable to load photos from database",
            "Photos": []models.Photo{},
        })
        return
    }

    // Jika data berhasil di-load, kirimkan data ke template
    c.HTML(http.StatusOK, "dashboard.html", gin.H{
        "Photos": photos,
    })
})


	// Route untuk halaman upload
	router.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})

	// Route untuk proses upload foto
	router.POST("/api/upload", controllers.AddPhoto)

	// Route untuk halaman edit photo
router.GET("/edit/:id", func(c *gin.Context) {
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
	router.POST("/api/edit/:id", controllers.EditPhoto)

	// Route untuk menghapus foto
	router.POST("/api/delete/:id", controllers.DeletePhoto)

	// Route untuk proses login
	router.POST("/api/login", controllers.LoginUser)

	// Route untuk logout
	router.GET("/logout", func(c *gin.Context) {
		// Add any session clearing logic if necessary
		c.Redirect(http.StatusFound, "/login")
	})
}
