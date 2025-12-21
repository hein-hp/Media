package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"media-app/pkg/file"
	"media-app/pkg/logger"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
)

// ShortcutConfig 快捷键配置
type ShortcutConfig struct {
	Key       string `json:"key"`       // 快捷键 (1-9, a-z)
	TargetDir string `json:"targetDir"` // 目标文件夹（绝对路径或相对路径）
	Label     string `json:"label"`     // 显示名称
}

// ShortcutsData 快捷键配置数据
type ShortcutsData struct {
	Shortcuts []ShortcutConfig `json:"shortcuts"`
}

// MoveRecord 移动记录（用于撤销）
type MoveRecord struct {
	SourcePath string    `json:"sourcePath"` // 原始路径
	TargetPath string    `json:"targetPath"` // 目标路径
	Timestamp  time.Time `json:"timestamp"`  // 操作时间
}

// ShortcutHandler 快捷键处理器
type ShortcutHandler struct {
	ctx         context.Context
	port        int
	dir         string // 当前工作目录
	mux         sync.Mutex
	configPath  string       // 配置文件路径
	undoStack   []MoveRecord // 撤销栈
	maxUndoSize int          // 最大撤销记录数
}

// NewShortcutHandler 创建快捷键处理器
func NewShortcutHandler(port int) *ShortcutHandler {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("获取用户目录失败", zap.Error(err))
		homeDir = "."
	}

	configDir := filepath.Join(homeDir, ".media-app")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		logger.Error("创建配置目录失败", zap.Error(err))
	}

	return &ShortcutHandler{
		port:        port,
		configPath:  filepath.Join(configDir, "shortcuts.json"),
		undoStack:   make([]MoveRecord, 0),
		maxUndoSize: 50,
	}
}

// SetContext 设置 wails 上下文
func (sh *ShortcutHandler) SetContext(ctx context.Context) {
	sh.ctx = ctx
}

// GetSelectedDir 获取当前选中的目录
func (sh *ShortcutHandler) GetSelectedDir() string {
	return sh.dir
}

// SetSelectedDir 设置当前选中的目录
func (sh *ShortcutHandler) SetSelectedDir(dir string) {
	sh.mux.Lock()
	defer sh.mux.Unlock()
	sh.dir = dir
	logger.Info("分类目录已选择", zap.String("dir", dir))
}

// GetShortcuts 获取所有快捷键配置
func (sh *ShortcutHandler) GetShortcuts() []ShortcutConfig {
	data, err := os.ReadFile(sh.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 返回默认配置
			return sh.getDefaultShortcuts()
		}
		logger.Error("读取快捷键配置失败", zap.Error(err))
		return sh.getDefaultShortcuts()
	}

	var config ShortcutsData
	if err := json.Unmarshal(data, &config); err != nil {
		logger.Error("解析快捷键配置失败", zap.Error(err))
		return sh.getDefaultShortcuts()
	}

	return config.Shortcuts
}

// getDefaultShortcuts 获取默认快捷键配置
func (sh *ShortcutHandler) getDefaultShortcuts() []ShortcutConfig {
	return []ShortcutConfig{
		{Key: "1", TargetDir: "", Label: "分类1"},
		{Key: "2", TargetDir: "", Label: "分类2"},
		{Key: "3", TargetDir: "", Label: "分类3"},
		{Key: "d", TargetDir: ".delete", Label: "待删除"},
		{Key: "s", TargetDir: ".star", Label: "收藏"},
	}
}

// SaveShortcuts 保存快捷键配置
func (sh *ShortcutHandler) SaveShortcuts(shortcuts []ShortcutConfig) error {
	sh.mux.Lock()
	defer sh.mux.Unlock()

	data := ShortcutsData{Shortcuts: shortcuts}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logger.Error("序列化快捷键配置失败", zap.Error(err))
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(sh.configPath, jsonData, 0644); err != nil {
		logger.Error("写入快捷键配置失败", zap.Error(err))
		return fmt.Errorf("写入配置失败: %w", err)
	}

	logger.Info("快捷键配置已保存", zap.Int("count", len(shortcuts)))
	return nil
}

// SelectTargetDir 弹出系统文件夹选择器
func (sh *ShortcutHandler) SelectTargetDir() string {
	dir, err := runtime.OpenDirectoryDialog(sh.ctx, runtime.OpenDialogOptions{
		Title: "选择目标文件夹",
	})
	if err != nil {
		logger.Error("选择目录失败", zap.Error(err))
		return ""
	}
	return dir
}

