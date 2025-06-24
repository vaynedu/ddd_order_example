-- 创建订单数据库
CREATE DATABASE IF NOT EXISTS orders;
USE orders;

-- 创建订单表
-- 订单主表，存储订单基本信息，与订单项表(t_order_items)为一对多关系
CREATE TABLE IF NOT EXISTS t_order (
    id VARCHAR(36) PRIMARY KEY,
    customer_id VARCHAR(36) NOT NULL,
    status ENUM('created', 'paid', 'shipped', 'completed', 'cancelled') NOT NULL,
    total_amount INT NOT NULL COMMENT '订单总金额，单位：分',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_customer_id (customer_id),
    INDEX idx_status (status)
);

-- 创建订单项表
-- 订单项子表，存储订单包含的商品信息，通过order_id与订单主表关联
CREATE TABLE IF NOT EXISTS t_order_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL COMMENT '关联订单主表的ID',
    product_id VARCHAR(36) NOT NULL,
    quantity INT NOT NULL,
    unit_price INT NOT NULL COMMENT '商品单价，单位：分',
    subtotal INT NOT NULL COMMENT '商品小计金额，单位：分',
    INDEX idx_order_id (order_id)
);