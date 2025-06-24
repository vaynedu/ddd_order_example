package dmoney

import "strconv"

// ConvertStringFloat64ToCent 将string的float数字 * 100， 返回int64
func ConvertStringFloat64ToCent(s string) (int64, error) {
	// 1. 转换为float64
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	// 2. 乘以100
	f *= 100

	// 3. 转换为int64
	return int64(f), nil
}

// ConvertCentToStringFloat64 将int64的分数字转换为float64的元数字
func ConvertCentToStringFloat64(cents int64) string {
	return strconv.FormatFloat(float64(cents) / 100, 'f', 2, 64)
}

// ConvertCentToFloat64 将int64的分数字转换为float64的元数字
func ConvertCentToFloat64(cents int64) float64 {
	return float64(cents) / 100
}

// ConvertFloat64ToCent 将float64的元数字转换为int64的分数字
func ConvertFloat64ToCent(f float64) int64 {
	return int64(f * 100)
}
