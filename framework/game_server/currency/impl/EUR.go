package impl

import "sl.framework.com/tool"

type EUR struct {
	BaseCurrency
}

/**
 * NewEUR
 * 创建对象
 *
 * @param src  float64- 原始数值
 * @param length  int- 有效长度
 * @return *EUR - 货币对象
 */

func NewEUR(src float64, length int) *EUR {
	if length == 0 {
		length = 2
	}
	return &EUR{
		BaseCurrency{
			Name:   "EUR",
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

func (c *EUR) CurrencyValue() float64 {
	return tool.Trunc(c.Value, c.Length)
}
