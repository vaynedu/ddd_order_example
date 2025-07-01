package database

import (
	"context"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func InitMySQL(ctx context.Context, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// // 自动迁移模型
	// err = db.AutoMigrate(&model.TShortURL{})
	// if err != nil {
	// 	return nil, err
	// }

	return db, nil
}
