package domain_payment_core

import "context"

// 支付仓储接口
type Repository interface {
	Save(ctx context.Context, payment *PaymentDO) error
	FindByID(ctx context.Context, id string) (*PaymentDO, error)
	FindByOrderID(ctx context.Context, orderID string) (*PaymentDO, error)
}