package render

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuandoio/klik-dokter/internal/app/transformer"
)

// JSON /**
func JSON(c *gin.Context, payload interface{}) {
	transformerManager := transformer.NewManager()

	payloadStruct := transformerManager.CreateData(payload)
	c.JSON(http.StatusOK, payloadStruct)
}
