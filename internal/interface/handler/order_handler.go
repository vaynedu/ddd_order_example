package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vaynedu/ddd_order_example/internal/application/service"
	"github.com/vaynedu/ddd_order_example/internal/interface/dto"
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
	var request dto.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "无效的请求格式: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2. 转换为领域模型（通过DTO）
	items := request.ToDomain()

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
	// 1. 从body中获取订单id
	var request struct {
		OrderID string `json:"order_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "无效的请求格式: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2. 调用应用服务
	order, err := h.orderService.GetOrder(r.Context(), request.OrderID)
	if err != nil {
		if err.Error() == "订单不存在" {
			http.Error(w, "订单不存在", http.StatusNotFound)
		} else {
			http.Error(w, "获取订单失败: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// 3. 转换为响应DTO并返回
	response := dto.NewOrderResponse(order)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
