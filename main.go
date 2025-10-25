package main

import (
	"time"

	"github.com/dis70rt/subpaper-backend/internal"
	reddit "github.com/dis70rt/subpaper-backend/internal/Reddit"
	wallpaper "github.com/dis70rt/subpaper-backend/internal/Wallpaper"
	"github.com/dis70rt/subpaper-backend/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middlewares.LoggingMiddleware())
	router.Use(middlewares.CORSMiddleware())

	client := reddit.NewClient()
	cache := cache.New(12*time.Hour, 1*time.Hour)

	router.GET("/", internal.GetRootInfo)
    router.GET("/health", internal.GetHealth)
	router.GET("/metrics", internal.GetMetrics)

	router.Use(middlewares.RateLimitMiddleware(30,5))
	router.Use(middlewares.APIAuthMiddleware())

	api := router.Group("/api/v1")

	wallpaper.RegisterRoutes(api, client, cache)

	log.Info("Starting server on :8080")
	router.Run(":8080")
}
