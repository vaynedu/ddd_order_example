### DDD架构下第三方商品接口集成方案

#### 一、设计思想与原则
1. **依赖倒置原则**：领域层定义抽象接口，基础设施层实现具体调用
2. **防腐层模式**：隔离外部系统变化对核心业务逻辑的影响
3. **接口适配**：将第三方接口数据结构转换为领域模型
4. **关注点分离**：业务规则与外部通信逻辑解耦

#### 二、目录结构设计
```
internal/
├── domain/
│   └── product/
│       ├── model.go        // 商品领域模型
│       └── service.go      // 商品服务抽象接口
├── infrastructure/
│   └── external/
│       └── product_api/
│           ├── adapter.go  // 数据转换适配器
│           ├── client.go   // 第三方API客户端
│           └── service.go  // 领域接口实现
└── application/
    └── service/
        └── order_service.go // 使用商品服务
```

#### 三、核心代码实现

##### 1. 领域层 - 抽象接口定义
```go:/c:/Users/lingze/.go/src/ddd_order_example/internal/domain/product/service.go
package product

import "context"

// Product 商品领域模型
type Product struct {
    ID        string
    Name      string
    Status    ProductStatus
    Price     float64
    // 其他领域属性...
}

// ProductStatus 商品状态枚举
type ProductStatus int

const (
    StatusValid   ProductStatus = iota // 有效
    StatusInvalid                      // 无效
    StatusDeleted                      // 已删除
)

// ProductService 商品服务抽象接口
type ProductService interface {
    // ValidateProduct 验证商品是否合法可用
    ValidateProduct(ctx context.Context, productID string) (*Product, error)
}
```

##### 2. 基础设施层 - 第三方接口实现
```go:/c:/Users/lingze/.go/src/ddd_order_example/internal/infrastructure/external/product_api/client.go
package product_api

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
)

// ThirdPartyProductAPI 第三方商品API客户端
type ThirdPartyProductAPI struct {
    baseURL    string
    timeout    time.Duration
    httpClient *http.Client
    apiKey     string
}

// NewThirdPartyProductAPI 创建客户端实例
func NewThirdPartyProductAPI(baseURL, apiKey string) *ThirdPartyProductAPI {
    return &ThirdPartyProductAPI{
        baseURL: baseURL,
        timeout: 5 * time.Second,
        httpClient: &http.Client{
            Timeout: 5 * time.Second,
        },
        apiKey: apiKey,
    }
}

// GetProductStatus 调用第三方API获取商品状态
func (c *ThirdPartyProductAPI) GetProductStatus(ctx context.Context, productID string) (*ThirdPartyProductResponse, error) {
    // 实现HTTP请求逻辑...
    req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/products/"+productID, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+c.apiKey)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result ThirdPartyProductResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    return &result, nil
}
```

##### 3. 基础设施层 - 数据转换适配器
```go:/c:/Users/lingze/.go/src/ddd_order_example/internal/infrastructure/external/product_api/adapter.go
package product_api

import (
    "context"
    "ddd_order_example/internal/domain/product"
    "errors"
)

// ThirdPartyProductResponse 第三方API响应结构
type ThirdPartyProductResponse struct {
    ProductID string `json:"product_id"`
    Name      string `json:"name"`
    IsActive  bool   `json:"is_active"`
    IsDeleted bool   `json:"is_deleted"`
    Price     string `json:"price"`
    // 其他第三方字段...
}

// ProductServiceAdapter 适配第三方API到领域接口
type ProductServiceAdapter struct {
    client *ThirdPartyProductAPI
}

// NewProductServiceAdapter 创建适配器实例
func NewProductServiceAdapter(client *ThirdPartyProductAPI) product.ProductService {
    return &ProductServiceAdapter{
        client: client,
    }
}

// ValidateProduct 实现领域接口
func (a *ProductServiceAdapter) ValidateProduct(ctx context.Context, productID string) (*product.Product, error) {
    resp, err := a.client.GetProductStatus(ctx, productID)
    if err != nil {
        return nil, err
    }

    // 转换第三方响应为领域模型
    domainProduct := &product.Product{
        ID:   resp.ProductID,
        Name: resp.Name,
    }

    // 状态映射
    switch {
    case resp.IsDeleted:
        domainProduct.Status = product.StatusDeleted
    case !resp.IsActive:
        domainProduct.Status = product.StatusInvalid
    default:
        domainProduct.Status = product.StatusValid
    }

    // 价格转换和验证
    // ...

    return domainProduct, nil
}
```

##### 4. 依赖注入配置
```go:/c:/Users/lingze/.go/src/ddd_order_example/internal/infrastructure/di/providers.go
package di

import (
    "ddd_order_example/internal/infrastructure/external/product_api"
    "os"
)

// 提供第三方商品API客户端
func NewProductAPIClient() *product_api.ThirdPartyProductAPI {
    return product_api.NewThirdPartyProductAPI(
        os.Getenv("PRODUCT_API_URL"),
        os.Getenv("PRODUCT_API_KEY"),
    )
}

// 提供商品服务适配器
func NewProductService(client *product_api.ThirdPartyProductAPI) product.ProductService {
    return product_api.NewProductServiceAdapter(client)
}
```

#### 四、使用示例（订单服务中）
```go:/c:/Users/lingze/.go/src/ddd_order_example/internal/application/service/order_service.go
package service

import (
    "context"
    "ddd_order_example/internal/domain/order"
    "ddd_order_example/internal/domain/product"
    "errors"
)

type OrderService struct {
    orderRepo     order.Repository
    productService product.ProductService // 依赖领域接口
}

func NewOrderService(orderRepo order.Repository, productService product.ProductService) *OrderService {
    return &OrderService{
        orderRepo:     orderRepo,
        productService: productService,
    }
}

func (s *OrderService) CreateOrder(ctx context.Context, orderID string, productID string) error {
    // 验证商品状态
    prod, err := s.productService.ValidateProduct(ctx, productID)
    if err != nil {
        return err
    }

    if prod.Status != product.StatusValid {
        return errors.New("product is not available")
    }

    // 创建订单逻辑...
    return s.orderRepo.Save(ctx, &order.Order{
        ID:        orderID,
        ProductID: productID,
        // 其他订单属性...
    })
}
```

#### 五、设计优势
1. **隔离性**：第三方接口变化仅影响基础设施层实现，不影响领域逻辑
2. **可测试性**：领域层可使用mock实现进行单元测试
3. **可替换性**：更换第三方服务只需实现新的适配器，无需修改业务代码
4. **关注点分离**：领域层专注业务规则，基础设施层处理技术细节
5. **符合依赖规则**：内层不依赖外层，通过接口抽象实现依赖注入

#### 六、扩展建议
1. 添加缓存层减少第三方API调用
2. 实现熔断器模式处理API调用失败
3. 添加请求重试机制提高可靠性
4. 配置化第三方服务地址和超时参数
5. 实现批量商品验证接口提高性能
        