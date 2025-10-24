package internal

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type APIInfo struct {
	Name          string      `json:"name"`
    Version       string      `json:"version"`
    Description   string      `json:"description"`
    Status        string      `json:"status"`
	Endpoints     interface{} `json:"endpoints"`
}

type HealthCheck struct {
    Status    string `json:"status"`
    Timestamp int64  `json:"timestamp"`
    Uptime    int64  `json:"uptime"`
}

var startTime = time.Now()

func GetRootInfo(c *gin.Context) {
    c.JSON(http.StatusOK, APIInfo{
        Name:        "Subpaper API",
        Version:     "2.0.0",
        Description: "Wallpaper API powered by Reddit",
        Status:      "operational",
        Endpoints: gin.H{
            "health":     "/health",
            "api_info":   "/api/v1",
            "wallpapers": "/api/v1/wallpapers/search",
        },
    })
}

func GetHealth(c *gin.Context) {
    c.JSON(http.StatusOK, HealthCheck{
        Status:    "healthy",
        Timestamp: time.Now().Unix(),
        Uptime:    int64(time.Since(startTime).Seconds()),
    })
}
