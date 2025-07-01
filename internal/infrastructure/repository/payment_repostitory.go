package repository

import (
	"context"

	"github.com/vaynedu/ddd_order_example/internal/domain/domain_payment_core"
	"gorm.io/gorm"
)

// PaymentRepositoryMySQL MySQL实现的支付仓储
type PaymentRepositoryMySQL struct {
	db *gorm.DB
}

// NewPaymentRepository 创建订单仓储实例
func NewPaymentRepository(db *gorm.DB) domain_payment_core.Repository {
	return &PaymentRepositoryMySQL{db: db}
}

// Save 保存支付记录
func (r *PaymentRepositoryMySQL) Save(ctx context.Context, payment *domain_payment_core.PaymentDO) error {
	return r.db.WithContext(ctx).Table("t_payment").Save(payment).Error
}

// FindByID 根据ID查询支付记录
func (r *PaymentRepositoryMySQL) FindByID(ctx context.Context, id string) (*domain_payment_core.PaymentDO, error) {
	var payment domain_payment_core.PaymentDO
	err := r.db.WithContext(ctx).Table("t_payment").Where("id = ?", id).First(&payment).Error
	return &payment, err
}

// FindByOrderID 根据订单ID查询支付记录
func (r *PaymentRepositoryMySQL) FindByOrderID(ctx context.Context, orderID string) (*domain_payment_core.PaymentDO, error) {
	var payment domain_payment_core.PaymentDO
	err := r.db.WithContext(ctx).Table("t_payment").Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}
