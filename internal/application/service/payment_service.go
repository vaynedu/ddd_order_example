package service

import (
	"context"

	"github.com/vaynedu/ddd_order_example/internal/domain/domain_payment_core"
	"github.com/vaynedu/ddd_order_example/internal/infrastructure/payment"
)

type PaymentService struct {
	domainService *domain_payment_core.PaymentDomainService
	paymentProxy  payment.PaymentProxy
}

func NewPaymentService(domainService *domain_payment_core.PaymentDomainService, paymentProxy payment.PaymentProxy) *PaymentService {
	return &PaymentService{
		domainService: domainService,
		paymentProxy:  paymentProxy,
	}
}

// 创建支付请求
func (s *PaymentService) CreatePayment(ctx context.Context, orderID string, amount int64, currency string, channel int) (string, error) {
	// 1. 创建支付记录
	paymentDO, err := s.domainService.CreatePayment(ctx, orderID, amount, currency, channel)
	if err != nil {
		return "", err
	}

	// 2. 调用外部支付系统
	transactionID, err := s.paymentProxy.CreatePayment(ctx, orderID, amount)
	if err != nil {
		// 支付失败，更新支付状态
		_ = s.domainService.ProcessPaymentResult(ctx, paymentDO.ID, "", false)
		return "", err
	}

	// 3. 更新支付记录
	return paymentDO.ID, s.domainService.ProcessPaymentResult(ctx, paymentDO.ID, transactionID, true)
}
