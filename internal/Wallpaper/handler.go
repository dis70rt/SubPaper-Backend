package wallpaper

import (
	"net/http"

	reddit "github.com/dis70rt/subpaper-backend/internal/Reddit"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

type Handler struct {
    service *Service
}

func NewHandler(client *reddit.RedditClient, cache *cache.Cache) *Handler {
    return &Handler{
        service: NewService(client, cache),
    }
}

func (handler *Handler) SearchWallpapers(c *gin.Context) {
	var req WallpaperRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	posts, err := handler.service.FetchWallpaper(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}