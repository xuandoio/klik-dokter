package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xuandoio/klik-dokter/internal/app/render"
)

// RecoveryMiddleware /**
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		render.Error(c, recovered)
		return
	})
}
