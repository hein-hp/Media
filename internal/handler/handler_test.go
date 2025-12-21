package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/corona10/goimagehash"
	"github.com/stretchr/testify/assert"
)

func TestImage(t *testing.T) {
	dir := "xxx"

	// 初始化映射
	fileMap := make(map[string]*os.File)
	hashMap := make(map[string]*goimagehash.ImageHash)
	var imgPaths []string

	// 遍历目录，收集所有.jpg文件
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk dir error: %w", err)
		}
		if d.IsDir() {
			return nil
		}
		// 只处理.jpg后缀
		if ext := filepath.Ext(path); ext != ".jpg" {
			return nil
		}
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open file %s error: %w", path, err)
		}
		fileMap[path] = f
		imgPaths = append(imgPaths, path)
		return nil
	})
	assert.Nil(t, err, "遍历目录收集文件失败")

	// 为每个文件生成图片哈希，并关闭文件
	defer func() {
		for _, f := range fileMap {
			_ = f.Close()
		}
	}()

	hashResults := calcAverageHash(imgPaths, fileMap)
	for _, hashResult := range hashResults {
		if hashResult.hash == nil {
			continue
		}
		hashMap[hashResult.path] = hashResult.hash
	}

	// 两两计算汉明距离
	imgCount := len(hashMap)
	t.Logf("共找到 %d 张图片，开始计算两两汉明距离\n", imgCount)

	for i := 0; i < imgCount; i++ {
		path1 := imgPaths[i]
		hash1 := hashMap[path1]
		for j := i + 1; j < imgCount; j++ {
			path2 := imgPaths[j]
			hash2 := hashMap[path2]

			// 计算汉明距离
			distance, err := hash1.Distance(hash2)
			if err != nil {
				t.Errorf("计算 %s 和 %s 汉明距离失败: %v", path1, path2, err)
				continue
			}

			// 生成相似度说明
			if distance == 0 {
				// 同时输出到测试日志
				t.Logf("【图片1】%s\n【图片2】%s\n【汉明距离】%d\n------------------------\n",
					path1, path2, distance)
			}
		}
	}
}
