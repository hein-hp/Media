package app

import (
	"context"
	"media-app/pkg/logger"

	"media-app/internal/handler"
	"media-app/internal/server"

	"go.uber.org/zap"
)

// App struct represents the main application
type App struct {
	ctx context.Context

	MediaHandler    *handler.MediaHandler
	SimilarHandler  *handler.SimilarHandler
	ShortcutHandler *handler.ShortcutHandler
	HttpServer      *server.HttpServer
}

// New creates a new App application struct
func New(filePort int) *App {
	mediaHandler := handler.NewMediaHandler(filePort)
	similarHandler := handler.NewSimilarHandler(filePort)
	shortcutHandler := handler.NewShortcutHandler(filePort)
	httpServer := server.NewHttpServer(filePort, mediaHandler)
	return &App{
		HttpServer:      httpServer,
		MediaHandler:    mediaHandler,
		SimilarHandler:  similarHandler,
		ShortcutHandler: shortcutHandler,
	}
}

// Startup is called when the app starts
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.MediaHandler.SetContext(ctx)
	a.SimilarHandler.SetContext(ctx)
	a.ShortcutHandler.SetContext(ctx)
	a.HttpServer.Start()
}

// Context returns the application context
func (a *App) Context() context.Context {
	return a.ctx
}

// RemoveMedia 删除媒体资源
func (a *App) RemoveMedia(path string) error {
	logger.Info("删除文件", zap.String("path", path))
	return a.MediaHandler.RemoveMedia(path)
}

// RemoveSimilarImage 删除相似图片
func (a *App) RemoveSimilarImage(path string) error {
	logger.Info("删除相似图片", zap.String("path", path))
	return a.SimilarHandler.RemoveSimilarImage(path)
}

// ==================== 快捷键分类相关 ====================

// GetShortcuts 获取快捷键配置
func (a *App) GetShortcuts() []handler.ShortcutConfig {
	return a.ShortcutHandler.GetShortcuts()
}

// SaveShortcuts 保存快捷键配置
func (a *App) SaveShortcuts(shortcuts []handler.ShortcutConfig) error {
	return a.ShortcutHandler.SaveShortcuts(shortcuts)
}

// SelectShortcutTargetDir 选择快捷键目标目录
func (a *App) SelectShortcutTargetDir() string {
	return a.ShortcutHandler.SelectTargetDir()
}

// MoveByShortcut 通过快捷键移动文件
func (a *App) MoveByShortcut(filePath string, shortcutKey string) error {
	logger.Info("快捷键移动文件", zap.String("path", filePath), zap.String("key", shortcutKey))
	return a.ShortcutHandler.MoveByShortcut(filePath, shortcutKey)
}

// UndoMove 撤销上一次移动
func (a *App) UndoMove() error {
	logger.Info("撤销移动操作")
	return a.ShortcutHandler.Undo()
}

// GetUndoCount 获取可撤销操作数量
func (a *App) GetUndoCount() int {
	return a.ShortcutHandler.GetUndoCount()
}

// SetClassifyDir 设置分类目录
func (a *App) SetClassifyDir(dir string) {
	a.ShortcutHandler.SetSelectedDir(dir)
}

// GetClassifyDir 获取分类目录
func (a *App) GetClassifyDir() string {
	return a.ShortcutHandler.GetSelectedDir()
}
