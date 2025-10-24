package main

import (
	"time"

	reddit "github.com/dis70rt/subpaper-backend/internal/Reddit"
	wallpaper "github.com/dis70rt/subpaper-backend/internal/Wallpaper"
	"github.com/dis70rt/subpaper-backend/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

func main() {
	router := gin.Default()
	client := reddit.NewClient()
	cache := cache.New(12*time.Hour, 1*time.Hour)

	api := router.Group("/api/v1")
	
	router.Use(middlewares.RateLimitMiddleware(30,6))
	router.Use(middlewares.APIAuthMiddleware())

	wallpaper.RegisterRoutes(api, client, cache)

	router.Run(":8080")
}
