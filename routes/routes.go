package routes

import (
	"99-api-public/handlers"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(r *gin.Engine, rdb *redis.Client) {
	api := r.Group("/public-api")

	// Listing Service Routes
	listingService := api.Group("/listings")
	listingService.GET("/", handlers.ProxyListingServiceGet)
	listingService.POST("/", handlers.ProxyListingServicePost)

	// User Service Routes
	userService := api.Group("/users")
	userService.POST("/", handlers.ProxyUserServicePost)
}
