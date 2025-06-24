package repository

import (
	"context"
	"errors"

	"github.com/vaynedu/ddd_order_example/internal/domain/order"
	"gorm.io/gorm"
)

// OrderRepositoryMySQL MySQL实现的订单仓储
type OrderRepositoryMySQL struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储实例
func NewOrderRepository(db *gorm.DB) order.OrderRepository {
	return &OrderRepositoryMySQL{db: db}
}

// Save 保存订单
func (r *OrderRepositoryMySQL) Save(ctx context.Context, o *order.OrderDO) error {
	// 开始事务
	tx := r.db.Begin()
	tx = tx.WithContext(ctx)
	if err := tx.Error; err != nil {
		return err
	}
	defer tx.Rollback()

	// 使用GORM保存订单主表
	if err := tx.Table("t_order").Create(o).Error; err != nil {
		return err
	}

	// 删除原有订单项
	if err := tx.Table("t_order_items").Where("order_id = ?", o.ID).Delete(&order.OrderItemDO{}).Error; err != nil {
		return err
	}

	// 批量插入新订单项
	orderItems := make([]order.OrderItemDO, len(o.Items))
	for i, item := range o.Items {
		orderItems[i] = order.OrderItemDO{
			OrderID:   o.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
		}
	}
	if err := tx.Table("t_order_items").Create(&orderItems).Error; err != nil {
		return err
	}

	// 提交事务
	err := tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

// FindByID 根据ID查找订单
func (r *OrderRepositoryMySQL) FindByID(ctx context.Context, id string) (*order.OrderDO, error) {
	// 查询订单主表
	var o order.OrderDO
	if err := r.db.WithContext(ctx).Table("t_order").First(&o, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}

	// 查询订单项
	query := `
        SELECT product_id, quantity, unit_price, subtotal
        FROM order_items
        WHERE order_id = ?
    `

	var items []order.OrderItemDO
	if err := r.db.WithContext(ctx).Table(order.OrderItemDO{}.TableName()).Raw(query, id).Scan(&items).Error; err != nil {
		return nil, err
	}

	o.Items = items
	return &o, nil
}
