package domain_payment_core

import (
	"context"
	"time"
)

type PaymentDomainService struct {
	repo Repository
}

// CreatePayment 创建支付
func (s *PaymentDomainService) CreatePayment(ctx context.Context, orderID string, amount int64, currency string, channel int) (*PaymentDO, error) {
	paymentDO := &PaymentDO{
		OrderID:   orderID,
		Amount:    amount,
		Currency:  currency,
		Channel:   channel,
		Status:    PaymentStatusCreated,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return paymentDO, s.repo.Save(ctx, paymentDO)
}

// 处理支付结果
func (s *PaymentDomainService) ProcessPaymentResult(ctx context.Context, paymentID, transactionID string, success bool) error {
	payment, err := s.repo.FindByID(ctx, paymentID)
	if err != nil {
		return err
	}

	if success {
		now := time.Now()
		payment.Status = PaymentStatusCompleted
		payment.TransactionID = transactionID
		payment.CompletedAt = &now
	} else {
		payment.Status = PaymentStatusFailed
	}

	return s.repo.Save(ctx, payment)
}
