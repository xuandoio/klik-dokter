package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xuandoio/klik-dokter/internal/app/handler"
	"github.com/xuandoio/klik-dokter/internal/app/middleware"
	"github.com/xuandoio/klik-dokter/internal/app/validator"
	"github.com/xuandoio/klik-dokter/internal/config"
)

func NewRouter(c *config.Config, h *handler.Handler) *gin.Engine {
	engine := gin.New()
	if err := validator.Register(); err != nil {
		log.Panicln(err)
	}
	engine.Use(
		middleware.Compressing(),            // compressing
		middleware.LoggingMiddleware(),      // std out logger
		middleware.RecoveryMiddleware(),     // recovery
		middleware.RateLimitingMiddleware(), // rate limiting
		middleware.ConfigMiddleware(c),      // attach config to gin context
	)

	engine = SetApiRoutes(engine, c, h)
	return engine
}
