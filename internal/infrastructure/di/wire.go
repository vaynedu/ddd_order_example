//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/vaynedu/ddd_order_example/internal/application/service"
	"github.com/vaynedu/ddd_order_example/internal/domain/domain_order_core"
	"github.com/vaynedu/ddd_order_example/internal/domain/domain_payment_core"
	"github.com/vaynedu/ddd_order_example/internal/domain/domain_product_core"
	"github.com/vaynedu/ddd_order_example/internal/infrastructure/external/mocks"
	"github.com/vaynedu/ddd_order_example/internal/infrastructure/payment"
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
		NewOrderRepository,    // 订单仓储
		NewOrderDomainService, // 订单领域服务
		NewMockProductService, // 商品服务

		NewPaymentRepository,    // 支付仓储
		NewPaymentDomainService, // 支付领域服务
		NewMockPaymentProxy,     // 支付代理
		NewPaymentService,       // 支付应用服务
		NewOrderService,
		NewOrderHandler,
	)
	return nil, nil
}

// NewOrderRepository - 初始化仓储
func NewOrderRepository(db *gorm.DB) domain_order_core.OrderRepository {
	return repository.NewOrderRepository(db)
}

// NewOrderDomainService - 初始化领域服务
func NewOrderDomainService(repo domain_order_core.OrderRepository) domain_order_core.OrderDomainService {
	return domain_order_core.NewOrderDomainService(repo)
}

// NewOrderHandler 初始化处理器
func NewOrderHandler(orderService *service.OrderService) *handler.OrderHandler {
	return handler.NewOrderHandler(orderService)
}

// NewMockProductService 创建商品服务的Mock实现
func NewMockProductService() domain_product_core.ProductService {
	return mocks.NewMockProductService()
}

// NewPaymentRepository 创建支付仓储
func NewPaymentRepository(db *gorm.DB) domain_payment_core.Repository {
	return repository.NewPaymentRepository(db)
}

// NewPaymentDomainService 创建支付领域服务
func NewPaymentDomainService(repo domain_payment_core.Repository) *domain_payment_core.PaymentDomainService {
	return domain_payment_core.NewPaymentDomainService(repo)
}

// // NewPaymentProxy 创建真实支付代理， 这里可创建函数，封装不同的支付方式，待定
// func NewPaymentProxy() payment.PaymentProxy {
// 	return payment.NewRealPaymentProxy()
// }

// NewMockPaymentProxy 创建Mock支付代理
func NewMockPaymentProxy() payment.PaymentProxy {
	return payment.NewMockPaymentProxy()
}

// NewPaymentService 创建支付应用服务
func NewPaymentService(domainService *domain_payment_core.PaymentDomainService, proxy payment.PaymentProxy) *service.PaymentService {
	return service.NewPaymentService(domainService, proxy)
}

// NewOrderService 创建订单应用服务, 包含订单领域服务, 支付应用服务, 商品服务
func NewOrderService(
	productService domain_product_core.ProductService,
	orderDomainService domain_order_core.OrderDomainService,
	paymentService *service.PaymentService,
) *service.OrderService {
	return service.NewOrderService(orderDomainService, paymentService, productService)
}
