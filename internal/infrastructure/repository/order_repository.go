package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/vaynedu/ddd_order_example/internal/domain/order"
)

// OrderRepositoryMySQL MySQL实现的订单仓储
type OrderRepositoryMySQL struct {
	db *sqlx.DB
}

// NewOrderRepository 创建订单仓储实例
func NewOrderRepository(db *sqlx.DB) order.OrderRepository {
	return &OrderRepositoryMySQL{db: db}
}

// Save 保存订单
func (r *OrderRepositoryMySQL) Save(ctx context.Context, order *order.Order) error {
	// 开始事务
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 插入或更新订单主表
	query := `
        INSERT INTO orders (id, customer_id, status, total_amount, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
        status = VALUES(status),
        total_amount = VALUES(total_amount),
        updated_at = VALUES(updated_at)
    `

	_, err = tx.ExecContext(ctx, query,
		order.ID,
		order.CustomerID,
		order.Status,
		order.TotalAmount,
		order.CreatedAt,
		order.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// 先删除原有订单项
	_, err = tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id = ?", order.ID)
	if err != nil {
		return err
	}

	// 插入新的订单项
	for _, item := range order.Items {
		query := `
            INSERT INTO order_items (order_id, product_id, quantity, unit_price, subtotal)
            VALUES (?, ?, ?, ?, ?)
        `

		_, err = tx.ExecContext(ctx, query,
			order.ID,
			item.ProductID,
			item.Quantity,
			item.UnitPrice,
			item.Subtotal,
		)

		if err != nil {
			return err
		}
	}

	// 提交事务
	return tx.Commit()
}

// FindByID 根据ID查找订单
func (r *OrderRepositoryMySQL) FindByID(ctx context.Context, id string) (*order.Order, error) {
	// 查询订单主表
	query := `
        SELECT id, customer_id, status, total_amount, created_at, updated_at
        FROM orders
        WHERE id = ?
    `

	var o order.Order
	err := r.db.GetContext(ctx, &o, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}

	// 查询订单项
	query = `
        SELECT product_id, quantity, unit_price, subtotal
        FROM order_items
        WHERE order_id = ?
    `

	var items []order.OrderItem
	err = r.db.SelectContext(ctx, &items, query, id)
	if err != nil {
		return nil, err
	}

	o.Items = items
	return &o, nil
}
