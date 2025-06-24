package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vaynedu/ddd_order_example/internal/application/service"
	"github.com/vaynedu/ddd_order_example/internal/domain/order"
)

// OrderHandler 订单HTTP处理器
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: service}
}

// CreateOrder 创建订单的HTTP处理函数
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求体
	var request struct {
		CustomerID string `json:"customer_id"`
		Items      []struct {
			ProductID string  `json:"product_id"`
			Quantity  int     `json:"quantity"`
			UnitPrice int `json:"unit_price"`
			Subtotal  int `json:"subtotal"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "无效的请求格式: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2. 转换为领域模型
	var items []order.OrderItemDO
	for _, item := range request.Items {
		items = append(items, order.OrderItemDO{	
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
		})
	}

	// 3. 调用应用服务
	orderID, err := h.orderService.CreateOrder(r.Context(), request.CustomerID, items)
	if err != nil {
		http.Error(w, "创建订单失败: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// 4. 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"order_id": orderID,
	})
}

// GetOrder 获取订单的HTTP处理函数
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// 1. 从URL中获取订单ID
	orderID := r.URL.Path[len("/api/orders/"):]
	if orderID == "" {
		http.Error(w, "订单ID不能为空", http.StatusBadRequest)
		return
	}

	// 2. 调用应用服务
	order, err := h.orderService.GetOrder(r.Context(), orderID)
	if err != nil {
		if err.Error() == "订单不存在" {
			http.Error(w, "订单不存在", http.StatusNotFound)
		} else {
			http.Error(w, "获取订单失败: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// 3. 返回订单数据
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":           order.ID,
		"customer_id":  order.CustomerID,
		"status":       order.Status,
		"total_amount": order.TotalAmount,
		"created_at":   order.CreatedAt,
		"updated_at":   order.UpdatedAt,
		"items":        order.Items,
	})
}
