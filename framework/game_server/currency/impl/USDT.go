package impl

import "sl.framework.com/tool"

type USDT struct {
	BaseCurrency
}

/**
 * NewUSDT
 * 创建对象
 *
 * @param src  float64- 原始数值
 * @param length  int- 有效长度
 * @return *USDT - 货币对象
 */

func NewUSDT(src float64, length int) *USDT {
	if length == 0 {
		length = 6
	}
	return &USDT{
		BaseCurrency{
			Name:   "USDT",
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

func (c *USDT) CurrencyValue() float64 {
	return tool.Trunc(c.Value, c.Length)
}
