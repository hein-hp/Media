package handler

import (
	"context"
	"fmt"
	"image/jpeg"
	"media-app/pkg/file"
	"media-app/pkg/logger"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/corona10/goimagehash"
	"go.uber.org/zap"
)

// SimilarHandler handles similarity analysis
type SimilarHandler struct {
	dir string
	mux sync.Mutex
	ctx context.Context
}

// HashResult 哈希结果
type HashResult struct {
	path string
	hash *goimagehash.ImageHash
}

// SimilarityResult 相似度分析结果
type SimilarityResult struct {
	Images []file.Meta
}

// NewSimilarHandler creates a new SimilarHandler instance
func NewSimilarHandler() *SimilarHandler {
	return &SimilarHandler{}
}

// SetContext sets the wails context
func (sh *SimilarHandler) SetContext(ctx context.Context) {
	sh.ctx = ctx
}

// GetSelectedDir returns the selected directory
func (sh *SimilarHandler) GetSelectedDir() string {
	return sh.dir
}

// SetSelectedDir sets the selected directory
func (sh *SimilarHandler) SetSelectedDir(dir string) {
	sh.mux.Lock()
	defer sh.mux.Unlock()
	sh.dir = dir
	logger.Info("目录已选择", zap.String("dir", dir))
}

// CalcSimilarity 计算相似度
func (sh *SimilarHandler) CalcSimilarity() []SimilarityResult {
	dir := sh.GetSelectedDir()
	if dir == "" {
		logger.Warn("未选择文件夹")
		return nil
	}

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
	if err != nil {
		logger.Error("遍历目录收集文件失败", zap.String("dir", dir), zap.Error(err))
	}

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

	// 过滤出有效的图片路径（有哈希值的）
	var validPaths []string
	for _, path := range imgPaths {
		if hashMap[path] != nil {
			validPaths = append(validPaths, path)
		}
	}

	imgCount := len(validPaths)
	logger.Infof("找到 %d 张有效图片，开始计算两两汉明距离", imgCount)

	if imgCount == 0 {
		return nil
	}

	// 初始化并查集
	uf := newUnionFind(validPaths)

	// 两两计算汉明距离，distance为0时合并到同一组
	for i := 0; i < imgCount; i++ {
		path1 := validPaths[i]
		hash1 := hashMap[path1]
		for j := i + 1; j < imgCount; j++ {
			path2 := validPaths[j]
			hash2 := hashMap[path2]

			// 计算汉明距离
			distance, err := hash1.Distance(hash2)
			if err != nil {
				logger.Errorf("计算 %s 和 %s 汉明距离失败: %v", path1, path2, err)
				continue
			}

			// 距离为0时，合并到同一组
			if distance == 0 {
				logger.Info("发现相同图片",
					zap.String("图片1", path1), zap.String("图片2", path2), zap.Int("汉明距离", distance))
				uf.union(path1, path2)
			}
		}
	}

	// 收集分组结果
	groups := uf.getGroups()

	// 构建 SimilarityResult 切片（只包含有多张图片的组）
	var results []SimilarityResult
	for _, paths := range groups {
		if len(paths) < 2 {
			// 跳过只有一张图片的组
			continue
		}

		var metas []file.Meta
		for _, path := range paths {
			meta, err := file.GetFileMeta(path)
			if err != nil {
				logger.Error("获取文件元数据失败", zap.String("path", path), zap.Error(err))
				continue
			}
			metas = append(metas, *meta)
		}

		if len(metas) >= 2 {
			results = append(results, SimilarityResult{Images: metas})
			logger.Infof("找到相似图片组，共 %d 张图片", len(metas))
		}
	}

	logger.Infof("相似度分析完成，共找到 %d 组相似图片", len(results))
	return results
}

// unionFind 并查集数据结构
type unionFind struct {
	parent map[string]string
	rank   map[string]int
}

// newUnionFind 创建并查集
func newUnionFind(paths []string) *unionFind {
	uf := &unionFind{
		parent: make(map[string]string),
		rank:   make(map[string]int),
	}
	for _, path := range paths {
		uf.parent[path] = path
		uf.rank[path] = 0
	}
	return uf
}

// find 查找根节点（带路径压缩）
func (uf *unionFind) find(x string) string {
	if uf.parent[x] != x {
		uf.parent[x] = uf.find(uf.parent[x])
	}
	return uf.parent[x]
}

// union 合并两个集合（按秩合并）
func (uf *unionFind) union(x, y string) {
	rootX := uf.find(x)
	rootY := uf.find(y)

	if rootX == rootY {
		return
	}

	// 按秩合并
	if uf.rank[rootX] < uf.rank[rootY] {
		uf.parent[rootX] = rootY
	} else if uf.rank[rootX] > uf.rank[rootY] {
		uf.parent[rootY] = rootX
	} else {
		uf.parent[rootY] = rootX
		uf.rank[rootX]++
	}
}

// getGroups 获取所有分组
func (uf *unionFind) getGroups() map[string][]string {
	groups := make(map[string][]string)
	for path := range uf.parent {
		root := uf.find(path)
		groups[root] = append(groups[root], path)
	}
	return groups
}

func calcAverageHash(imgPath []string, fileMap map[string]*os.File) []HashResult {
	groupSize := 200
	groups := splitSlice(imgPath, groupSize)
	logger.Infof("切片拆分为 %d 组，每组最多 %d 个元素", len(groups), groupSize)

	resChan := make(chan HashResult, len(groups))

	var wg sync.WaitGroup
	wg.Add(len(groups))
	for _, group := range groups {
		go func(group []string) {
			defer wg.Done()

			for _, path := range group {
				f := fileMap[path]
				if f == nil {
					resChan <- HashResult{path, nil}
					continue
				}
				img, err := jpeg.Decode(f)
				if err != nil {
					logger.Errorf("解码图片 %s 失败: %v", path, err)
					resChan <- HashResult{path, nil}
					continue
				}
				hash, err := goimagehash.AverageHash(img)
				if err != nil {
					logger.Errorf("计算哈希 %s 失败: %v", path, err)
					resChan <- HashResult{path, nil}
					continue
				}
				resChan <- HashResult{path, hash}
			}
		}(group)
	}

	go func() {
		wg.Wait()
		close(resChan)
	}()

	var results []HashResult
	for res := range resChan {
		results = append(results, res)
	}
	return results
}

func splitSlice(paths []string, groupSize int) [][]string {
	if groupSize <= 0 {
		groupSize = 200
	}
	var groups [][]string
	length := len(paths)

	for start := 0; start < length; start += groupSize {
		end := start + groupSize
		if end > length {
			end = length
		}
		groups = append(groups, paths[start:end])
	}
	return groups
}
