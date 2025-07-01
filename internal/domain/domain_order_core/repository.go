package domain_order_core

import "context"

// OrderRepository 订单仓储接口
type OrderRepository interface {
	Save(ctx context.Context, order *OrderDO) error
	FindByID(ctx context.Context, id string) (*OrderDO, error)
}
