package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"storage-manager-api/internal/config"
	"storage-manager-api/internal/handlers"
	"storage-manager-api/internal/middleware"
	"storage-manager-api/internal/repository"
	"storage-manager-api/internal/services"
	"storage-manager-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	if err := os.MkdirAll(cfg.StoragePath, os.ModePerm); err != nil {
		log.Fatalf("Gagal membuat direktori storage (%s): %v", cfg.StoragePath, err)
	}

	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Inisialisasi database gagal: %v", err)
	}

	gin.SetMode(cfg.GinMode)
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	if cfg.MaxUploadSize > 0 {
		router.MaxMultipartMemory = cfg.MaxUploadSize
	}

	// Dependency Injection: Autentikasi
	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	// Dependency Injection: File Management
	fileRepo := repository.NewFileRepository(db)
	fileService := services.NewFileService(fileRepo, cfg)
	fileHandler := handlers.NewFileHandler(fileService)

	// Routing API v1
	apiV1 := router.Group("/api/v1")
	{
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		// Group route terproteksi middleware JWT
		protected := apiV1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			protected.GET("/auth/me", func(c *gin.Context) {
				utils.SuccessResponse(c, http.StatusOK, "Profil user aktif", gin.H{
					"id":    c.GetString(middleware.ContextUserIDKey),
					"email": c.GetString(middleware.ContextEmailKey),
				})
			})

			filesGroup := protected.Group("/files")
			{
				filesGroup.POST("/upload", fileHandler.UploadFile)
				filesGroup.GET("", fileHandler.ListFiles)
				filesGroup.GET("/download/:id", fileHandler.DownloadFile)
				filesGroup.DELETE("/:id", fileHandler.DeleteFile)
			}
		}
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Storage Manager API berjalan pada http://localhost%s (mode: %s)", addr, cfg.GinMode)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Gagal menjalankan server HTTP: %v", err)
	}
}



