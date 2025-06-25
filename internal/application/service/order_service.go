package service

import (
	"context"
	"errors"
	"fmt"

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
func (s *OrderService) CreateOrder(ctx context.Context, customerID string, items []order.OrderItemDO) (string, error) {
	// 计算总金额并验证订单项
	var totalAmount int
	for _, item := range items {
		if item.ProductID == "" {
			return "", errors.New("product ID cannot be empty")
		}
		if item.Quantity <= 0 {
			return "", errors.New("quantity must be greater than zero")
		}
		if item.UnitPrice <= 0 {
			return "", errors.New("unit price must be greater than zero")
		}
		calculatedSubtotal := item.Quantity * item.UnitPrice
		if item.Subtotal != calculatedSubtotal {
			return "", fmt.Errorf("subtotal mismatch for product %s: expected %d, got %d", item.ProductID, calculatedSubtotal, item.Subtotal)
		}
		totalAmount += item.Subtotal
	}

	// 创建订单聚合
	newOrder := &order.OrderDO{
		ID:          uuid.New().String(),
		CustomerID:  customerID,
		Items:       items,
		TotalAmount: totalAmount,
		Status:      order.OrderStatusCreated,
	}

	// 委托给领域服务处理业务逻辑
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
