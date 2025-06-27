package payment

import (
	"context"
	"strconv"

	"github.com/smartwalle/alipay/v3"
)

// 支付宝适配器实现
type AlipayAdapter struct {
	client *alipay.Client
}

// NewAlipayAdapter 创建支付宝适配器实例
func NewAlipayAdapter(appID, privateKey string, debugMode bool) (*AlipayAdapter, error) {
	client, err := alipay.New(appID, privateKey, debugMode)
	if err != nil {
		return nil, err
	}

	return &AlipayAdapter{client: client}, nil
}

// CreatePayment 发起支付宝支付
func (a *AlipayAdapter) CreatePayment(ctx context.Context, orderID string, amount int64) (string, error) {
	// 构建支付请求参数
	payReq := alipay.TradePagePay{}
	payReq.NotifyURL = "http://notify.example.com" // 支付通知地址
	payReq.ReturnURL = "http://return.example.com" // 支付成功返回地址
	payReq.Subject = "订单支付"
	payReq.OutTradeNo = orderID
	payReq.TotalAmount = strconv.FormatInt(amount, 10)
	payReq.ProductCode = "FAST_INSTANT_TRADE_PAY"

	// 发起支付请求
	payResp, err := a.client.TradePagePay(payReq)
	if err != nil {
		return "", err
	}

	// fixme 这里先写mock数据
	return payResp.Fragment, nil
}
