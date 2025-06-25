package product_api

import (
	"context"
	"errors"

	"github.com/vaynedu/ddd_order_example/internal/domain/domain_product_core"
)

// ThirdPartyProductResponse 第三方API响应结构
type ThirdPartyProductResponse struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	// todo 先简单实现，理论上最好有公共库，表示商品的状态，这里暂时使用0 ,1 ,2
	Status int `json:"status"` //
	// 其他第三方字段...
}

// ProductServiceAdapter 适配第三方API到领域接口
type ProductServiceAdapter struct {
	client *ThirdPartyProductAPI
}

// NewProductServiceAdapter 创建适配器实例
func NewProductServiceAdapter(client *ThirdPartyProductAPI) domain_product_core.ProductService {
	return &ProductServiceAdapter{
		client: client,
	}
}

// ValidateProduct 实现领域接口
func (a *ProductServiceAdapter) ValidateProduct(ctx context.Context, req *domain_product_core.ValidateProductRequest) (*domain_product_core.ValidateProductResponse, error) {
	resp, err := a.client.GetProductStatus(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}
	// 名称校验
	if resp.Name == "" {
		return nil, errors.New("product name is empty")
	}
	// 状态校验
	if resp.Status != 0 {
		return nil, errors.New("product status is invalid")
	}
	// 价格校验
	if resp.Price != req.Price {
		return nil, errors.New("product price is invalid")
	}
	// 其他校验...

	// 转换第三方响应为领域模型
	domainProduct := &domain_product_core.Product{
		ID:    resp.ProductID,
		Name:  resp.Name,
		Price: resp.Price,
	}

	// 状态映射
	switch {
	case resp.Status == 1:
		domainProduct.Status = domain_product_core.StatusDeleted
	case resp.Status == 2:
		domainProduct.Status = domain_product_core.StatusInvalid
	default:
		domainProduct.Status = domain_product_core.StatusValid
	}

	// todo message信息待处理
	message := "待处理"

	return &domain_product_core.ValidateProductResponse{
		Product:  domainProduct,
		IsValid:  domainProduct.Status == 0,
		Messages: message,
	}, nil
}
