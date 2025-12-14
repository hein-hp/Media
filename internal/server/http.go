package server

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"media-app/internal/handler"
	"media-app/pkg/logger"

	"go.uber.org/zap"
)

// HTTPServer represents the HTTP file server
type HTTPServer struct {
	port         int
	httpServer   *http.Server
	mutex        sync.Mutex
	mediaHandler *handler.MediaHandler
}

// NewHTTPServer creates a new HTTPServer instance
func NewHTTPServer(port int, mediaHandler *handler.MediaHandler) *HTTPServer {
	return &HTTPServer{
		port:         port,
		mediaHandler: mediaHandler,
	}
}

// Start starts the HTTP server
func (hs *HTTPServer) Start() {
	hs.Stop()

	fileHandler := hs.fileHandler()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", hs.port))
	if err != nil {
		logger.L().Error("HTTP server listen error", zap.Error(err))
		return
	}

	server := &http.Server{
		Handler: fileHandler,
	}

	hs.mutex.Lock()
	hs.httpServer = server
	hs.mutex.Unlock()

	logger.L().Info("HTTP server started", zap.Int("port", hs.port))

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			logger.L().Error("HTTP server error", zap.Error(err))
		}
	}()
}

// Stop stops the HTTP server
func (hs *HTTPServer) Stop() {
	hs.mutex.Lock()
	defer hs.mutex.Unlock()

	if hs.httpServer != nil {
		logger.L().Info("HTTP server stopping...")
		if err := hs.httpServer.Close(); err != nil {
			logger.L().Error("HTTP server stop error", zap.Error(err))
		}
		hs.httpServer = nil
	}
}

// GetPort returns the HTTP server port
func (hs *HTTPServer) GetPort() int {
	return hs.port
}

func (hs *HTTPServer) fileHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		dir := hs.mediaHandler.GetSelectedDir()
		if dir == "" {
			http.Error(w, "未选择文件目录", http.StatusForbidden)
			return
		}

		uri := strings.TrimPrefix(r.URL.Path, "/")
		if uri == "" {
			http.Error(w, "无效的文件路径", http.StatusBadRequest)
			return
		}
		http.ServeFile(w, r, filepath.Join(dir, uri))
	})
}
