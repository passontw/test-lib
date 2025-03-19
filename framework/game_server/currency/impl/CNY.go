package impl

import "sl.framework.com/tool"

type CNY struct {
	BaseCurrency
}

/**
 * NewCNY
 * 创建对象
 *
 * @param src  float64- 原始数值
 * @param length  int- 有效长度
 * @return *CNY - 货币对象
 */

func NewCNY(src float64, length int) *CNY {
	if length == 0 {
		length = 2
	}
	return &CNY{
		BaseCurrency{
			Name:   "CNY",
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

func (c *CNY) CurrencyValue() float64 {
	return tool.Trunc(c.Value, c.Length)
}
