package order

import (
    "errors"
    "time"
)

// Order 订单聚合根
type Order struct {
    ID          string      `json:"id"`
    CustomerID  string      `json:"customer_id"`
    Items       []OrderItem `json:"items"`
    Status      OrderStatus `json:"status"`
    TotalAmount float64     `json:"total_amount"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

// OrderItem 订单项
type OrderItem struct {
    ProductID   string  `json:"product_id"`
    Quantity    int     `json:"quantity"`
    UnitPrice   float64 `json:"unit_price"`
    Subtotal    float64 `json:"subtotal"`
}

// OrderStatus 订单状态
type OrderStatus string

const (
    OrderStatusCreated   OrderStatus = "created"
    OrderStatusPaid      OrderStatus = "paid"
    OrderStatusShipped   OrderStatus = "shipped"
    OrderStatusCompleted OrderStatus = "completed"
    OrderStatusCancelled OrderStatus = "cancelled"
)

// Validate 创建订单时的业务规则校验
func (o *Order) Validate() error {
    if o.CustomerID == "" {
        return errors.New("客户ID不能为空")
    }
    
    if len(o.Items) == 0 {
        return errors.New("订单商品不能为空")
    }
    
    var calculatedTotal float64
    for _, item := range o.Items {
        if item.ProductID == "" {
            return errors.New("商品ID不能为空")
        }
        
        if item.Quantity <= 0 {
            return errors.New("商品数量必须大于0")
        }
        
        if item.UnitPrice < 0 {
            return errors.New("商品单价不能为负数")
        }
        
        calculatedTotal += item.Subtotal
    }
    
    if o.TotalAmount != calculatedTotal {
        return errors.New("订单总金额与商品小计之和不匹配")
    }
    
    return nil
}

// CanBeCancelled 检查订单是否可以被取消
func (o *Order) CanBeCancelled() bool {
    return o.Status == OrderStatusCreated || o.Status == OrderStatusPaid
}

// Cancel 取消订单的行为方法
func (o *Order) Cancel() error {
    if !o.CanBeCancelled() {
        return errors.New("当前订单状态不允许取消")
    }
    
    o.Status = OrderStatusCancelled
    o.UpdatedAt = time.Now()
    return nil
}  