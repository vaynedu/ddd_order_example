package order

import (
	"errors"
	"time"
)

// OrderDO 订单聚合根
// 修正：将 gorm 标签移至结构体定义处

// TableName 指定模型对应的数据库表名
func (OrderDO) TableName() string {
	return "t_order"
}

type OrderDO struct {
	ID          string        `json:"id" gorm:"column:id"`
	CustomerID  string        `json:"customer_id" gorm:"column:customer_id"`
	Items       []OrderItemDO `json:"items" gorm:"foreignKey:OrderID"`
	Status      OrderStatus   `json:"status" gorm:"column:status"`
	TotalAmount int64           `json:"total_amount" gorm:"column:total_amount"`
	CreatedAt   time.Time     `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time     `json:"updated_at" gorm:"column:updated_at"`
}

// OrderItemDOs 订单项集合
type OrderItemDOs []OrderItemDO

// OrderItemDO 订单项
// 将 gorm 标签移至 OrderItemDO 结构体定义处

// TableName 指定模型对应的数据库表名
func (OrderItemDO) TableName() string {
	return "t_order_items"
}

type OrderItemDO struct {
	OrderID   string `json:"order_id" gorm:"column:order_id"`
	ProductID string `json:"product_id" gorm:"column:product_id"`
	Quantity  int64    `json:"quantity" gorm:"column:quantity"`
	UnitPrice int64    `json:"unit_price" gorm:"column:unit_price"`
	Subtotal  int64    `json:"subtotal" gorm:"column:subtotal"`
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
func (o *OrderDO) Validate() error {
	if o.CustomerID == "" {
		return errors.New("客户ID不能为空")
	}

	if len(o.Items) == 0 {
		return errors.New("订单商品不能为空")
	}

	var calculatedTotal int64
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

	if o.TotalAmount != calculatedTotal{
		return errors.New("订单总金额与商品小计之和不匹配")
	}

	return nil
}

// CanBeCancelled 检查订单是否可以被取消
func (o *OrderDO) CanBeCancelled() bool {
	return o.Status == OrderStatusCreated || o.Status == OrderStatusPaid
}

// Cancel 取消订单的行为方法
func (o *OrderDO) Cancel() error {
	if !o.CanBeCancelled() {
		return errors.New("当前订单状态不允许取消")
	}

	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}
