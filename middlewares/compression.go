package middlewares

import (
	"compress/gzip"
	"io"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
    gin.ResponseWriter
    writer *gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
    return g.writer.Write(data)
}

var gzipPool = sync.Pool{
    New: func() interface{} {
        return gzip.NewWriter(io.Discard)
    },
}

func GzipMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
            c.Next()
            return
        }

        gz := gzipPool.Get().(*gzip.Writer)
        defer gzipPool.Put(gz)
        gz.Reset(c.Writer)

        c.Header("Content-Encoding", "gzip")
        c.Header("Vary", "Accept-Encoding")
        
        c.Writer = &gzipWriter{
            ResponseWriter: c.Writer,
            writer:         gz,
        }
        
        defer gz.Close()
        c.Next()
    }
}