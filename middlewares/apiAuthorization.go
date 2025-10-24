package middlewares

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
)

func APIAuthMiddleware() gin.HandlerFunc {
    secretKey := os.Getenv("API_SECRET_KEY")
    
    if secretKey == "" {
        secretKey = "X-_zCtQ44jM6HPIfdLNljiRJrNU31ODaOTPyWx1HZY6G0Lu0wBOFiY9zBpkpU31k"
    }
    
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader != "" {
            const bearerPrefix = "Bearer "
            if len(authHeader) > len(bearerPrefix) {
                token := authHeader[len(bearerPrefix):]
                if token == secretKey {
                    c.Next()
                    return
                }
            }
        }
        
        apiKey := c.GetHeader("X-API-Key")
        if apiKey == secretKey {
            c.Next()
            return
        }
        
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Unauthorized: Invalid or missing API key",
        })
        c.Abort()
    }
}