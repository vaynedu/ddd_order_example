package domain_product_core

import "context"

// ValidateProductRequest 商品验证请求参数
type ValidateProductRequest struct {
	ProductID string
	Name      string
	Price     int64
	Quantity  int64
}

// ValidateProductResponse 商品验证响应结果
type ValidateProductResponse struct {
	Product  *Product
	IsValid  bool
	Messages string

	// 后续考虑新增全局错误码
}

// ProductService 商品服务抽象接口
type ProductService interface {
	// ValidateProduct 验证商品是否合法可用
	ValidateProduct(ctx context.Context, req *ValidateProductRequest) (*ValidateProductResponse, error)
}
