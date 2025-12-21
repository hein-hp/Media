package handler

import (
	"context"
	"fmt"
	"io/fs"
	"media-app/pkg/file"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

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

// MediaInfo represents information about a media file
type MediaInfo struct {
	Path    string         `json:"path"`
	Name    string         `json:"name"`
	Size    int64          `json:"size"`
	Url     string         `json:"url"`
	Type    file.MediaType `json:"type"`
	ModTime time.Time      `json:"modTime"`
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
	logger.Info("目录已选择", zap.String("dir", dir))
}

// GetMediaFiles scans the given directory and returns a list of media file paths
func (mh *MediaHandler) GetMediaFiles() []MediaInfo {
	var medias []MediaInfo

	dirStat, err := os.Stat(mh.dir)
	if err != nil || !dirStat.IsDir() {
		logger.Error("无效的目录", zap.String("dir", mh.dir), zap.Error(err))
		return medias
	}

	err = filepath.WalkDir(mh.dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Error("遍历文件失败", zap.String("path", path), zap.Error(err))
			return nil
		}
		if d.IsDir() {
			if path == mh.dir {
				return nil
			}
			return filepath.SkipDir
		}

		abs, err := filepath.Abs(path)
		if err != nil {
			logger.Error("获取绝对路径失败", zap.String("path", path), zap.Error(err))
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			logger.Error("获取文件信息失败", zap.String("path", abs), zap.Error(err))
			return nil
		}
		relPath, err := filepath.Rel(mh.dir, abs)
		if err != nil {
			logger.Error("获取相对路径失败", zap.String("path", abs), zap.Error(err))
			return nil
		}

		mediaType := file.GetFileTypeByExt(relPath)
		if mediaType != file.MediaTypeImage && mediaType != file.MediaTypeVideo {
			logger.Error("非图片视频", zap.String("abs", abs), zap.Any("mediaType", mediaType))
			return nil
		}
		urlPath := strings.ReplaceAll(relPath, string(filepath.Separator), "/")
		// 添加修改时间戳防止浏览器缓存
		modTimeUnix := fileInfo.ModTime().Unix()
		medias = append(medias, MediaInfo{
			Path:    abs,
			Name:    d.Name(),
			Size:    fileInfo.Size(),
			Url:     fmt.Sprintf("http://localhost:%d/%s?t=%d", mh.port, urlPath, modTimeUnix),
			Type:    file.GetFileTypeByExt(relPath),
			ModTime: fileInfo.ModTime(),
		})
		return nil
	})

	if err != nil {
		logger.Error("遍历目录失败", zap.String("dir", mh.dir), zap.Error(err))
	}

	logger.Info("扫描完成", zap.Int("count", len(medias)), zap.String("dir", mh.dir))

	// 最新的文件在前面
	sort.Slice(medias, func(i, j int) bool {
		return medias[i].ModTime.After(medias[j].ModTime)
	})
	return medias
}

// SendMediaFiles sends the media files to the frontend
func (mh *MediaHandler) SendMediaFiles(mediaInfos []MediaInfo, fileCount int, filepath string) {
	runtime.EventsEmit(mh.ctx, "media-list", mediaInfos)
	runtime.EventsEmit(mh.ctx, "media-count", fileCount)
	runtime.EventsEmit(mh.ctx, "selected-dir", filepath)
	logger.Debug("已发送媒体列表到前端", zap.Int("size", len(mediaInfos)), zap.Int("count", fileCount))
}

// FixMediaFilename fix the media filename
func (mh *MediaHandler) FixMediaFilename() {
	selected := mh.GetSelectedDir()
	err := file.WithOrderly(selected, 4)
	if err != nil {
		logger.Error("重排序文件失败", zap.Error(err))
		return
	}
	files, err := file.CountFiles(selected)
	if err != nil {
		logger.Error("读取文件数量失败", zap.Error(err))
		files = -1
	}
	mh.SendMediaFiles(mh.GetMediaFiles(), files, selected)
}

// BatchFixMediaFilename batch fix the media filename
func (mh *MediaHandler) BatchFixMediaFilename() {
	dir := mh.GetSelectedDir()
	dirStat, err := os.Stat(dir)
	if err != nil {
		logger.Error("目录不存在或无法访问", zap.String("dir", dir), zap.Error(err))
		return
	}
	if !dirStat.IsDir() {
		logger.Error("指定路径不是文件夹", zap.String("dir", dir))
		return
	}

	// 读取目录下所有文件，收集文件信息
	entries, err := os.ReadDir(dir)
	if err != nil {
		logger.Error("读取目录失败", zap.String("dir", dir), zap.Error(err))
		return
	}

	// 过滤出目录
	var dirs []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		}
	}

	logger.Info("开始批量修复文件名", zap.Int("目录数量", len(dirs)))

	sem := make(chan struct{}, 5)
	wg := &sync.WaitGroup{}
	wg.Add(len(dirs))

	for _, entry := range dirs {
		sem <- struct{}{}
		go func(e os.DirEntry) {
			defer wg.Done()
			defer func() { <-sem }()

			sub := filepath.Join(dir, e.Name())
			logger.Debug("处理子目录", zap.String("dir", sub))

			err := file.WithOrderly(sub, 4)
			if err != nil {
				logger.Error("修复文件名错误", zap.String("dir", sub), zap.Error(err))
			} else {
				logger.Debug("子目录处理完成", zap.String("dir", sub))
			}
		}(entry)
	}

	wg.Wait()
	logger.Info("批量修复文件名完成")
}

// RemoveMedia 删除媒体资源
func (mh *MediaHandler) RemoveMedia(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		logger.Error("源文件不存在或无法访问", zap.String("filePath", filePath), zap.Error(err))
		return err
	}
	srcDir := filepath.Dir(filePath)

	deleteDir := filepath.Join(srcDir, ".delete")
	if err := os.MkdirAll(deleteDir, 0755); err != nil {
		logger.Error("创建 .delete 目录失败", zap.String("filePath", filePath), zap.Error(err))
		return err
	}

	// 提取源文件名，拼接目标文件路径
	fileName := filepath.Base(filePath)
	targetFilePath := filepath.Join(deleteDir, fileName)

	err = file.RenameFile(filePath, targetFilePath, true, 100)
	if err != nil {
		logger.Error("移动到 .delete 目录失败", zap.String("filePath", filePath), zap.Error(err))
	}
	return nil
}
