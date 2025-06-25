//go:build wireinject
// +build wireinject

//go:generate wire
package di

import (
	"github.com/google/wire"
	"github.com/vaynedu/ddd_order_example/internal/application/service"
	"github.com/vaynedu/ddd_order_example/internal/domain/order"
	"github.com/vaynedu/ddd_order_example/internal/infrastructure/repository"
	"github.com/vaynedu/ddd_order_example/internal/interface/handler"
	"gorm.io/gorm"
)

// InitializeOrderHandler 依赖注入入口
func InitializeOrderHandler(db *gorm.DB) (*handler.OrderHandler, error) {
	wire.Build(NewOrderRepository, NewOrderService, NewOrderHandler)
	return nil, nil
}

// NewOrderRepository - 初始化仓储
func NewOrderRepository(db *gorm.DB) order.OrderRepository {
	return repository.NewOrderRepository(db)
}

// NewOrderService 初始化应用服务
func NewOrderService(repo order.OrderRepository) *service.OrderService {
	return service.NewOrderService(repo)
}

// NewOrderHandler 初始化处理器
func NewOrderHandler(orderService *service.OrderService) *handler.OrderHandler {
	return handler.NewOrderHandler(orderService)
}
