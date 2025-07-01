package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vaynedu/ddd_order_example/internal/domain/domain_order_core"
	"github.com/vaynedu/ddd_order_example/internal/domain/domain_payment_core"
	"github.com/vaynedu/ddd_order_example/internal/domain/domain_product_core"
	"gorm.io/gorm"
)

type OrderService struct {
	productService     domain_product_core.ProductService   // 依赖商品领域接口
	orderDomainService domain_order_core.OrderDomainService // 依赖订单领域服务
	paymentService     *PaymentService                      // 注入依赖支付服务
}

func NewOrderService(orderDomainService domain_order_core.OrderDomainService, paymentService *PaymentService, productService domain_product_core.ProductService) *OrderService {
	return &OrderService{
		orderDomainService: orderDomainService,
		paymentService:     paymentService,
		productService:     productService,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, customerID string, items []*domain_order_core.OrderItemDO) (string, error) {
	// todo 考虑分布式锁、业务幂等， 防止重复创建单子
	// todo 假装这里只有一个商品
	// todo 这块逻辑要重新封装处理
	// 验证商品状态, 并获取商品信息
	req := &domain_product_core.ValidateProductRequest{
		ProductID: items[0].ProductID,
		Name:      "",
		Price:     items[0].UnitPrice,
		Quantity:  items[0].Quantity,
	}
	resp, err := s.productService.ValidateProduct(ctx, req)
	if err != nil {
		return "", err
	}
	if resp.IsValid == false {
		return "", errors.New(resp.Messages)
	}
	if resp.Product.Status != domain_product_core.StatusValid {
		return "", errors.New("product is not available")
	}

	p := resp.Product

	// 创建订单
	newOrder := &domain_order_core.OrderDO{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		Status:     domain_order_core.OrderStatusCreated,
		Items: []domain_order_core.OrderItemDO{
			{
				ProductID: p.ID,
				Quantity:  items[0].Quantity,
				UnitPrice: items[0].UnitPrice,
				Subtotal:  items[0].Subtotal,
			},
		},
		TotalAmount: items[0].Subtotal,
	}

	// 委托领域服务处理业务逻辑
	return newOrder.ID, s.orderDomainService.CreateOrder(ctx, newOrder)
}

// GetOrder 获取订单
func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*domain_order_core.OrderDO, error) {
	return s.orderDomainService.GetOrderByID(ctx, orderID)
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(ctx context.Context, orderID string) error {
	order, err := s.orderDomainService.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 委托给领域对象处理业务逻辑
	if err := order.Cancel(); err != nil {
		return err
	}

	// 持久化更新
	return s.orderDomainService.UpdateOrder(ctx, order)
}

// PayOrder 支付订单
func (s *OrderService) PayOrder(ctx context.Context, orderID string) error {
	// 1. 获取订单
	orderDO, err := s.orderDomainService.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 2. 业务规则检查：订单必须是已创建的状态
	if orderDO.Status != domain_order_core.OrderStatusCreated {
		return fmt.Errorf("订单状态异常，当前状态: %s,无法发起支付", domain_order_core.GetOrderStatusDetail(orderDO.Status))
	}

	// 3. 检查是否已存在支付单
	existingPayment, err := s.paymentService.GetPaymentByOrderID(ctx, orderDO.ID)
	if err != nil {
		// 仅当支付单不存在时才返回错误，其他错误正常返回
		if !errors.Is(err, domain_payment_core.ErrPaymentNotFound) {
			return err
		}
	}

	// 4. 处理支付单
	var paymentID string
	if existingPayment != nil {
		switch existingPayment.Status {
		case domain_payment_core.PaymentStatusPaid:
			return errors.New("订单已支付，无需重复操作")
		case domain_payment_core.PaymentStatusPending, domain_payment_core.PaymentStatusCreated:
			paymentID = existingPayment.ID
		default:
			return fmt.Errorf("支付单状态异常: %s", domain_payment_core.GetPaymentStatusDetail(existingPayment.Status))
		}
	} else {
		// 创建新支付单
		newPaymentID, err := s.paymentService.CreatePayment(ctx, orderDO.ID, orderDO.TotalAmount, "CNY", 1)
		if err != nil {
			return fmt.Errorf("创建支付单失败: %w", err)
		}
		paymentID = newPaymentID
	}

	// 5. 更新订单状态为待支付
	if err := orderDO.MarkAsPendingPayment(); err != nil {
		return fmt.Errorf("更新订单为待支付状态失败: %w", err)
	}

	// 6. 持久化订单状态变更
	if err := s.orderDomainService.UpdateOrder(ctx, orderDO); err != nil {
		return fmt.Errorf("保存订单状态失败: %w", err)
	}

	// 7. 这里应该调用支付网关获取支付链接或发起支付处理
	// 实际项目中这里会有支付网关的交互逻辑

	if paymentID != "" {
		// 这里占位， 打印下paymentID， 这里还是重新修改
		fmt.Println("支付链接或支付处理信息:", paymentID)
	}

	return nil
}

// UpdateOrder 更新订单
func (s *OrderService) UpdateOrder(ctx context.Context, orderDO *domain_order_core.OrderDO) error {
	if err := s.orderDomainService.UpdateOrder(ctx, orderDO); err != nil {
		// 检查是否为乐观锁冲突错误 (处理乐观锁冲突（v2 特定写法）)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("订单已被其他操作更新，请刷新后重试: %w", err)
		}
		return fmt.Errorf("更新订单失败: %w", err)
	}
	return nil
}
