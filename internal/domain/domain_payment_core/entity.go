package domain_payment_core

import (
	"context"
	"time"
)

// PaymentDO 支付领域对象
type PaymentDO struct {
	ID                  string
	OrderID             string
	Amount              int64
	Currency            string
	Channel             int
	Status              PaymentStatus
	TransactionID       string
	RefundTransactionID string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	CompletedAt         *time.Time
}

// 支付状态
type PaymentStatus int

const (
	PaymentStatusCreated         PaymentStatus = iota // 已创建
	PaymentStatusPaid                                 // 已支付
	PaymentStatusRefunded                             // 已退款
	PaymentStatusFailed                               // 支付失败
	PaymentStatusExpired                              // 已过期
	PaymentStatusCanceled                             // 已取消
	PaymentStatusRefunding                            // 退款中
	PaymentStatusRefundFailed                         // 退款失败
	PaymentStatusRefundedSuccess                      // 退款成功
	PaymentStatusCompleted                            // 已完成
	PaymentStatusClosed                               // 已关闭
	PaymentStatusPending                              // 待支付
)

// 支付仓储接口
type Repository interface {
	Save(ctx context.Context, payment *PaymentDO) error
	FindByID(ctx context.Context, id string) (*PaymentDO, error)
	FindByOrderID(ctx context.Context, orderID string) (*PaymentDO, error)
}
