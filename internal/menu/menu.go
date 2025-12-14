package menu

import (
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
	fileMenu.AddText("选择文件夹", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
		filepath, err := wailsruntime.OpenDirectoryDialog(app.Context(), wailsruntime.OpenDialogOptions{})
		if err != nil {
			logger.L().Error("打开目录对话框失败", zap.Error(err))
			return
		}
		if filepath == "" {
			logger.L().Debug("用户取消选择目录")
			return
		}
		app.MediaHandler.SetSelectedDir(filepath)
		mediaInfos := app.MediaHandler.GetMediaFiles()
		if len(mediaInfos) != 0 {
			app.MediaHandler.SendMediaFiles(mediaInfos)
		}
	})
	fileMenu.AddSeparator()
	return appMenu
}
