package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vaynedu/ddd_order_example/internal/domain/domain_order_core"
	"github.com/vaynedu/ddd_order_example/internal/domain/domain_product_core"
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

func (s *OrderService) PayOrder(ctx context.Context, orderID string) error {
	// 1. 获取订单
	orderDO, err := s.orderDomainService.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 2. 业务规则检查：订单必须是已创建的状态
	if orderDO.Status != domain_order_core.OrderStatusCreated {
		return errors.New("订单状态异常，无法支付")
	}

	// 3. 调用支付服务创建支付
	_, err = s.paymentService.CreatePayment(ctx, orderDO.ID, orderDO.TotalAmount, "CNY", 1)
	if err != nil {
		return err
	}

	// 4. 更新订单状态
	if err := orderDO.MarkAsPaid(); err != nil {
		return err
	}

	// 5. 持久化订单变更
	return s.orderDomainService.UpdateOrder(ctx, orderDO)
}
