package main

import (
	"embed"
	"log"
	"os"

	"media-app/internal/app"
	"media-app/internal/menu"
	"media-app/pkg/logger"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"go.uber.org/zap"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("获取用户主目录失败，%v", err)
	}
	logger.InitWithFile(home)
	defer logger.Sync()
	logger.L().Info("应用启动中...")

	app := app.New(8080)

	err = wails.Run(&options.App{
		Title:  "Media",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: options.NewRGBA(255, 255, 255, 0),
		Menu:             menu.Create(app),
		OnStartup:        app.Startup,
		Bind: []any{
			app,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				HideTitle: true,
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
	})

	if err != nil {
		logger.L().Fatal("应用运行失败", zap.Error(err))
	}
}
