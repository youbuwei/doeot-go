package domain

import "context"

// Pay 是 pay 模块的领域模型示例，你可以按需扩展字段。
type Pay struct {
	ID   int64
	Name string
}

// NotFoundError 用于标识未找到该资源。
type NotFoundError struct {
	msg string
}

func (e *NotFoundError) Error() string { return e.msg }

// ErrPayNotFound 在仓储查不到数据时返回。
var ErrPayNotFound = &NotFoundError{msg: "pay not found"}

// Repo 抽象了针对 Pay 的持久化操作。
type Repo interface {
	FindByID(ctx context.Context, id int64) (*Pay, error)
	Create(ctx context.Context, m *Pay) (*Pay, error)
	List(ctx context.Context) ([]*Pay, error)
}
