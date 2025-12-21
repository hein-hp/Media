package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

func TestMoveImageFile(t *testing.T) {
	srcDir := "xxx"
	destDir := "xxx"

	// 1. 源目录校验
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("源目录不存在: %s", srcDir)
			return
		}
		t.Errorf("获取源目录信息失败: %v", err)
		return
	}
	if !srcInfo.IsDir() {
		t.Errorf("srcDir 不是一个有效的目录: %s", srcDir)
		return
	}

	// 2. 创建目标目录
	if err := os.MkdirAll(destDir, 0755); err != nil {
		t.Errorf("创建目标目录失败: %v", err)
		return
	}

	// 3. 图片扩展名定义
	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".webp": true,
	}

	// 4. 多协程相关初始化
	const maxConcurrency = 10 // 最大并发数（可根据机器性能调整，机械硬盘建议5-10，固态硬盘可20-50）
	// 协程池：带缓冲通道，控制并发数
	workerChan := make(chan struct{}, maxConcurrency)
	// 协程同步等待组
	var wg sync.WaitGroup
	// 互斥锁：保护共享计数器cnt
	var mutex sync.Mutex
	// 错误通道：收集所有协程的错误信息
	errChan := make(chan error, 100) // 缓冲大小可根据文件数量调整
	cnt := 0                         // 成功移动的图片数量

	// 5. 递归遍历源目录（不变）
	traverseErr := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("遍历文件失败 %s: %w", path, err)
		}
		// 跳过目录
		if d.IsDir() {
			return nil
		}
		// 筛选图片文件
		fileExt := strings.ToLower(filepath.Ext(path))
		if !imageExts[fileExt] {
			return nil
		}

		// 6. 多协程执行移动逻辑
		workerChan <- struct{}{} // 占用协程池资源，达到最大并发则阻塞
		wg.Add(1)                // 等待组计数+1

		go func(srcPath string) {
			// 协程退出时释放资源
			defer func() {
				<-workerChan // 释放协程池资源
				wg.Done()    // 等待组计数-1
			}()

			// 构造目标路径
			fileName := filepath.Base(srcPath)
			destPath := filepath.Join(destDir, fileName)

			// 执行移动
			if err := RenameFile(srcPath, destPath, true, 1000); err != nil {
				errChan <- fmt.Errorf("移动文件 %s 失败: %w", srcPath, err)
				return
			}

			// 安全修改共享计数器
			mutex.Lock()
			cnt++
			mutex.Unlock()

			fmt.Printf("成功移动图片: %s -> %s\n", srcPath, destPath)
		}(path) // 传入当前文件路径，避免闭包延迟绑定问题

		return nil
	})

	// 7. 等待所有移动协程执行完成
	go func() {
		wg.Wait()
		close(errChan)    // 所有协程执行完后关闭错误通道
		close(workerChan) // 关闭协程池通道
	}()

	// 8. 收集并处理所有错误
	var allErrors []error
	for err := range errChan {
		allErrors = append(allErrors, err)
	}

	// 9. 处理遍历过程中的错误
	if traverseErr != nil {
		t.Errorf("批量移动图片失败: %v", traverseErr)
	}

	// 10. 输出错误信息（若有）
	if len(allErrors) > 0 {
		t.Errorf("共发生 %d 个移动错误：", len(allErrors))
		for _, err := range allErrors {
			t.Error(err)
		}
	}

	// 11. 输出最终统计结果
	fmt.Printf("共移动 %d 张图片\n", cnt)
}
func TestRemoveHiddenFile(t *testing.T) {
	dir := "/Users/hejin/.mounty/难以描述/个人珍藏/temp/新建文件夹 (2)"

	// 读取当前目录下的所有目录条目
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("读取当前目录失败：%v\n", err)
		return
	}

	// 遍历所有目录条目，筛选并删除符合条件的文件
	for _, entry := range entries {
		filename := entry.Name()

		// 排除特殊目录 "." 和 ".."
		if filename == "." || filename == ".." {
			continue
		}

		// 筛选：文件名以 "." 开头
		if len(filename) == 0 || filename[0] != '.' {
			continue
		}

		// 筛选：普通文件（对应 -type f），排除目录、链接等其他类型
		fileInfo, err := entry.Info()
		if err != nil {
			fmt.Printf("获取文件信息失败（%s）：%v\n", filename, err)
			continue
		}
		if !fileInfo.Mode().IsRegular() { // 判断是否为普通文件
			continue
		}

		// 拼接完整文件路径
		filePath := filepath.Join(dir, filename)

		// 删除文件（对应 -delete）
		err = os.Remove(filePath)
		if err != nil {
			fmt.Printf("删除文件失败（%s）：%v\n", filePath, err)
		} else {
			fmt.Printf("成功删除文件：%s\n", filePath)
		}
	}

	fmt.Println("所有符合条件的文件处理完毕")
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
		"MP4":  "mp4",
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
		targetExt := renameMapping[srcExt]
		if targetExt == "" {
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

		if targetPath == path {
			skipCount++
			return nil
		}

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
