package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"personal-cockpit/config"
	"personal-cockpit/database"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cfg := config.MustLoad()

	log.Println("Inicializando banco de dados...")
	db, err := database.NewDB()
	if err != nil {
		log.Fatal("Erro ao criar banco:", err)
	}
	defer db.Close()
	log.Println("Banco de dados pronto!")

	app := NewApp()

	err = wails.Run(&options.App{
		Title:  cfg.WindowTitle,
		Width:  cfg.WindowWidth,
		Height: cfg.WindowHeight,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal("Erro ao iniciar app:", err)
	}
}
