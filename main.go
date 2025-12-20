package main

import (
	"embed"
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

//go:embed build/appicon.png
var icon []byte

func main() {
	err := logger.Init(logger.DefaultConfig())
	if err != nil {
		panic("初始化日志失败: " + err.Error())
	}
	defer logger.Sync()
	logger.Info("应用启动中...")

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
			About: &mac.AboutInfo{
				Title:   "Media",
				Message: "一个媒体工具",
				Icon:    icon,
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
	})

	if err != nil {
		logger.Fatal("应用运行失败", zap.Error(err))
	}
}