// MoveByShortcut 通过快捷键移动文件
func (sh *ShortcutHandler) MoveByShortcut(filePath string, shortcutKey string) error {
	sh.mux.Lock()
	defer sh.mux.Unlock()

	// 查找快捷键配置
	shortcuts := sh.GetShortcuts()
	var targetConfig *ShortcutConfig
	for _, sc := range shortcuts {
		if strings.EqualFold(sc.Key, shortcutKey) {
			targetConfig = &sc
			break
		}
	}

	if targetConfig == nil {
		return fmt.Errorf("未找到快捷键 %s 的配置", shortcutKey)
	}

	if targetConfig.TargetDir == "" {
		return fmt.Errorf("快捷键 %s 未配置目标文件夹", shortcutKey)
	}

	// 检查源文件是否存在
	if _, err := os.Stat(filePath); err != nil {
		logger.Error("源文件不存在", zap.String("path", filePath), zap.Error(err))
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	// 解析目标目录
	targetDir := sh.resolveTargetDir(filePath, targetConfig.TargetDir)

	// 创建目标目录
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		logger.Error("创建目标目录失败", zap.String("dir", targetDir), zap.Error(err))
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 构建目标文件路径
	fileName := filepath.Base(filePath)
	targetPath := filepath.Join(targetDir, fileName)

	// 处理同名文件
	targetPath = sh.resolveConflict(targetPath)

	// 移动文件
	if err := file.RenameFile(filePath, targetPath, true, 100); err != nil {
		logger.Error("移动文件失败", zap.String("source", filePath), zap.String("target", targetPath), zap.Error(err))
		return fmt.Errorf("移动文件失败: %w", err)
	}

	// 记录到撤销栈
	sh.pushUndoRecord(MoveRecord{
		SourcePath: filePath,
		TargetPath: targetPath,
		Timestamp:  time.Now(),
	})

	logger.Info("文件已移动",
		zap.String("source", filePath),
		zap.String("target", targetPath),
		zap.String("shortcut", shortcutKey),
		zap.String("label", targetConfig.Label))

	return nil
}

// resolveTargetDir 解析目标目录（支持相对路径）
func (sh *ShortcutHandler) resolveTargetDir(sourcePath, targetDir string) string {
	// 如果是绝对路径，直接返回
	if filepath.IsAbs(targetDir) {
		return targetDir
	}

	// 相对路径：相对于源文件所在目录
	sourceDir := filepath.Dir(sourcePath)
	return filepath.Join(sourceDir, targetDir)
}

// resolveConflict 处理文件名冲突
func (sh *ShortcutHandler) resolveConflict(targetPath string) string {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return targetPath
	}

	dir := filepath.Dir(targetPath)
	ext := filepath.Ext(targetPath)
	name := strings.TrimSuffix(filepath.Base(targetPath), ext)

	for i := 1; i < 1000; i++ {
		newPath := filepath.Join(dir, fmt.Sprintf("%s_%d%s", name, i, ext))
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
	}

	// 如果还是冲突，使用时间戳
	timestamp := time.Now().UnixNano()
	return filepath.Join(dir, fmt.Sprintf("%s_%d%s", name, timestamp, ext))
}

// pushUndoRecord 添加撤销记录
func (sh *ShortcutHandler) pushUndoRecord(record MoveRecord) {
	sh.undoStack = append(sh.undoStack, record)
	if len(sh.undoStack) > sh.maxUndoSize {
		sh.undoStack = sh.undoStack[1:]
	}
}

// Undo 撤销上一次移动操作
func (sh *ShortcutHandler) Undo() error {
	sh.mux.Lock()
	defer sh.mux.Unlock()

	if len(sh.undoStack) == 0 {
		return fmt.Errorf("没有可撤销的操作")
	}

	// 弹出最后一条记录
	record := sh.undoStack[len(sh.undoStack)-1]
	sh.undoStack = sh.undoStack[:len(sh.undoStack)-1]

	// 检查目标文件是否存在
	if _, err := os.Stat(record.TargetPath); err != nil {
		logger.Error("撤销失败：目标文件不存在", zap.String("path", record.TargetPath), zap.Error(err))
		return fmt.Errorf("文件已被删除或移动: %s", record.TargetPath)
	}

	// 确保源目录存在
	sourceDir := filepath.Dir(record.SourcePath)
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		logger.Error("创建源目录失败", zap.String("dir", sourceDir), zap.Error(err))
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 移回原位置
	if err := file.RenameFile(record.TargetPath, record.SourcePath, true, 100); err != nil {
		logger.Error("撤销移动失败", zap.String("source", record.TargetPath), zap.String("target", record.SourcePath), zap.Error(err))
		return fmt.Errorf("撤销失败: %w", err)
	}

	logger.Info("已撤销移动操作",
		zap.String("from", record.TargetPath),
		zap.String("to", record.SourcePath))

	return nil
}

// GetUndoCount 获取可撤销操作数量
func (sh *ShortcutHandler) GetUndoCount() int {
	return len(sh.undoStack)
}

// ClearUndoStack 清空撤销栈
func (sh *ShortcutHandler) ClearUndoStack() {
	sh.mux.Lock()
	defer sh.mux.Unlock()
	sh.undoStack = make([]MoveRecord, 0)
	logger.Info("撤销栈已清空")
}

// GetLastMoveRecord 获取最后一次移动记录（用于前端显示）
func (sh *ShortcutHandler) GetLastMoveRecord() *MoveRecord {
	if len(sh.undoStack) == 0 {
		return nil
	}
	record := sh.undoStack[len(sh.undoStack)-1]
	return &record
}
