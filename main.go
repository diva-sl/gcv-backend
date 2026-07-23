package main

import (
	"log"
	"os"

	"gcv-backend/config"
	"gcv-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environments")
	}

	// Initialize MongoDB database connection
	config.ConnectDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	router := gin.Default()

	// Configure CORS Middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup Routes
	routes.RegisterRoutes(router)

	log.Printf("Server running securely on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}