package wallpaper

import (
	reddit "github.com/dis70rt/subpaper-backend/internal/Reddit"
	"github.com/dis70rt/subpaper-backend/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)


func RegisterRoutes(router *gin.RouterGroup, client *reddit.RedditClient, cache *cache.Cache) {
	handler := NewHandler(client, cache)
	w := router.Group("/wallpapers")
	w.Use(middlewares.GzipMiddleware())
	w.GET("/search", handler.SearchWallpapers)
}
