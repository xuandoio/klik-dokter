package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/xuandoio/klik-dokter/internal/app/common"
	"github.com/xuandoio/klik-dokter/internal/app/middleware"
	"github.com/xuandoio/klik-dokter/internal/app/model"
	"github.com/xuandoio/klik-dokter/internal/app/render"
	"github.com/xuandoio/klik-dokter/internal/app/request"
	"github.com/xuandoio/klik-dokter/internal/app/transformer"
)

func (h *Handler) UserRegister(c *gin.Context) {
	userRegister := request.User{}
	if err := c.ShouldBind(&userRegister); err != nil {
		render.Error(c, err)
		return
	}

	user := model.User{
		Email:    userRegister.Email,
		Password: userRegister.Password,
	}

	exist, err := model.FindUserByEmail(c, user.Email, h.db)
	if err == nil && exist.ID > 0 {
		render.Error(c, errors.New("Email already exists"))
		return
	}

	if err := user.Create(c, h.db); err != nil {
		render.Error(c, err)
		return
	}

	// render user item as JSON
	userTransformer := &transformer.UserTransformer{}
	userItem := transformer.NewItem(user, userTransformer)
	render.JSON(c, userItem)
}

func (h *Handler) UserLogin(c *gin.Context) {
	userRegister := request.User{}
	if err := c.ShouldBind(&userRegister); err != nil {
		render.Error(c, err)
		return
	}

	// find user by email in database
	currentUser, err := model.FindUserByEmail(c, userRegister.Email, h.db)
	if err != nil {
		render.Error(c, err)
		return
	}

	// compare password
	if !model.CheckPasswordHash(userRegister.Password, currentUser.Password) {
		common.PanicUnauthorized()
		return
	}

	// expired at
	expirationTime := time.Now().Add(1440 * time.Minute)
	claims := &middleware.Payload{
		Email: &currentUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		common.PanicInternalServerError(err)
		return
	}

	// return token
	tokenModel := model.Token{
		AccessToken: tokenString,
		ExpiredAt:   expirationTime,
	}
	tokenTransformer := &transformer.TokenTransformer{}
	tokenItem := transformer.NewItem(tokenModel, tokenTransformer)
	render.JSON(c, tokenItem)
}
