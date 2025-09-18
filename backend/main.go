package main

import (
	"log"

	"cidr-viewer/handlers"
	_ "cidr-viewer/docs" // Import generated docs

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title CIDR Viewer API
// @version 1.0
// @description API for analyzing and validating CIDR ranges
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Cache-Control", "Pragma", "Expires"}
	r.Use(cors.New(config))

	// Swagger documentation at root path
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Redirect root to swagger docs
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	// API routes
	api := r.Group("/api")
	{
		api.POST("/analyze", handlers.AnalyzeCIDRs)
		api.POST("/validate", handlers.ValidateCIDR)
		api.GET("/health", handlers.HealthCheck)
	}

	// Start server
	log.Println("Starting CIDR Viewer API server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}