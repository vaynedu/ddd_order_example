package order

import (
    "context"
    "errors"
    "time"
)

// OrderDomainService 订单领域服务
type OrderDomainService struct {
    orderRepo OrderRepository
}

// NewOrderDomainService 创建订单领域服务
func NewOrderDomainService(repo OrderRepository) *OrderDomainService {
    return &OrderDomainService{orderRepo: repo}
}

// CreateOrder 创建订单
func (s *OrderDomainService) CreateOrder(ctx context.Context, order *OrderDO) error {
    // 应用业务规则
    if err := order.Validate(); err != nil {
        return err
    }
    
    // 设置默认值
    order.Status = OrderStatusCreated
    order.CreatedAt = time.Now()
    order.UpdatedAt = order.CreatedAt
    
    // 持久化订单
    return s.orderRepo.Save(ctx, order)
}

// PayOrder 支付订单
func (s *OrderDomainService) PayOrder(ctx context.Context, orderID string) error {
    order, err := s.orderRepo.FindByID(ctx, orderID)
    if err != nil {
        return err
    }
    
    if order.Status != OrderStatusCreated {
        return errors.New("只有已创建的订单可以支付")
    }
    
    order.Status = OrderStatusPaid
    order.UpdatedAt = time.Now()
    
    return s.orderRepo.Save(ctx, order)
}

// GetOrderByID 获取订单
func (s *OrderDomainService) GetOrderByID(ctx context.Context, orderID string) (*OrderDO, error) {
    return s.orderRepo.FindByID(ctx, orderID)
}

// UpdateOrder 更新订单
func (s *OrderDomainService) UpdateOrder(ctx context.Context, order *OrderDO) error {
    return s.orderRepo.Save(ctx, order)
}