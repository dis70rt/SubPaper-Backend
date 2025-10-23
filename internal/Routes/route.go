package routes

import (
	wallpaper "github.com/dis70rt/subpaper-backend/internal/Wallpaper"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/v1")
	wallpaper.RegisterRoutes(api)
	return router
}