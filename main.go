package main

import (
	"99-api-public/config"
	"99-api-public/routes"
	"99-api-public/user"
	"99-api-public/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load environment variables
	config.Load()
	utils.InitURLs()

	// Create a new Gin engine
	r := gin.Default()

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.GetEnv("REDIS_ADDR", "localhost:6379"),
		Password: config.GetEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	// Get all users and cache the result
	users, err := user.GetAllUsers()
	if err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}

	err = user.CacheUsers(rdb, users)
	if err != nil {
		log.Fatalf("Failed to cache users: %v", err)
	}

	// Set up routes
	routes.SetupRoutes(r, rdb)

	// Start the server
	port := config.GetEnv("PORT", "8081")

	log.Printf("API Gateway running on port %s", port)
	r.Run(":" + port)
}
