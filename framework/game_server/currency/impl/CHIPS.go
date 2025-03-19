package impl

import "sl.framework.com/tool"

type CHIPS struct {
	BaseCurrency
}

/**
 * NewCHIPS
 * 创建对象
 *
 * @param src  float64- 原始数值
 * @param length  int- 有效长度
 * @return *CHIPS - 货币对象
 */

func NewCHIPS(src float64, length int) *CHIPS {
	if length == 0 {
		length = 2
	}
	return &CHIPS{
		BaseCurrency{
			Name:   "CHIPS",
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

func (c *CHIPS) CurrencyValue() float64 {
	return tool.Trunc(c.Value, c.Length)
}
