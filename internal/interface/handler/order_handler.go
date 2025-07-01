package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/vaynedu/ddd_order_example/internal/application/service"
	"github.com/vaynedu/ddd_order_example/internal/domain/domain_order_core"
	"github.com/vaynedu/ddd_order_example/internal/interface/dto"
	"gorm.io/gorm"
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
	var req dto.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求格式: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2. 转换为领域模型（通过DTO）
	items := req.ToDomain()

	// 3. 调用应用服务
	orderID, err := h.orderService.CreateOrder(r.Context(), req.CustomerID, items)
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

// PayOrder 处理订单支付请求
func (h *OrderHandler) PayOrder(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求体获取订单ID
	var req struct {
		OrderID string `json:"order_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求格式: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.OrderID == "" {
		http.Error(w, "订单ID不能为空", http.StatusBadRequest)
		return
	}

	// 2. 调用应用服务执行支付
	if err := h.orderService.PayOrder(r.Context(), req.OrderID); err != nil {
		http.Error(w, "支付失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "支付成功",
	})
}

// UpdateOrder 更新订单的HTTP处理函数
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求体
	var req dto.UpdateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求格式: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2. 验证订单是否存在
	// todo 关于error，和返回值统一处理
	existingOrder, err := h.orderService.GetOrder(r.Context(), req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "订单不存在", http.StatusNotFound)
		} else {
			http.Error(w, "查询订单失败: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// 3. 合并更新数据（保留原有必要字段并重新计算金额）
	orderDO := req.ToDomain()
	if orderDO.CustomerID == "" {
		orderDO.CustomerID = existingOrder.CustomerID
	}
	orderDO.CreatedAt = existingOrder.CreatedAt
	orderDO.UpdatedAt = time.Now()
	// 状态变更
	if orderDO.Status != domain_order_core.OrderStatusUnknown {
		orderDO.Status = existingOrder.Status
	}

	// 如果更新了订单项，则重新计算总金额
	if len(req.Items) > 0 {
		// 调用领域层方法计算总金额
		if err := orderDO.CalculateTotalAmount(); err != nil {
			http.Error(w, "计算订单金额失败: "+err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		// 未更新订单项，保留原金额
		orderDO.TotalAmount = existingOrder.TotalAmount
	}

	// 4. 调用应用服务
	if err := h.orderService.UpdateOrder(r.Context(), orderDO); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			http.Error(w, "订单已被其他操作更新，请刷新后重试", http.StatusConflict)
		} else {
			http.Error(w, "更新订单失败: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// 5. 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "订单更新成功",
	})
}
