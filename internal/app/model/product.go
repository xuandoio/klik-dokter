package model

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10/orm"
	"github.com/xuandoio/klik-dokter/internal/app/common"
)

var ErrForbiddenResource = errors.New("forbidden resource")

type Product struct {
	tableName struct{}        `pg:"products,alias:product"`
	ID        int             `pg:"id,pk"`
	SKU       string          `pg:"sku"`
	Name      string          `pg:"name"`
	Quantity  uint64          `pg:"quantity"`
	Price     float64         `pg:"price"`
	Unit      string          `pg:"unit"`
	Status    int             `pg:"status"`
	CreatedBy int             `pg:"created_by"`
	CreatedAt common.DateTime `pg:"created_at"`
	UpdatedAt common.DateTime `pg:"updated_at"`
}

func FindProducts(c *gin.Context, db orm.DB) ([]Product, error) {
	products := make([]Product, 0)
	err := db.Model(&products).Select()
	return products, err
}

func (product *Product) Create(c *gin.Context, db orm.DB) error {
	u := c.MustGet("user")
	user, ok := u.(*User)
	if !ok {
		panic("user object not found from interface of context")
	}

	product.CreatedBy = user.ID
	product.CreatedAt = common.DateTime{Time: time.Now().UTC()}
	product.UpdatedAt = common.DateTime{Time: time.Now().UTC()}

	_, err := db.Model(product).Insert()
	return err
}

func (product *Product) Update(c *gin.Context, db orm.DB) error {
	product.UpdatedAt = common.DateTime{Time: time.Now().UTC()}

	_, err := db.Model(product).
		Where("id=?", product.ID).
		ExcludeColumn("created_by", "created_at").
		Update()
	return err
}

func (product *Product) Delete(c *gin.Context, db orm.DB) error {
	_, err := db.Model(product).Where("id=?", product.ID).Delete()
	return err
}

func FindProductByID(c *gin.Context, id int, db orm.DB) (Product, error) {
	product := Product{
		ID: id,
	}

	err := db.Model(&product).WherePK().First()
	if err != nil {
		return product, err
	}

	return product, err
}

func FindProductBySKU(c *gin.Context, sku string, db orm.DB) (Product, error) {
	product := Product{
		SKU: sku,
	}

	err := db.Model(&product).Where("sku=?", product.SKU).Select()
	return product, err
}
