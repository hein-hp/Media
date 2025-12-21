package menu

import (
	"media-app/pkg/file"
	"runtime"

	"media-app/internal/app"
	"media-app/pkg/logger"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
)

// Create creates the application menu
func Create(app *app.App) *menu.Menu {
	appMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.AppMenu())
	}

	fileMenu := appMenu.AddSubmenu("文件")
	fileMenu.AddText("选择文件夹", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) { openDirectory(app) })
	fileMenu.AddSeparator()

	operMenu := appMenu.AddSubmenu("操作")
	operMenu.AddText("修复文件名", &keys.Accelerator{}, func(_ *menu.CallbackData) { app.MediaHandler.FixMediaFilename() })
	operMenu.AddText("修复文件名（批量）", &keys.Accelerator{}, func(_ *menu.CallbackData) { app.MediaHandler.BatchFixMediaFilename() })
	operMenu.AddSeparator()
	operMenu.AddText("查找相同图片", keys.CmdOrCtrl("f"), func(_ *menu.CallbackData) { findSimilarImages(app) })
	operMenu.AddText("快捷分类", keys.CmdOrCtrl("k"), func(_ *menu.CallbackData) { openClassify(app) })

	return appMenu
}

func openDirectory(app *app.App) string {
	filepath, err := wailsruntime.OpenDirectoryDialog(app.Context(), wailsruntime.OpenDialogOptions{})
	if err != nil {
		logger.Error("打开目录对话框失败", zap.Error(err))
		return ""
	}
	if filepath == "" {
		logger.Debug("用户取消选择目录")
		return ""
	}
	app.MediaHandler.SetSelectedDir(filepath)
	app.SimilarHandler.SetSelectedDir(filepath)
	app.ShortcutHandler.SetSelectedDir(filepath)
	fileCount, err := file.CountFiles(filepath)
	if err != nil {
		logger.Error("读取文件失败", zap.Error(err))
		return ""
	}
	if fileCount != 0 {
		app.MediaHandler.SendMediaFiles(app.MediaHandler.GetMediaFiles(), fileCount, filepath)
	}
	// 跳转回首页
	Goto(app, "/")
	return filepath
}

// findSimilarImages 查找相似图片
func findSimilarImages(app *app.App) {
	// 发送加载状态
	wailsruntime.EventsEmit(app.Context(), "similar-loading", true)
	// 跳转到相似图片页面
	Goto(app, "/similar")

	// 异步执行相似度分析
	go func() {
		results := app.SimilarHandler.CalcSimilarity()
		app.SimilarHandler.SendSimilarResults(results)
	}()
}

func Goto(app *app.App, route string) {
	wailsruntime.EventsEmit(app.Context(), "router", route)
}

// openClassify 打开快捷分类页面
func openClassify(app *app.App) {
	Goto(app, "/classify")
}
