package payment

import (
	"context"



	"github.com/vaynedu/ddd_order_example/internal/domain/domain_payment_core"
)

// 支付代理接口（与外部支付系统通信）
type PaymentProxy interface {
	CreatePayment(ctx context.Context, orderID string, amount int64) (string, error)
	QueryPaymentStatus(ctx context.Context, paymentID string) (*domain_payment_core.PaymentStatus, error)
	// QueryPayment(ctx context.Context, paymentID string) (*domain_payment_core.PaymentDO, error)
}
