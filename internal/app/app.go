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
	similarHandler := handler.NewSimilarHandler()
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
func (a *App) RemoveMedia(mediaInfo handler.MediaInfo) {
	logger.Info("删除文件", zap.Any("mediaInfo", mediaInfo))
	a.MediaHandler.RemoveMedia(mediaInfo)
}
