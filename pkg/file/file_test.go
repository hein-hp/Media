package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountFiles(t *testing.T) {
	files, err := CountFiles("/Users/xxx/xxx")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 20, files)
}

func TestGetFileMeta(t *testing.T) {
	meta, err := GetFileMeta("xxx")
	assert.Nil(t, err)
	fmt.Println(meta)
}

func TestRemoveFile(t *testing.T) {

	ext := "xmp"
	dir := "xxx"
	recursive := true

	// 标准化扩展名
	targetExt := strings.ToLower(ext)
	if !strings.HasPrefix(targetExt, ".") {
		targetExt = "." + targetExt
	}

	// 校验目标文件夹是否存在
	dirInfo, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("文件夹不存在：%s", dir)
			return
		}
		t.Errorf("获取文件夹信息失败：%s", err)
		return
	}
	if !dirInfo.IsDir() {
		t.Errorf("指定路径不是文件夹：%s", dir)
		return
	}

	// 遍历目录并删除对应文件
	var deleteCount int
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		// 处理遍历过程中的错误（如权限不足无法访问子目录）
		if err != nil {
			return fmt.Errorf("遍历路径 %s 失败：%w", path, err)
		}

		// 跳过目录（若不递归，直接返回跳过子目录）
		if d.IsDir() {
			if !recursive && path != dir {
				return filepath.SkipDir // 不递归时，跳过所有子目录
			}
			return nil
		}

		// 匹配文件扩展名
		fileExt := strings.ToLower(filepath.Ext(path))
		if fileExt != targetExt {
			return nil
		}

		// 删除文件
		if err := os.Remove(path); err != nil {
			fmt.Printf("删除文件 %s 失败：%s", path, err)
			return nil
		}
		deleteCount++
		fmt.Printf("已删除文件：%s\n", path)
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("操作完成！共删除 %d 个 %s 类型文件\n", deleteCount, targetExt)
}

func TestRenameExt(t *testing.T) {
	renameMapping := map[string]string{
		"JPG":  "jpg",
		"HEIC": "jpg",
		"MOV":  "mp4",
		"PNG":  "png",
	}

	dir := "xxx"
	recursive := true
	skipHiddenFiles := false // 是否跳过隐藏文件

	dirInfo, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("文件夹不存在：%s", dir)
		}
		t.Fatalf("获取文件夹信息失败：%v", err)
	}
	if !dirInfo.IsDir() {
		t.Fatalf("指定路径不是文件夹：%s", dir)
	}

	// 统计信息
	var (
		successCount int // 重命名成功数量
		failCount    int // 重命名失败数量
		skipCount    int // 跳过的文件数量
	)

	// 遍历目录并处理文件重命名
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		// 处理遍历错误
		if err != nil {
			t.Errorf("遍历路径 %s 失败：%v", path, err)
			return nil
		}

		// 跳过目录
		if d.IsDir() {
			if !recursive && path != dir {
				skipCount++
				return filepath.SkipDir
			}
			return nil
		}

		// 跳过隐藏文件
		if skipHiddenFiles && strings.HasPrefix(d.Name(), ".") {
			skipCount++
			return nil
		}

		// 提取文件扩展名（带点，如 .JPG）
		rawExt := filepath.Ext(path)
		if rawExt == "" {
			skipCount++
			return nil // 无扩展名文件，跳过
		}

		// 标准化扩展名：去掉点 + 转为大写，实现大小写不敏感匹配
		srcExt := strings.TrimPrefix(strings.ToUpper(rawExt), ".")
		targetExt, needRename := renameMapping[srcExt]
		if !needRename {
			skipCount++
			return nil // 无需重命名的扩展名，跳过
		}

		// 获取文件目录
		fileDir := filepath.Dir(path)
		// 获取文件名
		fileNameWithoutExt := strings.TrimSuffix(filepath.Base(path), rawExt)
		// 拼接新文件名
		newFileName := fmt.Sprintf("%s.%s", fileNameWithoutExt, targetExt)
		// 拼接完整目标路径
		targetPath := filepath.Join(fileDir, newFileName)
		// 标准化路径
		targetPath = filepath.Clean(targetPath)

		// 处理目标文件已存在的冲突
		if _, err := os.Stat(targetPath); err == nil {
			t.Errorf("目标文件已存在，跳过重命名：%s -> %s", path, targetPath)
			failCount++
			return nil
		} else if !os.IsNotExist(err) {
			t.Errorf("检查目标文件状态失败：%s，错误：%v", targetPath, err)
			failCount++
			return nil
		}

		// 执行重命名
		if err := os.Rename(path, targetPath); err != nil {
			t.Errorf("重命名失败：%s -> %s，错误：%v", path, targetPath, err)
			failCount++
			return nil
		}

		successCount++
		t.Logf("重命名成功：%s -> %s", path, targetPath)
		return nil
	})

	// 处理遍历整体错误
	if err != nil {
		t.Fatalf("目录遍历异常终止：%v", err)
	}

	// 输出统计结果
	t.Logf("==================== 操作汇总 ====================")
	t.Logf("目标文件夹：%s", dir)
	t.Logf("重命名成功：%d 个", successCount)
	t.Logf("重命名失败：%d 个", failCount)
	t.Logf("跳过文件：%d 个", skipCount)
	t.Logf("==================================================")
}
