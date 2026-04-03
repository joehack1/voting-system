package main

import (
    "fmt"
    "log"
    "os"

    "voting-system/db"
    "voting-system/handlers"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    godotenv.Load()

    // Initialize database
    db.InitDB()
    defer db.DB.Close()

    // Create Gin router
    r := gin.Default()

    // Enable CORS for frontend
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    // API routes
    api := r.Group("/api")
    {
        api.POST("/polls", handlers.CreatePoll)
        api.GET("/polls/:id", handlers.GetPoll)
        api.POST("/vote", handlers.Vote)
        api.GET("/results/:pollId", handlers.GetResults)
    }

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "OK"})
    })

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    fmt.Printf("🚀 Server running on http://localhost:%s\n", port)
    log.Fatal(r.Run(":" + port))
}