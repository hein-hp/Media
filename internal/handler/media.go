package handler

import (
	"context"
	"errors"
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
			return nil
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
		medias = append(medias, MediaInfo{
			Path:    abs,
			Name:    d.Name(),
			Size:    fileInfo.Size(),
			Url:     fmt.Sprintf("http://localhost:%d/%s", mh.port, urlPath),
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
	err := fixMediaFilename(selected)
	if err != nil {
		logger.Error("修复资源名称失败", zap.Error(err))
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

	wg := &sync.WaitGroup{}
	wg.Add(len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		go func() {
			err := fixMediaFilename(filepath.Join(dir, entry.Name()))
			if err != nil {
				logger.Error("修复文件名错误", zap.String("dir", dir), zap.Error(err))
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

// fixMediaFilename 修复dir文件夹下的所有文件名称
func fixMediaFilename(dir string) error {
	dirStat, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("目录不存在或无法访问：%w", err)
	}
	if !dirStat.IsDir() {
		return errors.New("指定路径不是文件夹")
	}

	metas, err := file.GetFileMetas(dir)
	if err != nil {
		return err
	}

	// 修改时间逆序
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].ModTime.After(metas[j].ModTime)
	})

	lastIndex, err := file.CountFiles(dir)
	if err != nil {
		return err
	}
	successCount := 0
	for _, meta := range metas {
		numStr := fmt.Sprintf("%05d", lastIndex)
		newFileName := numStr + meta.Ext
		newFullPath := filepath.Join(dir, newFileName)

		if meta.FullPath == newFullPath {
			successCount++
		} else {
			logger.Debug("重命名文件", zap.String("old", meta.FullPath), zap.String("new", newFullPath))
			err := file.RenameFile(meta.FullPath, newFullPath, false, -1)
			if err != nil {
				return err
			}
			successCount++
			logger.Info("文件重命名成功",
				zap.String("oldName", meta.FileName), zap.String("newName", newFileName), zap.Time("modTime", meta.ModTime),
			)
		}
		lastIndex--
	}

	logger.Info("批量重命名完成", zap.Int("成功数量", successCount), zap.String("目录", dir))
	return nil
}
