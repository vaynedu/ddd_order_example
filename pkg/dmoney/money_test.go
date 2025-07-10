package dmoney

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestConvertStringFloat64ToCent(t *testing.T) {
	convey.Convey("Test ConvertStringFloat64ToCent function", t, func() {
		convey.Convey("When input is valid float string", func() {
			input := "10.50"
			expected := int64(1050)
			result, err := ConvertStringFloat64ToCent(input)
			convey.So(err, convey.ShouldBeNil)
			convey.So(result, convey.ShouldEqual, expected)
		})

		convey.Convey("When input is invalid float string", func() {
			input := "invalid"
			result, err := ConvertStringFloat64ToCent(input)
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(result, convey.ShouldEqual, int64(0))
		})

		convey.Convey("When input has more than two decimal places", func() {
			input := "10.123"
			expected := int64(1012) // 注意: 这里会截断而非四舍五入
			result, err := ConvertStringFloat64ToCent(input)
			convey.So(err, convey.ShouldBeNil)
			convey.So(result, convey.ShouldEqual, expected)
		})
	})
}

func TestConvertCentToStringFloat64(t *testing.T) {
	convey.Convey("Test ConvertCentToStringFloat64 function", t, func() {
		convey.Convey("When input is positive cent value", func() {
			input := int64(1050)
			expected := "10.50"
			result := ConvertCentToStringFloat64(input)
			convey.So(result, convey.ShouldEqual, expected)
		})

		convey.Convey("When input is zero", func() {
			input := int64(0)
			expected := "0.00"
			result := ConvertCentToStringFloat64(input)
			convey.So(result, convey.ShouldEqual, expected)
		})
	})
}

func TestConvertCentToFloat64(t *testing.T) {
	convey.Convey("Test ConvertCentToFloat64 function", t, func() {
		convey.Convey("When input is positive cent value", func() {
			input := int64(1050)
			expected := 10.50
			result := ConvertCentToFloat64(input)
			convey.So(result, convey.ShouldEqual, expected)
		})

		convey.Convey("When input is negative cent value", func() {
			input := int64(-500)
			expected := -5.00
			result := ConvertCentToFloat64(input)
			convey.So(result, convey.ShouldEqual, expected)
		})
	})
}

func TestConvertFloat64ToCent(t *testing.T) {
	convey.Convey("Test ConvertFloat64ToCent function", t, func() {
		convey.Convey("When input is positive float", func() {
			input := 10.50
			expected := int64(1050)
			result := ConvertFloat64ToCent(input)
			convey.So(result, convey.ShouldEqual, expected)
		})

		convey.Convey("When input has decimal that needs rounding", func() {
			input := 10.995
			expected := int64(1099) // 注意: 这里会截断而非四舍五入
			result := ConvertFloat64ToCent(input)
			convey.So(result, convey.ShouldEqual, expected)
		})
	})
}