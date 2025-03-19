package impl

import "sl.framework.com/tool"

type BRL struct {
	BaseCurrency
}

/**
 * NewBRL
 * 创建对象
 *
 * @param src  float64- 原始数值
 * @param length  int- 有效长度
 * @return *BRL - 货币对象
 */

func NewBRL(src float64, length int) *BRL {
	if length == 0 {
		length = 2
	}
	return &BRL{
		BaseCurrency{
			Name:   "BRL",
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

func (c *BRL) CurrencyValue() float64 {
	return tool.Trunc(c.Value, c.Length)
}
