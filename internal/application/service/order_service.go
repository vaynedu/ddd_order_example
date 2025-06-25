package service

import (
	"context"
	"errors"

	"github.com/vaynedu/ddd_order_example/internal/domain/domain_product_core"
	"github.com/vaynedu/ddd_order_example/internal/domain/order"
)

type OrderService struct {
	productService     domain_product_core.ProductService // 依赖商品领域接口
	orderDomainService order.OrderDomainService           // 依赖订单领域服务
}

func NewOrderService(orderDomainService order.OrderDomainService, productService domain_product_core.ProductService) *OrderService {
	return &OrderService{
		orderDomainService: orderDomainService,
		productService:     productService,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, orderID string, items []*order.OrderItemDO) (string, error) {
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
	newOrder := &order.OrderDO{
		ID:     orderID,
		Status: order.OrderStatusCreated,
		Items: []order.OrderItemDO{
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
func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*order.OrderDO, error) {
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
