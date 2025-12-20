package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithOrderly(t *testing.T) {
	root := "xxx"
	err := WithOrderly(root, 3)
	assert.Nil(t, err)
}

// TestHasHiddenFiles 检查是否有 . 起头的文件（除了系统隐藏文件）
func TestHasHiddenFiles(t *testing.T) {
	rootDir := "xxx"

	err := filepath.WalkDir(rootDir, func(fullPath string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("遍历错误 %s", err)
			return nil
		}

		if strings.HasPrefix(d.Name(), ".") || filterFiles[d.Name()] {
			if filterFiles[d.Name()] {
				return nil
			}
			fmt.Printf("发现 %s\n", fullPath)
			return nil
		}
		return nil
	})

	// 处理整体遍历的致命错误
	if err != nil {
		fmt.Printf("递归遍历失败: %v", err)
		os.Exit(1)
	}

	fmt.Printf("目录 %s 遍历完成", rootDir)
}
