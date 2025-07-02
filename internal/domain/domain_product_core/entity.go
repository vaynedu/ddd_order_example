package domain_product_core

// Product 商品领域模型
type Product struct {
	ID     string
	Name   string
	Status ProductStatus
	Price  int64
	// 其他领域属性...
}

// ProductStatus 商品状态枚举
type ProductStatus int

const (
	StatusValid      ProductStatus = iota // 有效
	StatusInvalid                         // 无效
	StatusDeleted                         // 已删除
	StatusOutOfStock                      // 售罄
)

func GetProductStatusDetail(status ProductStatus) string {
	switch status {
	case StatusValid:
		return "有效"
	case StatusInvalid:
		return "无效"
	case StatusDeleted:
		return "已删除"
	case StatusOutOfStock:
		return "售罄"
	default:
		return "未知"
	}
}
