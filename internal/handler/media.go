package handler

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"media-app/pkg/logger"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
)

// MediaHandler handles media file operations
type MediaHandler struct {
	dir  string
	mux  sync.Mutex
	ctx  context.Context
	port int
}

// MediaType represents the type of media
type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

// MediaInfo represents information about a media file
type MediaInfo struct {
	Path string    `json:"path"`
	Name string    `json:"name"`
	Size int64     `json:"size"`
	Url  string    `json:"url"`
	Type MediaType `json:"type"`
}

// NewMediaHandler creates a new MediaHandler instance
func NewMediaHandler(port int) *MediaHandler {
	return &MediaHandler{
		port: port,
	}
}

// SetContext sets the wails context
func (mh *MediaHandler) SetContext(ctx context.Context) {
	mh.ctx = ctx
}

// GetSelectedDir returns the selected directory
func (mh *MediaHandler) GetSelectedDir() string {
	return mh.dir
}

// SetSelectedDir sets the selected directory
func (mh *MediaHandler) SetSelectedDir(dir string) {
	mh.mux.Lock()
	defer mh.mux.Unlock()
	mh.dir = dir
	logger.L().Info("目录已选择", zap.String("dir", dir))
}

// GetMediaFiles scans the given directory and returns a list of media file paths
func (mh *MediaHandler) GetMediaFiles() []MediaInfo {
	var medias []MediaInfo

	dirStat, err := os.Stat(mh.dir)
	if err != nil || !dirStat.IsDir() {
		logger.L().Warn("无效的目录", zap.String("dir", mh.dir), zap.Error(err))
		return medias
	}

	err = filepath.WalkDir(mh.dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.L().Warn("遍历文件失败", zap.String("path", path), zap.Error(err))
			return nil
		}
		if d.IsDir() {
			return nil
		}

		abs, err := filepath.Abs(path)
		if err != nil {
			logger.L().Warn("获取绝对路径失败", zap.String("path", path), zap.Error(err))
			return nil
		}
		fileType, err := detectMediaType(abs)
		if err != nil {
			logger.L().Debug("识别媒体类型失败", zap.String("path", abs), zap.Error(err))
			return nil
		}

		var mediaType MediaType
		if strings.HasPrefix(fileType, "image/") {
			mediaType = MediaTypeImage
		} else if strings.HasPrefix(fileType, "video/") {
			mediaType = MediaTypeVideo
		} else {
			// 不是图片或视频，跳过
			return nil
		}

		fileinfo, err := d.Info()
		if err != nil {
			logger.L().Warn("获取文件信息失败", zap.String("path", abs), zap.Error(err))
			return nil
		}
		relPath, err := filepath.Rel(mh.dir, abs)
		if err != nil {
			logger.L().Warn("获取相对路径失败", zap.String("path", abs), zap.Error(err))
			return nil
		}
		urlPath := strings.ReplaceAll(relPath, string(filepath.Separator), "/")
		medias = append(medias, MediaInfo{
			Path: abs,
			Name: d.Name(),
			Size: fileinfo.Size(),
			Url:  fmt.Sprintf("http://localhost:%d/%s", mh.port, urlPath),
			Type: mediaType,
		})
		return nil
	})

	if err != nil {
		logger.L().Error("遍历目录失败", zap.String("dir", mh.dir), zap.Error(err))
	}

	logger.L().Info("扫描完成", zap.Int("count", len(medias)), zap.String("dir", mh.dir))
	return medias
}

// SendMediaFiles sends the media files to the frontend
func (mh *MediaHandler) SendMediaFiles(mediaInfos []MediaInfo) {
	runtime.EventsEmit(mh.ctx, "media-list", mediaInfos)
	logger.L().Debug("已发送媒体列表到前端", zap.Int("count", len(mediaInfos)))
}

// detectMediaType returns the MIME type of the given file
func detectMediaType(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer[:n])
	return contentType, nil
}

