package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vaynedu/ddd_order_example/internal/domain/domain_order_core"
	"github.com/vaynedu/ddd_order_example/internal/domain/domain_payment_core"
	"github.com/vaynedu/ddd_order_example/internal/domain/domain_product_core"
	"github.com/vaynedu/ddd_order_example/internal/infrastructure/mocks"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// TestOrderService_CreateOrder_Success 创建订单成功场景
func TestOrderService_CreateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建mock依赖
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	mockProductService := mocks.NewMockProductService(ctrl)
	mockPaymentRepo := mocks.NewMockRepository(ctrl)
	mockPaymentProxy := mocks.NewMockPaymentProxy(ctrl)

	// 初始化领域服务和应用服务
	paymentDomainService := domain_payment_core.NewPaymentDomainService(mockPaymentRepo)
	mockPaymentService := NewPaymentService(paymentDomainService, mockPaymentProxy)
	orderDomainService := domain_order_core.NewOrderDomainService(mockOrderRepo)
	service := NewOrderService(orderDomainService, mockPaymentService, mockProductService)

	// 准备测试数据
	ctx := context.Background()
	customerID := "cust_123"
	items := []*domain_order_core.OrderItemDO{
		{
			ProductID: "prod_123",
			Quantity:  2,
			UnitPrice: 100,
			Subtotal:  200,
		},
	}

	// 设置mock预期
	mockProductService.EXPECT().ValidateProduct(gomock.Any(), gomock.Any()).Return(&domain_product_core.ValidateProductResponse{
		IsValid: true,
		Product: &domain_product_core.Product{
			ID:     "prod_123",
			Status: domain_product_core.StatusValid,
		},
	}, nil)

	mockOrderRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

	// 执行测试
	orderID, err := service.CreateOrder(ctx, customerID, items)

	// 验证结果
	assert.NoError(t, err)
	assert.NotEmpty(t, orderID)
}

// TestOrderService_UpdateOrder_OptimisticLockConflict 更新订单乐观锁冲突场景
func TestOrderService_UpdateOrder_OptimisticLockConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建mock依赖
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	orderDomainService := domain_order_core.NewOrderDomainService(mockOrderRepo)
	service := NewOrderService(orderDomainService, nil, nil)

	// 准备测试数据
	ctx := context.Background()
	orderDO := &domain_order_core.OrderDO{
		ID:      "order_123",
		Status:  domain_order_core.OrderStatusCreated,
	}

	// 设置mock预期 - 返回乐观锁冲突错误
	mockOrderRepo.EXPECT().Save(gomock.Any(), orderDO).Return(gorm.ErrDuplicatedKey)

	// 执行测试
	err := service.UpdateOrder(ctx, orderDO)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "订单已被其他操作更新，请刷新后重试")
}

// TestOrderService_GetOrder_Success 获取订单成功场景
func TestOrderService_GetOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建mock依赖
	mockOrderRepo := mocks.NewMockOrderRepository(ctrl)
	orderDomainService := domain_order_core.NewOrderDomainService(mockOrderRepo)
	service := NewOrderService(orderDomainService, nil, nil)

	// 准备测试数据
	ctx := context.Background()
	orderID := "order_123"
	expectedOrder := &domain_order_core.OrderDO{
		ID:      orderID,
		Status:  domain_order_core.OrderStatusCreated,
	}

	// 设置mock预期
	mockOrderRepo.EXPECT().FindByID(ctx, orderID).Return(expectedOrder, nil)

	// 执行测试
	result, err := service.GetOrder(ctx, orderID)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, result)
}

// MockProductService 模拟ProductService接口
type MockProductService struct {
	ctrl *gomock.Controller
}

func NewMockProductService(ctrl *gomock.Controller) *MockProductService {
	return &MockProductService{ctrl: ctrl}
}

func (m *MockProductService) ValidateProduct(ctx context.Context, req *domain_product_core.ValidateProductRequest) (*domain_product_core.ValidateProductResponse, error) {
	panic("implement me")
}

// MockPaymentService 模拟PaymentService结构体
type MockPaymentService struct {
	ctrl *gomock.Controller
}

func NewMockPaymentService(ctrl *gomock.Controller) *MockPaymentService {
	return &MockPaymentService{ctrl: ctrl}
}

func (m *MockPaymentService) GetPaymentByOrderID(ctx context.Context, orderID string) (*domain_payment_core.PaymentDO, error) {
	return nil, domain_payment_core.ErrPaymentNotFound
}

func (m *MockPaymentService) CreatePayment(ctx context.Context, orderID string, amount int64, currency string, payType int) (string, error) {
	return "pay_123", nil
}
