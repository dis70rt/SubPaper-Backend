package middlewares

import (
    "time"
    
    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
)

func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery
        
        c.Next()
        
        latency := time.Since(start)
        statusCode := c.Writer.Status()
        
        if raw != "" {
            path = path + "?" + raw
        }
        
        log.WithFields(log.Fields{
            "status":     statusCode,
            "method":     c.Request.Method,
            "path":       path,
            "ip":         c.ClientIP(),
            "latency":    latency,
            "user-agent": c.Request.UserAgent(),
        }).Info("Request processed")
    }
}