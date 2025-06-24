package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/vaynedu/ddd_order_example/internal/domain/order"
)

// OrderService 订单应用服务
type OrderService struct {
	orderDomainService *order.OrderDomainService
}

// NewOrderService 创建订单应用服务
func NewOrderService(repo order.OrderRepository) *OrderService {
	return &OrderService{
		orderDomainService: order.NewOrderDomainService(repo),
	}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, customerID string, items []order.OrderItem) (string, error) {
	// 计算总金额
	var totalAmount float64
	for _, item := range items {
		totalAmount += item.Subtotal
	}

	// 创建订单聚合
	newOrder := &order.Order{
		ID:          uuid.New().String(),
		CustomerID:  customerID,
		Items:       items,
		TotalAmount: totalAmount,
	}

	// 委托给领域服务处理业务逻辑
	return newOrder.ID, s.orderDomainService.CreateOrder(ctx, newOrder)
}

// GetOrder 获取订单
func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*order.Order, error) {
	return s.orderDomainService.orderRepo.FindByID(ctx, orderID)
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(ctx context.Context, orderID string) error {
	order, err := s.orderDomainService.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 委托给领域对象处理业务逻辑
	if err := order.Cancel(); err != nil {
		return err
	}

	// 持久化更新
	return s.orderDomainService.orderRepo.Save(ctx, order)
}
