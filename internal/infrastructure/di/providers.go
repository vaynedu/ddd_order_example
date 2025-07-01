package di

import (
	"os"

	"github.com/vaynedu/ddd_order_example/internal/domain/domain_product_core"
	"github.com/vaynedu/ddd_order_example/internal/infrastructure/external/product_api"
)

// 提供第三方商品API客户端
func NewProductAPIClient() *product_api.ThirdPartyProductAPI {
	return product_api.NewThirdPartyProductAPI(
		os.Getenv("PRODUCT_API_URL"),
		os.Getenv("PRODUCT_API_KEY"),
	)
}

// NewProductService 创建商品服务实例
func NewProductService(client *product_api.ThirdPartyProductAPI) domain_product_core.ProductService {
	return product_api.NewProductServiceAdapter(client)
}
