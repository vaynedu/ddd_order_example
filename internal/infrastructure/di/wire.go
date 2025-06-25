//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/vaynedu/ddd_order_example/internal/application/service"
	"github.com/vaynedu/ddd_order_example/internal/domain/domain_product_core"
	"github.com/vaynedu/ddd_order_example/internal/domain/order"
	"github.com/vaynedu/ddd_order_example/internal/infrastructure/external/mocks"
	"github.com/vaynedu/ddd_order_example/internal/infrastructure/repository"
	"github.com/vaynedu/ddd_order_example/internal/interface/handler"
	"gorm.io/gorm"
)

// 生产环境依赖注入
// func InitializeOrderHandler(db *gorm.DB) (*handler.OrderHandler, error) {
// 	wire.Build(
// 		repository.NewOrderRepository,
// 		NewProductService,
// 		service.NewOrderService,
// 		handler.NewOrderHandler,
// 	)
// 	return nil, nil
// }

// 测试环境依赖注入 - 使用Mock商品服务
func InitializeTestOrderHandler(db *gorm.DB) (*handler.OrderHandler, error) {
	wire.Build(
		NewOrderRepository,
		NewOrderDomainService,
		NewMockProductService,
		NewOrderService,
		NewOrderHandler,
	)
	return nil, nil
}

// NewOrderRepository - 初始化仓储
func NewOrderRepository(db *gorm.DB) order.OrderRepository {
	return repository.NewOrderRepository(db)
}

// NewOrderDomainService - 初始化领域服务
func NewOrderDomainService(repo order.OrderRepository) order.OrderDomainService {
	return order.NewOrderDomainService(repo)
}

// NewOrderService 初始化应用服务
// 传入实例化的order，传入实例化的productService
func NewOrderService(domainService order.OrderDomainService, productService domain_product_core.ProductService) *service.OrderService {
	return service.NewOrderService(domainService, productService)
}

// NewOrderHandler 初始化处理器
func NewOrderHandler(orderService *service.OrderService) *handler.OrderHandler {
	return handler.NewOrderHandler(orderService)
}

// NewMockProductService 创建商品服务的Mock实现
func NewMockProductService() domain_product_core.ProductService {
	return mocks.NewMockProductService()
}
