package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"personal-cockpit/database"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Criar e testar banco de dados
	fmt.Println("🔄 Inicializando banco de dados...")
	db, err := database.NewDB()
	if err != nil {
		log.Fatal("❌ Erro ao criar banco:", err)
	}
	defer db.Close()
	fmt.Println("✅ Banco de dados pronto!\n")

	// Criar instância do App
	app := NewApp()

	// Criar aplicação Wails
	err = wails.Run(&options.App{
		Title:  "Personal Cockpit",
		Width:  1024,
		Height: 768,
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
		println("Error:", err.Error())
	}
}
