package database

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
)

// MySQLConfig MySQL配置
type MySQLConfig struct {
    Username string
    Password string
    Host     string
    Port     int
    DBName   string
    Timeout  int
}

// NewMySQLDB 创建MySQL连接
func NewMySQLDB(ctx context.Context, dsn string) (*sqlx.DB, error) {
    db, err := sqlx.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("打开数据库连接失败: %w", err)
    }

    // 设置连接池参数
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(30 * time.Minute)

    // 测试连接
    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("ping数据库失败: %w", err)
    }

    return db, nil
}  