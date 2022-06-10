package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuandoio/klik-dokter/internal/app/common"
	"github.com/xuandoio/klik-dokter/internal/app/model"
	"github.com/xuandoio/klik-dokter/internal/app/render"
	"github.com/xuandoio/klik-dokter/internal/app/request"
	"github.com/xuandoio/klik-dokter/internal/app/transformer"
)

func (h *Handler) ProductIndex(c *gin.Context) {
	// prepare data response
	products, err := model.FindProducts(c, h.db)
	if err != nil {
		render.Error(c, err)
		return
	}

	// transform result to Collection for generation JSON
	var productResult []interface{}
	for _, v := range products {
		productResult = append(productResult, v)
	}

	// render Collection to JSON
	productTransformer := &transformer.ProductTransformer{}
	productCollection := transformer.NewCollection(productResult, productTransformer)
	render.JSON(c, productCollection)
}

func (h *Handler) ProductShow(c *gin.Context) {
	pId, ok := c.Params.Get("id")
	if !ok {
		common.PanicNotFound()
	}
	productId, err := strconv.Atoi(pId)
	if err != nil {
		common.PanicNotFound()
	}

	// prepare data response
	product, err := model.FindProductByID(c, productId, h.db)
	if err != nil {
		common.PanicInternalServerError(err)
	}

	// render product item as JSON
	productTransformer := &transformer.ProductTransformer{}
	productItem := transformer.NewItem(product, productTransformer)
	render.JSON(c, productItem)
}

func (h *Handler) ProductCreate(c *gin.Context) {
	productRequest := request.Product{}
	if err := c.ShouldBind(&productRequest); err != nil {
		render.Error(c, err)
		return
	}

	// product entity
	product := model.Product{
		Name:     productRequest.Name,
		SKU:      productRequest.SKU,
		Quantity: productRequest.Quantity,
		Price:    productRequest.Price,
		Unit:     productRequest.Unit,
		Status:   productRequest.Status,
	}
	if err := product.Create(c, h.db); err != nil {
		render.Error(c, err)
		return
	}

	// render product item as JSON
	productTransformer := &transformer.ProductTransformer{}
	productItem := transformer.NewItem(product, productTransformer)
	render.JSON(c, productItem)
}

func (h *Handler) ProductUpdate(c *gin.Context) {
	pId, ok := c.Params.Get("id")
	if !ok {
		common.PanicNotFound()
	}
	productId, err := strconv.Atoi(pId)
	if err != nil {
		common.PanicNotFound()
	}

	//parse request
	productRequest := request.Product{}
	if err = c.ShouldBindJSON(&productRequest); err != nil {
		render.Error(c, err)
		return
	}

	// product entity
	product := model.Product{
		ID:       productId,
		Name:     productRequest.Name,
		SKU:      productRequest.SKU,
		Quantity: productRequest.Quantity,
		Price:    productRequest.Price,
		Unit:     productRequest.Unit,
		Status:   productRequest.Status,
	}

	if err = product.Update(c, h.db); err != nil {
		render.Error(c, err)
		return
	}

	// render product item as JSON
	productTransformer := &transformer.ProductTransformer{}
	productItem := transformer.NewItem(product, productTransformer)
	render.JSON(c, productItem)
}

func (h *Handler) ProductDestroy(c *gin.Context) {
	pId, ok := c.Params.Get("id")
	if !ok {
		common.PanicNotFound()
	}
	productId, err := strconv.Atoi(pId)
	if err != nil {
		common.PanicNotFound()
	}

	// product entity
	product := model.Product{
		ID: productId,
	}

	err = product.Delete(c, h.db)
	if err != nil {
		render.Error(c, err)
		return
	}
	render.NoContent(c)
}
