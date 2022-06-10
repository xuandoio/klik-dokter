package handler

import (
	"context"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/xuandoio/klik-dokter/internal/app/model"
	"github.com/xuandoio/klik-dokter/internal/config"
)

type Handler struct {
	db     *pg.DB
	ctx    context.Context
	config *config.Config
}

func (h *Handler) GetDB() *pg.DB {
	return h.db
}

func New(ctx context.Context, c *config.Config) *Handler {

	connection, err := model.NewConnection(c)
	if err != nil {
		log.Fatalf("Error connecting to database...")
	}

	return &Handler{
		db:     connection,
		ctx:    ctx,
		config: c,
	}
}
