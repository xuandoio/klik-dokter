package transformer

import (
	"github.com/xuandoio/klik-dokter/internal/app/common"
	"github.com/xuandoio/klik-dokter/internal/app/model"
)

type UserTransformer struct {
	ID        int             `json:"id"`
	Email     string          `json:"email"`
	IsActive  bool            `json:"is_active"`
	CreateAt  common.DateTime `json:"create_at"`
	UpdatedAt common.DateTime `json:"updated_at"`
}

func (user *UserTransformer) Transform(e interface{}) interface{} {
	userModel, ok := e.(model.User)
	if !ok {
		return e
	}

	user.ID = userModel.ID
	user.Email = userModel.Email
	user.IsActive = userModel.IsActive
	user.CreateAt = userModel.CreatedAt
	user.UpdatedAt = userModel.UpdatedAt
	return *user
}
