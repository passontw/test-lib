package impl

import (
	"sl.framework.com/tool"
)

type JPY struct {
	BaseCurrency
}

/**
 * NewJPY
 * 创建对象
 *
 * @param src  float64- 原始数值
 * @param length  int- 有效长度
 * @return *JPY - 货币对象
 */

func NewJPY(src float64, length int) *JPY {
	if length == 0 {
		length = 2
	}
	return &JPY{
		BaseCurrency{
			Name:   "JPY",
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

func (c *JPY) CurrencyValue() float64 {
	return tool.Trunc(c.Value, c.Length)
}
