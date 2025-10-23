package wallpaper

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup) {
	w := rg.Group("/wallpapers")
	w.GET("/search", SearchWallpapers)
}
