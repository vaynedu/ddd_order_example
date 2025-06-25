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