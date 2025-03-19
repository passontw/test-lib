package impl

import "sl.framework.com/tool"

type USD struct {
	BaseCurrency
}

/**
 * NewUSD
 * 创建对象
 *
 * @param src  float64- 原始数值
 * @param length  int- 有效长度
 * @return *USD - 货币对象
 */

func NewUSD(src float64, length int) *USD {
	if length == 0 {
		length = 2
	}
	return &USD{
		BaseCurrency{
			Name:   "USD",
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

func (c *USD) CurrencyValue() float64 {
	return tool.Trunc(c.Value, c.Length)
}
