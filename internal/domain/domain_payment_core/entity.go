package domain_payment_core

import (
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
	PaymentStatusUnknown         PaymentStatus = iota // 未知
	PaymentStatusCreated                              // 已创建
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

func GetPaymentStatusDetail(status PaymentStatus) string {
	switch status {
	case PaymentStatusCreated:
		return "已创建"
	case PaymentStatusPaid:
		return "已支付"
	case PaymentStatusRefunded:
		return "已退款"
	case PaymentStatusFailed:
		return "支付失败"
	case PaymentStatusExpired:
		return "已过期"
	case PaymentStatusCanceled:
		return "已取消"
	case PaymentStatusRefunding:
		return "退款中"
	case PaymentStatusRefundFailed:
		return "退款失败"
	case PaymentStatusRefundedSuccess:
		return "退款成功"
	case PaymentStatusCompleted:
		return "已完成"
	case PaymentStatusClosed:
		return "已关闭"
	case PaymentStatusPending:
		return "待支付"
	default:
		return "未知"
	}
}

// 支付渠道枚举
type PaymentChannel int

const (
	PaymentChannelAlipay   PaymentChannel = iota // 支付宝
	PaymentChannelWechat                         // 微信
	PaymentChannelUnionPay                       // 银联
	PaymentChannelApplePay                       // ApplePay
	PaymentChannelJDPay                          // 京东支付
)

func GetPaymentChannelDetail(channel PaymentChannel) string {
	switch channel {
	case PaymentChannelAlipay:
		return "支付宝"
	case PaymentChannelWechat:
		return "微信"
	case PaymentChannelUnionPay:
		return "银联"
	case PaymentChannelApplePay:
		return "ApplePay"
	case PaymentChannelJDPay:
		return "京东支付"
	default:
		return "未知"
	}
}