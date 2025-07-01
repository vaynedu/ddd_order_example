package payment

import (
	"context"
	"fmt"

	"github.com/vaynedu/ddd_order_example/internal/domain/domain_payment_core"
)

// MockPaymentProxy 支付代理的Mock实现
type MockPaymentProxy struct {
	// 控制返回结果的标志
	ReturnSuccess bool
	// 自定义返回的交易ID
	TransactionID string
	// 自定义错误
	CustomError error
}

// NewMockPaymentProxy 创建Mock支付代理,具体实现
func NewMockPaymentProxy() PaymentProxy {
	return &MockPaymentProxy{
		ReturnSuccess: true,
		TransactionID: "mock_transaction_123",
	}
}

// CreatePayment 模拟创建支付请求
func (m *MockPaymentProxy) CreatePayment(ctx context.Context, orderID string, amount int64) (string, error) {
	if m.CustomError != nil {
		return "", m.CustomError
	}

	if !m.ReturnSuccess {
		return "", fmt.Errorf("payment failed: mock error")
	}

	return m.TransactionID, nil
}

// WithReturnSuccess 设置是否返回成功
func (m *MockPaymentProxy) WithReturnSuccess(success bool) *MockPaymentProxy {
	m.ReturnSuccess = success
	return m
}

// WithTransactionID 设置返回的交易ID
func (m *MockPaymentProxy) WithTransactionID(txID string) *MockPaymentProxy {
	m.TransactionID = txID
	return m
}

// WithError 设置自定义错误
func (m *MockPaymentProxy) WithError(err error) *MockPaymentProxy {
	m.CustomError = err
	return m
}

func (m *MockPaymentProxy) WithCustomError(err error) *MockPaymentProxy {
	m.CustomError = err
	return m
}

func (m *MockPaymentProxy) QueryPaymentStatus(ctx context.Context, paymentID string) (domain_payment_core.PaymentStatus, error) {
	if m.CustomError != nil {
		return domain_payment_core.PaymentStatusCanceled, m.CustomError
	}

	return domain_payment_core.PaymentStatusCompleted, nil
}
