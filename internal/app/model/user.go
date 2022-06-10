package model

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10/orm"
	"github.com/xuandoio/klik-dokter/internal/app/common"
)

type User struct {
	tableName struct{}        `pg:"users,alias:user"`
	ID        int             `pg:"id,pk"`
	Email     string          `pg:"email"`
	Password  string          `pg:"password"`
	IsActive  bool            `pg:"is_active,default:false"`
	CreatedAt common.DateTime `pg:"created_at"`
	UpdatedAt common.DateTime `pg:"updated_at"`
}

func FindUserByEmail(c *gin.Context, email string, db orm.DB) (User, error) {
	user := User{
		Email: email,
	}

	err := db.Model(&user).Where("email=?", email).Select()
	return user, err
}

func (user *User) Create(c *gin.Context, db orm.DB) error {
	// hash password using bcrypt before saving
	if hashedPassword, err := HashPassword(user.Password); err != nil {
		return err
	} else {
		user.Password = hashedPassword
	}

	user.IsActive = true // default is true when user register
	user.CreatedAt = common.DateTime{Time: time.Now().UTC()}
	user.UpdatedAt = common.DateTime{Time: time.Now().UTC()}
	if _, err := db.Model(user).Returning("*").Insert(); err != nil {
		return err
	}
	return nil
}
