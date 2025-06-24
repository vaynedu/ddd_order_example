# Order DDD Example

这是一个使用领域驱动设计(DDD)思想构建的Go语言订单系统示例。

## 项目结构

项目采用洋葱架构，主要分为以下几层：

- cmd：应用入口
- internal：应用核心代码
  - domain：领域模型层
  - application：应用服务层
  - infrastructure：基础设施层
  - interface：接口层
- pkg：公共工具包
- config：配置文件

## 技术栈

- 语言：Go 1.23
- 数据库：MySQL
- Web框架：标准库net/http
- 配置管理：viper
- ORM：gorm

## 运行步骤

1. 创建MySQL数据库：CREATE DATABASE orders;
2. 执行数据库脚本：mysql -u username -p orders < internal/infrastructure/persistence/schema.sql
3. 配置环境变量：cp .env.example .env
# 编辑.env文件，设置正确的数据库连接信息
4. 安装依赖：go mod tidy
5. 启动应用：go run cmd/server/main.go
## API接口

### 创建订单POST /api/orders
```json
{
    "customer_id": "123456",
    "items": [
        {
            "product_id": "P001",
            "quantity": 2,
            "unit_price": 9.99,
            "subtotal": 19.98
        },
        {
            "product_id": "P002",
            "quantity": 1,
            "unit_price": 19.99,
            "subtotal": 19.99
        }
    ]
}
```
### 获取订单POST /api/orders/list
```json
{
    "order_id": "9b958247-5511-4d78-ac98-a9ecee7538b3"
}

```
## 设计思想

本项目遵循DDD的核心原则：

1. 领域模型驱动设计
2. 清晰的职责分离
3. 依赖倒置原则
4. 聚合根模式
5. 领域事件模式(未来可扩展)

通过这种设计，系统具有良好的可维护性和可扩展性，业务逻辑与技术实现清晰分离。  


```
order_ddd_example/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口点
├── internal/
│   ├── domain/
│   │   └── order/
│   │       ├── entity.go        # 订单实体
│   │       ├── repository.go    # 仓储接口
│   │       └── service.go       # 领域服务
│   ├── application/
│   │   └── service/
│   │       └── order_service.go # 应用服务
│   ├── infrastructure/
│   │   ├── repository/
│   │   │   └── order_repository.go # 仓储实现
│   │   └── persistence/
│   │       └── schema.sql       # 数据库模式
│   └── interface/
│       └── handler/
│           └── order_handler.go # HTTP处理器
├── pkg/
│   └── database/
│       └── mysql.go             # MySQL连接
├── config/
│   └── config.yaml              # 配置文件
├── .env.example                 # 环境变量示例
├── go.mod                       # Go模块文件
└── README.md                    # 项目说明

```