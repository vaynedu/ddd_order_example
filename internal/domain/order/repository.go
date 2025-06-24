package order

import "context"

// OrderRepository 订单仓储接口
type OrderRepository interface {
    Save(ctx context.Context, order *Order) error
    FindByID(ctx context.Context, id string) (*Order, error)
}  