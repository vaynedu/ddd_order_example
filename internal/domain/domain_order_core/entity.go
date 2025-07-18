package domain_order_core

import (
	"errors"
	"time"

	"gorm.io/plugin/optimisticlock"
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
	TotalAmount int64         `json:"total_amount" gorm:"column:total_amount"`
	CreatedAt   time.Time     `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time     `json:"updated_at" gorm:"column:updated_at"`
	// Version     int64         `json:"version" gorm:"column:version;optimistic_lock"` // 乐观锁版本号
	Version optimisticlock.Version `json:"version" gorm:"column:version;optimistic_lock"` // 乐观锁版本号
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
	Quantity  int64  `json:"quantity" gorm:"column:quantity"`
	UnitPrice int64  `json:"unit_price" gorm:"column:unit_price"`
	Subtotal  int64  `json:"subtotal" gorm:"column:subtotal"`
}

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusUnknown   OrderStatus = "unknown"
	OrderStatusCreated   OrderStatus = "created"
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

func GetOrderStatusDetail(status OrderStatus) string {
	switch status {
	case OrderStatusCreated:
		return "已创建"
	case OrderStatusPending:
		return "待支付"
	case OrderStatusPaid:
		return "已支付"
	case OrderStatusShipped:
		return "已发货"
	case OrderStatusCompleted:
		return "已完成"
	case OrderStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

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

	if o.TotalAmount != calculatedTotal {
		return errors.New("订单总金额与商品小计之和不匹配")
	}

	return nil
}

// ValidateUpdate 更新订单
func (o *OrderDO) ValidateUpdate() error {
	if o.ID == "" {
		return errors.New("订单ID不能为空")
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

	if o.TotalAmount != calculatedTotal {
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

// 标记订单为待支付状态
func (o *OrderDO) MarkAsPendingPayment() error {
	// 状态验证：只能从CREATED状态转为PENDING_PAYMENT
	if o.Status != OrderStatusCreated {
		return errors.New("只有已创建的订单可以标记为待支付")
	}

	// 更新订单状态和支付ID
	o.Status = OrderStatusPending
	o.UpdatedAt = time.Now()

	return nil
}

// MarkAsPaid 标记订单为待支付状态
func (o *OrderDO) MarkAsPaid() error {
	// 状态验证：只能从PENDING_PAYMENT状态转为PAID
	if o.Status != OrderStatusPending {
		return errors.New("只有待支付的订单可以标记为已支付")
	}

	// 验证是否有关联的支付ID
	// 订单表没必要关联支付ID， 因为一个订单可能有多次支付

	o.Status = OrderStatusPaid
	o.UpdatedAt = time.Now()
	return nil
}

// CalculateTotalAmount 计算订单总金额
func (o *OrderDO) CalculateTotalAmount() error {
	var totalAmount int64
	for _, item := range o.Items {
		totalAmount += item.Subtotal
	}

	o.TotalAmount = totalAmount
	return nil
}
