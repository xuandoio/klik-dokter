package render

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xuandoio/klik-dokter/internal/app/common"
	"github.com/xuandoio/klik-dokter/internal/app/transformer"
	validation "github.com/xuandoio/klik-dokter/internal/app/validator"
)

// Error /**
func Error(c *gin.Context, payload interface{}) {

	var statusCode = http.StatusOK
	var data interface{}
	jsonSerializer := transformer.NewJSONSerializer()

	switch e := payload.(type) {
	case validator.ValidationErrors: // bad request
		statusCode = http.StatusBadRequest
		data = validation.Translate(e)
		jsonSerializer.Messages = []string{http.StatusText(http.StatusNotFound)}
		break
	case common.Error: // common error
		statusCode = http.StatusBadRequest
		jsonSerializer.Messages = []string{http.StatusText(e.StatusCode)}
		break
	case error: // internal error
		statusCode = http.StatusInternalServerError
		jsonSerializer.Messages = []string{e.Error()}
		break
	default: // default
		statusCode = http.StatusInternalServerError
		jsonSerializer.Messages = []string{http.StatusText(http.StatusInternalServerError)}
	}

	transformerManager := transformer.Manager{
		Serializer: jsonSerializer,
	}
	errorItem := transformer.NewError(data)

	payloadStruct := transformerManager.CreateData(errorItem)
	c.JSON(statusCode, payloadStruct)
}
