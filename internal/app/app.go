package app

import (
	"context"

	"media-app/internal/handler"
	"media-app/internal/server"
)

// App struct represents the main application
type App struct {
	ctx context.Context

	MediaHandler *handler.MediaHandler
	HTTPServer   *server.HTTPServer
}

// New creates a new App application struct
func New(filePort int) *App {
	mediaHandler := handler.NewMediaHandler(filePort)
	httpServer := server.NewHTTPServer(filePort, mediaHandler)
	return &App{
		HTTPServer:   httpServer,
		MediaHandler: mediaHandler,
	}
}

// Startup is called when the app starts
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.MediaHandler.SetContext(ctx)
	a.HTTPServer.Start()
}

// Context returns the application context
func (a *App) Context() context.Context {
	return a.ctx
}

