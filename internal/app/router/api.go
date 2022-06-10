package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xuandoio/klik-dokter/internal/app/handler"
	"github.com/xuandoio/klik-dokter/internal/app/middleware"
	"github.com/xuandoio/klik-dokter/internal/config"
)

func SetApiRoutes(e *gin.Engine, c *config.Config, h *handler.Handler) *gin.Engine {
	apiRoutes := e.Group("/api")

	//authentication
	apiRoutes.POST("/register", h.UserRegister)
	apiRoutes.POST("/login", h.UserLogin)

	//products
	apiRoutes.Use(middleware.AuthenticateMiddleware(c, h.GetDB()))
	{
		apiRoutes.GET("/products", h.ProductIndex)
		apiRoutes.GET("/products/:id", h.ProductShow)
		apiRoutes.DELETE("/products/:id", h.ProductDestroy)
		apiRoutes.POST("/products", h.ProductCreate)
		apiRoutes.PUT("/products/:id", h.ProductUpdate)
	}
	return e
}
