package dto

import (
	"time"

	"github.com/vaynedu/ddd_order_example/internal/domain/order"
	"github.com/vaynedu/ddd_order_example/pkg/dmoney"
)

// CreateOrderRequest 订单创建请求DTO
type CreateOrderRequest struct {
	CustomerID string          `json:"customer_id"`
	Items      []OrderItemRequest `json:"items"`
}

// OrderItemRequest 订单项请求DTO
type OrderItemRequest struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"` // 元
	Subtotal  float64 `json:"subtotal"`   // 元
}

// ToDomain 将DTO转换为领域模型
func (r *CreateOrderRequest) ToDomain() []order.OrderItemDO {
	var items []order.OrderItemDO
	for _, item := range r.Items {
		items = append(items, order.OrderItemDO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: int(dmoney.ConvertFloat64ToCent(item.UnitPrice)),
			Subtotal:  int(dmoney.ConvertFloat64ToCent(item.Subtotal)),
		})
	}
	return items
}

// OrderResponse 订单响应DTO
type OrderResponse struct {
	ID          string            `json:"id"`
	CustomerID  string            `json:"customer_id"`
	Status      string            `json:"status"`
	TotalAmount float64           `json:"total_amount"` // 元
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Items       []OrderItemResponse `json:"items"`
}

// OrderItemResponse 订单项响应DTO
type OrderItemResponse struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"` // 元
	Subtotal  float64 `json:"subtotal"`   // 元
}

// NewOrderResponse 从领域模型创建响应DTO
func NewOrderResponse(order *order.OrderDO) *OrderResponse {
	items := make([]OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = OrderItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: dmoney.ConvertCentToFloat64(int64(item.UnitPrice)),
			Subtotal:  dmoney.ConvertCentToFloat64(int64(item.Subtotal)),
		}
	}

	return &OrderResponse{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      string(order.Status),
		TotalAmount: dmoney.ConvertCentToFloat64(int64(order.TotalAmount)),
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		Items:       items,
	}
}