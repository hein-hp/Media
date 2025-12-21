package handler

import "context"

// Handler 处理器
type Handler interface {
	SetContext(ctx context.Context)

	GetSelectedDir() string
	SetSelectedDir(dir string)
}
