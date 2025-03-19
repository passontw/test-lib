package impl

import "sl.framework.com/tool"

type INR struct {
	BaseCurrency
}

/**
 * NewINR
 * 创建对象
 *
 * @param src  float64- 原始数值
 * @param length  int- 有效长度
 * @return *INR - 货币对象
 */

func NewINR(src float64, length int) *INR {
	if length == 0 {
		length = 2
	}
	return &INR{
		BaseCurrency{
			Name:   "INR",
			Value:  src,
			Length: length,
		},
	}
}

/**
 * CurrencyValue
 * 获取货币有效值
 *
 * @return float64  -
 */

func (c *INR) CurrencyValue() float64 {
	return tool.Trunc(c.Value, c.Length)
}
