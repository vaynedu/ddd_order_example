package domain_payment_core

import "errors"

// 支付领域错误定义
var (
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrPaymentAlreadyExists  = errors.New("payment already exists")
	ErrInvalidPaymentStatus  = errors.New("invalid payment status")
	ErrPaymentAmountMismatch = errors.New("payment amount mismatch")
	ErrPaymentPaid           = errors.New("订单已支付，无需重复操作")
)
