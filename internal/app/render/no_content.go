package render

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{})
}
