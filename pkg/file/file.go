package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MediaType 定义文件类型枚举
type MediaType string

const (
	MediaTypeVideo   MediaType = "video"
	MediaTypeImage   MediaType = "image"
	MediaTypeDoc     MediaType = "docx"
	MediaTypeAudio   MediaType = "audio"
	MediaTypeUnknown MediaType = "unknown"
)

// ExtToFileType 后缀名 -> 文件类型
var ExtToFileType = map[string]MediaType{
	// 视频后缀
	".mp4":  MediaTypeVideo,
	".mov":  MediaTypeVideo,
	".avi":  MediaTypeVideo,
	".mkv":  MediaTypeVideo,
	".flv":  MediaTypeVideo,
	".wmv":  MediaTypeVideo,
	".webm": MediaTypeVideo,
	// 图片后缀
	".jpg":  MediaTypeImage,
	".jpeg": MediaTypeImage,
	".png":  MediaTypeImage,
	".gif":  MediaTypeImage,
	".bmp":  MediaTypeImage,
	".webp": MediaTypeImage,
	".tiff": MediaTypeImage,
	// 文档后缀
	".txt":  MediaTypeDoc,
	".pdf":  MediaTypeDoc,
	".docx": MediaTypeDoc,
	".xlsx": MediaTypeDoc,
	".pptx": MediaTypeDoc,
	// 音频后缀
	".mp3":  MediaTypeAudio,
	".wav":  MediaTypeAudio,
	".flac": MediaTypeAudio,
	".aac":  MediaTypeAudio,
}

// filterFiles 需过滤的隐藏文件/目录
var filterFiles = map[string]bool{
	".DS_Store": true,
	".Trash":    true,
	"Thumbs.db": true,
	".deleted":  true,
}

// Meta 文件元数据
type Meta struct {
	FullPath string
	FileName string
	Ext      string
	ModTime  time.Time
}

// GetFileMetas 获取dir下的文件元数据信息
func GetFileMetas(dir string) ([]Meta, error) {
	// 校验目标目录有效性
	dirStat, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("目录不存在或无法访问：%w", err)
	}
	if !dirStat.IsDir() {
		return nil, errors.New("指定路径不是文件夹")
	}

	// 读取目录下所有文件，收集文件信息
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败：%w", err)
	}

	var metaList []Meta

	for _, entry := range entries {
		// 跳过目录
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		// 过滤隐藏文件/系统无用文件
		if filterFiles[fileName] {
			continue
		}

		fullPath := filepath.Join(dir, fileName)
		fileStat, err := os.Stat(fullPath)
		if err != nil {
			return nil, fmt.Errorf("获取文件信息失败 %s %w", fileName, err)
		}

		// 提取文件后缀
		ext := filepath.Ext(fileName)
		// 存储文件信息
		metaList = append(metaList, Meta{
			FullPath: fullPath,
			FileName: fileName,
			Ext:      ext,
			ModTime:  fileStat.ModTime(),
		})
	}

	return metaList, nil
}

// GetFileMeta 获取path文件路径的元数据信息
func GetFileMeta(path string) (*Meta, error) {
	fileStat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败：%s %w", path, err)
	}
	if fileStat.IsDir() {
		return nil, errors.New("指定路径是文件夹")
	}

	// 提取文件后缀
	ext := filepath.Ext(path)
	base := filepath.Base(path)
	// 存储文件信息
	return &Meta{
		FullPath: path,
		FileName: base,
		Ext:      ext,
		ModTime:  fileStat.ModTime(),
	}, nil
}

// CountFiles 对dir文件夹下的文件计数，不递归
func CountFiles(dir string) (int, error) {
	if dir == "" {
		return -1, fmt.Errorf("文件件为空")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return -1, fmt.Errorf("读取目录失败: %v", err)
	}

	fileCount := 0
	for _, entry := range entries {
		if !entry.IsDir() && !FilterFile(entry.Name()) {
			fileCount++
		}
	}

	return fileCount, nil
}

// GetFileTypeByExt 通过文件后缀名获取文件类型
func GetFileTypeByExt(filePath string) MediaType {
	ext := strings.ToLower(filepath.Ext(filePath))
	if fileType, ok := ExtToFileType[ext]; ok {
		return fileType
	}
	return MediaTypeUnknown
}

// RenameFile 重命名文件
func RenameFile(oldPath, newPath string, suffix bool, maxTry int) error {
	// 校验原文件是否存在
	oldFileInfo, err := os.Stat(oldPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("原文件不存在：%s", oldPath)
		}
		return fmt.Errorf("获取原文件信息失败：%w", err)
	}
	// 校验原路径是否为文件
	if oldFileInfo.IsDir() {
		return fmt.Errorf("原路径是目录，不支持重命名：%s", oldPath)
	}

	// 确定最终目标路径
	finalNewPath, err := getFinalTargetPath(newPath, suffix, maxTry)
	if err != nil {
		return fmt.Errorf("生成最终目标路径失败：%w", err)
	}

	// 执行文件重命名
	if err := os.Rename(oldPath, finalNewPath); err != nil {
		return fmt.Errorf("文件重命名失败：%w", err)
	}

	return nil
}

// FilterFile 是否过滤文件
func FilterFile(fileName string) bool {
	return filterFiles[fileName]
}

// getFinalTargetPath 生成最终目标路径（处理递增数字后缀逻辑）
func getFinalTargetPath(baseNewPath string, suffix bool, maxTry int) (string, error) {
	// 不开启重复后缀，直接校验并返回基础路径
	if !suffix {
		_, err := os.Stat(baseNewPath)
		if err == nil {
			return "", fmt.Errorf("目标文件已存在：%s", baseNewPath)
		}
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("获取目标路径信息失败：%w", err)
		}
		return baseNewPath, nil
	}

	// 开启重复后缀，生成_1、_2递增路径
	// 解析路径组件：目录、完整文件名、后缀、纯文件名
	dir := filepath.Dir(baseNewPath)
	baseName := filepath.Base(baseNewPath)
	ext := filepath.Ext(baseName)
	fileNameWithoutExt := strings.TrimSuffix(baseName, ext)

	// 循环尝试递增后缀路径，直到找到不存在的路径
	tryCount := 1
	currentTargetPath := baseNewPath
	for {
		// 检查当前路径是否存在
		_, err := os.Stat(currentTargetPath)
		if os.IsNotExist(err) {
			// 路径不存在，即为最终目标路径
			return currentTargetPath, nil
		}
		if err != nil {
			return "", fmt.Errorf("检查路径 %s 失败：%w", currentTargetPath, err)
		}

		// 路径已存在，拼接_+数字后缀生成新路径
		newFileName := fmt.Sprintf("%s_%d%s", fileNameWithoutExt, tryCount, ext)
		currentTargetPath = filepath.Join(dir, newFileName)
		tryCount++

		// 添加最大尝试次数限制，防止无限循环
		if tryCount > maxTry {
			return "", fmt.Errorf("超出最大尝试次数，未找到可用路径")
		}
	}
}
