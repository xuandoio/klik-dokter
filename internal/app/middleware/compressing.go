package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

type gzipResponseWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipResponseWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func Compressing() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldCompress(c) {
			c.Next()
			return
		}

		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestCompression)
		if err != nil {
			c.Next()
			return
		}

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")
		c.Writer = &gzipResponseWriter{c.Writer, gz}
		defer func() {
			c.Header("Content-Length", "0")
			gz.Close()
		}()

		c.Next()
	}
}

func shouldCompress(c *gin.Context) bool {
	if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") || c.Request.URL.Path == "/metrics" {
		return false
	}

	extension := filepath.Ext(c.Request.URL.Path)
	if len(extension) < 4 {
		return true
	}

	switch extension {
	case ".png", ".gif", ".jpeg", ".jpg":
		return false
	default:
		return true
	}
}
