package impl

import "sl.framework.com/tool"

type GBP struct {
	BaseCurrency
}

/**
 * NewGBP
 * 创建对象
 *
 * @param src  float64- 原始数值
 * @param length  int- 有效长度
 * @return *GBP - 货币对象
 */

func NewGBP(src float64, length int) *GBP {
	if length == 0 {
		length = 2
	}
	return &GBP{
		BaseCurrency{
			Name:   "GBP",
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

func (c *GBP) CurrencyValue() float64 {
	return tool.Trunc(c.Value, c.Length)
}
