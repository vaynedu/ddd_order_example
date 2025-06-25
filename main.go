package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"github.com/vaynedu/ddd_order_example/internal/infrastructure/di"
	"github.com/vaynedu/ddd_order_example/pkg/database"
)

func initConfig() error {
	viper.SetConfigName("config_prod")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}

func main() {
	// 解析命令行参数
	flag.String("config", "config/config_prod.yaml", "配置文件路径")
	flag.Parse()

	// 初始化配置
	if err := initConfig(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	ctx := context.Background()

	// 从配置文件读取数据库连接信息
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.name"),
	)

	// 初始化数据库连接
	db, err := database.InitMySQL(ctx, dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 通过Wire依赖注入初始化处理器
	orderHandler, err := di.InitializeOrderHandler(db)
	if err != nil {
		log.Fatalf("依赖注入初始化失败: %v", err)
	}

	// 注册路由
	mux := http.NewServeMux()
	mux.HandleFunc("/api/orders/create", orderHandler.CreateOrder)
	mux.HandleFunc("/api/orders/list", orderHandler.GetOrder)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    viper.GetString("server.address"),
		Handler: mux,
	}

	// 启动服务器
	go func() {
		log.Printf("服务器启动，监听地址: %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在优雅关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭失败: %v", err)
	}

	log.Println("服务器已关闭")
}
