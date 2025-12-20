package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

// renameItem 存储文件重命名的原路径、临时路径、最终路径
type renameItem struct {
	oldPath   string // 原始文件路径
	tempPath  string // 临时文件路径
	finalPath string // 最终文件路径
}

// WithOrderly 对目录下文件按更新时间逆序进行有序重命名
// dir: 目标目录
// length: 序号位数（如4位生成 0001、0002）
func WithOrderly(dir string, length int) error {
	// 前置参数校验
	if dir == "" {
		return fmt.Errorf("目标目录不能为空")
	}
	if length <= 0 {
		return fmt.Errorf("序号位数必须大于0，当前为%d", length)
	}

	// 获取目录下有效文件元数据
	metaList, err := GetFileMetas(dir)
	if err != nil {
		return fmt.Errorf("获取文件元数据失败：%w", err)
	}
	fileCount := len(metaList)
	if fileCount == 0 {
		return fmt.Errorf("无需排序，文件数量为 0, path: %s ", dir)
	}

	// 校验序号位数是否足够容纳文件数量
	maxSeq := fileCount // 最大序号（更新时间最新的文件对应此序号）
	digitCount := len(strconv.Itoa(maxSeq))
	if digitCount > length {
		return fmt.Errorf("序号位数不足：文件数量为%d（需%d位序号），当前指定位数为%d",
			fileCount, digitCount, length)
	}

	// 按文件更新时间逆序排序
	sort.Slice(metaList, func(i, j int) bool {
		return metaList[i].ModTime.Before(metaList[j].ModTime)
	})

	// 生成重命名映射
	var renameItems []renameItem
	for idx, meta := range metaList {
		seqStr := fmt.Sprintf("%0*d", length, idx+1)
		finalPath := filepath.Join(dir, seqStr+meta.Ext)
		tempPath := meta.FullPath + ".tmp"

		renameItems = append(renameItems, renameItem{
			oldPath:   meta.FullPath,
			tempPath:  tempPath,
			finalPath: finalPath,
		})
	}

	// 将所有文件重命名为临时文件（避免直接重命名导致的文件名冲突）
	for _, item := range renameItems {
		if err := os.Rename(item.oldPath, item.tempPath); err != nil {
			return fmt.Errorf("临时重命名失败 %s -> %s：%w", item.oldPath, item.tempPath, err)
		}
	}

	// 将临时文件重命名为最终有序文件
	for _, item := range renameItems {
		if err := os.Rename(item.tempPath, item.finalPath); err != nil {
			return fmt.Errorf("最终重命名失败 %s -> %s：%w", item.tempPath, item.finalPath, err)
		}
	}

	return nil
}
