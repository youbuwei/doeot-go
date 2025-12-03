package domain

import "context"

// Order 是 order 模块的领域模型示例，你可以按需扩展字段。
type Order struct {
	ID   int64
	Name string
}

// NotFoundError 用于标识未找到该资源。
type NotFoundError struct {
	msg string
}

func (e *NotFoundError) Error() string { return e.msg }

// ErrOrderNotFound 在仓储查不到数据时返回。
var ErrOrderNotFound = &NotFoundError{msg: "order not found"}

// Repo 抽象了针对 Order 的持久化操作。
type Repo interface {
	FindByID(ctx context.Context, id int64) (*Order, error)
	Create(ctx context.Context, m *Order) (*Order, error)
	List(ctx context.Context) ([]*Order, error)
}
