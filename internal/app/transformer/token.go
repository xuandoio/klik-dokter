package transformer

import (
	"time"

	"github.com/xuandoio/klik-dokter/internal/app/model"
)

type TokenTransformer struct {
	AccessToken string    `json:"access_token"`
	ExpiredAt   time.Time `json:"expired_at"`
}

// Transform /**
func (token *TokenTransformer) Transform(e interface{}) interface{} {
	tokenModel, ok := e.(model.Token)
	if !ok {
		return e
	}

	token.AccessToken = tokenModel.AccessToken
	token.ExpiredAt = tokenModel.ExpiredAt
	return *token
}
