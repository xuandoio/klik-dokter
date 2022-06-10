package main

import (
	"github.com/xuandoio/klik-dokter/internal/app"
	"github.com/xuandoio/klik-dokter/internal/app/common"
)

func main() {
	config := common.ReadConfig()
	application := app.NewApp(config)
	application.Run()
}
