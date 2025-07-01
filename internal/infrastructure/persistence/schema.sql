-- 创建订单数据库
CREATE DATABASE IF NOT EXISTS orders;
USE orders;

-- 删除表
drop table t_order;
drop table t_order_items;
drop table t_payment;

-- 创建订单表
-- 订单主表，存储订单基本信息，与订单项表(t_order_items)为一对多关系
CREATE TABLE IF NOT EXISTS t_order (
    id VARCHAR(36) PRIMARY KEY COMMENT '主键id',
    customer_id VARCHAR(36) NOT NULL COMMENT '客户id, todo感觉可以作为标识id',
    status ENUM('unknown','created','pending', 'paid', 'shipped', 'completed', 'cancelled') NOT NULL COMMENT '订单状态',
    total_amount BIGINT(20) NOT NULL COMMENT '订单总金额，单位：分',
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间,精确到毫秒',
    updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间，精确到毫秒',
    INDEX idx_customer_id (customer_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 创建订单项表
-- 订单项子表，存储订单包含的商品信息，通过order_id与订单主表关联
CREATE TABLE IF NOT EXISTS t_order_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键id',
    order_id VARCHAR(36) NOT NULL COMMENT '关联订单主表的ID',
    product_id VARCHAR(36) NOT NULL COMMENT '商品id',
    quantity BIGINT NOT NULL COMMENT '商品数量',
    unit_price BIGINT(20) NOT NULL COMMENT '商品单价，单位：分',
    subtotal BIGINT(20) NOT NULL COMMENT '商品小计金额，单位：分',
    INDEX idx_order_id (order_id)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='订单商品项表';

-- 创建支付表
-- 存储订单支付信息，通过order_id与订单主表关联
CREATE TABLE IF NOT EXISTS t_payment (
    id VARCHAR(36) PRIMARY KEY COMMENT '主键id: 后续使用雪花算法生成，使用整数',
    order_id VARCHAR(36) NOT NULL COMMENT '关联订单主表的ID',
    amount BIGINT NOT NULL COMMENT '支付金额，单位：分',
    currency CHAR(3) NOT NULL DEFAULT 'CNY' COMMENT '货币类型,如CNY/USD/EUR',
    channel TINYINT UNSIGNED NOT NULL COMMENT '支付渠道(1:支付宝 2:微信 3:银行卡)',
    --status ENUM('created', 'paid', 'refunded', 'failed', 'expired', 'canceled', 'refunding', 'refund_failed', 'refunded_success') NOT NULL COMMENT '支付状态',
    status TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '支付状态(0:创建 1:已支付 2:退款中 3:退款成功 4:支付失败 5:已过期 6:退款ing 7:退款失败 8:退款成功)',
    transaction_id VARCHAR(64) COMMENT '第三方交易流水号',
    refund_transaction_id VARCHAR(64) COMMENT '第三方退款交易流水号',
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间,精确到毫秒',
    updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间，精确到毫秒',
    completed_at TIMESTAMP(3) NULL COMMENT '支付完成时间,精确到毫秒',
    INDEX idx_order_id (order_id),
    INDEX idx_status_updated (status, updated_at)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='支付表';
