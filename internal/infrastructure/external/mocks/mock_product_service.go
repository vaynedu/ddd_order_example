package mocks

import (
	"context"

	"github.com/vaynedu/ddd_order_example/internal/domain/domain_product_core"
)

// MockProductService 商品服务的mock实现
type MockProductService struct {
	// ValidateFunc 自定义验证函数，用于测试不同场景
	ValidateFunc func(ctx context.Context, req *domain_product_core.ValidateProductRequest) (*domain_product_core.ValidateProductResponse, error)
}

// NewMockProductService 创建mock实例
func NewMockProductService() *MockProductService {
	return &MockProductService{}
}

// ValidateProduct 实现ProductService接口
func (m *MockProductService) ValidateProduct(ctx context.Context, req *domain_product_core.ValidateProductRequest) (*domain_product_core.ValidateProductResponse, error) {
	if m.ValidateFunc != nil {
		return m.ValidateFunc(ctx, req)
	}

	// 默认返回有效的商品信息
	return &domain_product_core.ValidateProductResponse{
		Product: &domain_product_core.Product{
			ID:     req.ProductID,
			Name:   req.Name,
			Price:  req.Price,
			Status: domain_product_core.StatusValid,
		},
		IsValid:  true,
		Messages: "",
	}, nil
}

// WithProductNotFound 设置商品不存在场景
func (m *MockProductService) WithProductNotFound() {
	m.ValidateFunc = func(ctx context.Context, req *domain_product_core.ValidateProductRequest) (*domain_product_core.ValidateProductResponse, error) {
		return &domain_product_core.ValidateProductResponse{
			Product:  nil,
			IsValid:  false,
			Messages: "product not found",
		}, nil
	}
}

// WithInvalidProduct 设置商品无效场景
func (m *MockProductService) WithInvalidProduct() {
	m.ValidateFunc = func(ctx context.Context, req *domain_product_core.ValidateProductRequest) (*domain_product_core.ValidateProductResponse, error) {
		return &domain_product_core.ValidateProductResponse{
			Product: &domain_product_core.Product{
				ID:     req.ProductID,
				Name:   req.Name,
				Price:  req.Price,
				Status: domain_product_core.StatusInvalid,
			},
			IsValid:  false,
			Messages: "product is invalid",
		}, nil
	}
}
