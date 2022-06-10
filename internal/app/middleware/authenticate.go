package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10/orm"
	"github.com/golang-jwt/jwt"
	"github.com/xuandoio/klik-dokter/internal/app/common"
	"github.com/xuandoio/klik-dokter/internal/app/model"
	"github.com/xuandoio/klik-dokter/internal/config"
)

func AuthenticateMiddleware(config *config.Config, db orm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		if len(bearerToken) > 0 {
			bearerToken = strings.TrimPrefix(bearerToken, "Bearer ")
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		payload := &Payload{}
		tkn, err := jwt.ParseWithClaims(bearerToken, payload, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}
			return []byte(config.JWTSecret), nil
		})

		if err != nil || tkn == nil || !tkn.Valid {
			log.Println(err)
			common.PanicUnauthorized()
			return
		}

		payload, ok := tkn.Claims.(*Payload)
		if !ok {
			common.PanicUnauthorized()
			return
		}

		user, err := model.FindUserByEmail(c, *payload.Email, db)
		if err == nil && user.Email != *payload.Email {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("user", &user)
		c.Next()
	}
}

type Payload struct {
	Subject int     `json:"sub,omitempty"`
	Prv     string  `json:"prv,omitempty"`
	Uid     string  `json:"uid,omitempty"`
	Name    *string `json:"name,omitempty"`
	Email   *string `json:"email,omitempty"`
	Avatar  *string `json:"avatar,omitempty"`
	jwt.StandardClaims
}
