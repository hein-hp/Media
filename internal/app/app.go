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

	MediaHandler   *handler.MediaHandler
	SimilarHandler *handler.SimilarHandler
	HttpServer     *server.HttpServer
}

// New creates a new App application struct
func New(filePort int) *App {
	mediaHandler := handler.NewMediaHandler(filePort)
	similarHandler := handler.NewSimilarHandler(filePort)
	httpServer := server.NewHttpServer(filePort, mediaHandler)
	return &App{
		HttpServer:     httpServer,
		MediaHandler:   mediaHandler,
		SimilarHandler: similarHandler,
	}
}

// Startup is called when the app starts
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.MediaHandler.SetContext(ctx)
	a.SimilarHandler.SetContext(ctx)
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
