package middlewares

import (
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

type tokenBucket struct {
    tokens         float64
    capacity       float64
    refillRate     float64
    lastRefillTime time.Time
    mu             sync.Mutex
}

func newTokenBucket(capacity, refillRate float64) *tokenBucket {
    return &tokenBucket{
        tokens:         capacity,
        capacity:       capacity,
        refillRate:     refillRate,
        lastRefillTime: time.Now(),
    }
}

func (tb *tokenBucket) allowRequest() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.lastRefillTime).Seconds()

    tb.tokens += elapsed * tb.refillRate
    if tb.tokens > tb.capacity {
        tb.tokens = tb.capacity
    }
    tb.lastRefillTime = now

    if tb.tokens >= 1 {
        tb.tokens--
        return true
    }

    return false
}

type rateLimiter struct {
    buckets map[string]*tokenBucket
    mu      sync.RWMutex
}

func newRateLimiter() *rateLimiter {
    rl := &rateLimiter{
        buckets: make(map[string]*tokenBucket),
    }

    go rl.cleanup()

    return rl
}

func (rl *rateLimiter) cleanup() {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        rl.mu.Lock()
        now := time.Now()
        for ip, bucket := range rl.buckets {
            bucket.mu.Lock()
            if now.Sub(bucket.lastRefillTime) > 30*time.Minute {
                delete(rl.buckets, ip)
            }
            bucket.mu.Unlock()
        }
        rl.mu.Unlock()
    }
}

func (rl *rateLimiter) getBucket(ip string, capacity, refillRate float64) *tokenBucket {
    rl.mu.RLock()
    bucket, exists := rl.buckets[ip]
    rl.mu.RUnlock()

    if exists {
        return bucket
    }

    rl.mu.Lock()
    bucket, exists = rl.buckets[ip]
    if !exists {
        bucket = newTokenBucket(capacity, refillRate)
        rl.buckets[ip] = bucket
    }
    rl.mu.Unlock()

    return bucket
}

var limiter = newRateLimiter()

func RateLimitMiddleware(capacity, refillRate float64) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()

        bucket := limiter.getBucket(ip, capacity, refillRate)

        if !bucket.allowRequest() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded. Please try again later.",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}