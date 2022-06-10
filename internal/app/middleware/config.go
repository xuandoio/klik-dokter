package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xuandoio/klik-dokter/internal/config"
)

func ConfigMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", config)
		c.Next()
	}
}
