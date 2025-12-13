package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	
	_ "github.com/Candoo/thornton-pickard-api/docs"
	
	"github.com/Candoo/thornton-pickard-api/internal/database"
	"github.com/Candoo/thornton-pickard-api/internal/handlers"
	"github.com/Candoo/thornton-pickard-api/internal/middleware"
)

// @title Thornton Pickard Camera API
// @version 2.0
// @description Complete API for Thornton Pickard cameras and ephemera data with authentication, pagination, search, and image uploads
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Seed database if SEED=true
	if os.Getenv("SEED") == "true" {
		if err := database.SeedDatabase(db); err != nil {
			log.Println("Warning: Database seeding failed:", err)
		}
	}

	// Set Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())

	// Serve static files (uploads)
	r.Static("/uploads", "./uploads")

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "version": "2.0"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register(db))
			auth.POST("/login", handlers.Login(db))
			auth.GET("/profile", middleware.AuthRequired(), handlers.GetProfile(db))
		}

		// Public camera routes (read-only)
		cameras := v1.Group("/cameras")
		{
			cameras.GET("", handlers.GetCameras(db))
			cameras.GET("/:id", handlers.GetCamera(db))
			cameras.GET("/search", handlers.SearchCameras(db))
		}

		// Protected camera routes (require auth)
		camerasProtected := v1.Group("/cameras")
		camerasProtected.Use(middleware.AuthRequired())
		{
			camerasProtected.POST("", handlers.CreateCamera(db))
			camerasProtected.PUT("/:id", handlers.UpdateCamera(db))
			camerasProtected.DELETE("/:id", middleware.AdminRequired(), handlers.DeleteCamera(db))
		}

		// Ephemera routes
		ephemera := v1.Group("/ephemera")
		{
			ephemera.GET("", handlers.GetEphemera(db))
			ephemera.GET("/:id", handlers.GetEphemeraItem(db))
		}

		ephemeraProtected := v1.Group("/ephemera")
		ephemeraProtected.Use(middleware.AuthRequired())
		{
			ephemeraProtected.POST("", handlers.CreateEphemeraItem(db))
		}

		// Manufacturer routes
		manufacturers := v1.Group("/manufacturers")
		{
			manufacturers.GET("", handlers.GetManufacturers(db))
		}

		// Upload routes (require auth)
		upload := v1.Group("/upload")
		upload.Use(middleware.AuthRequired())
		{
			upload.POST("", handlers.UploadImage())
			upload.POST("/multiple", handlers.UploadMultipleImages())
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📚 API Documentation: http://localhost:%s/swagger/index.html", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}