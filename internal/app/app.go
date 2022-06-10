package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuandoio/klik-dokter/internal/app/handler"
	"github.com/xuandoio/klik-dokter/internal/app/router"
	"github.com/xuandoio/klik-dokter/internal/config"
	"github.com/xuandoio/klik-dokter/internal/db/migration"
)

type App struct {
	router  *gin.Engine
	handler *handler.Handler
	ctx     context.Context
	server  *http.Server
	config  *config.Config
}

func (app *App) Run() {

	err := app.Migrate()
	if err != nil {
		log.Fatalf("Migrating fail:%s\n", err)
		return
	}

	app.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port),
		Handler:      app.router,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	// Start the reminder schedule
	//app.reminder.Run()

	// Start the
	log.Printf("Listening on %s:%v...\n", app.config.Server.Host, app.config.Server.Port)
	log.Fatalln(app.server.ListenAndServe())
}

func NewApp(config *config.Config) *App {
	if config.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	c := context.Background()
	h := handler.New(c, config)
	r := router.NewRouter(config, h)

	return &App{
		router:  r,
		handler: h,
		ctx:     c,
		config:  config,
	}
}

// Migrate /**
func (app *App) Migrate() (err error) {
	migrateEngine := migration.NewEngine(app.config)
	return migrateEngine.Migrate()
}
