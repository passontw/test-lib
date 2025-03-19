package impl

import "sl.framework.com/tool"

type PHP struct {
	BaseCurrency
}

/**
 * NewPHP
 * 创建对象
 *
 * @param src  float64- 原始数值
 * @param length  int- 有效长度
 * @return *PHP - 货币对象
 */

func NewPHP(src float64, length int) *PHP {
	if length == 0 {
		length = 4
	}
	return &PHP{
		BaseCurrency{
			Name:   "PHP",
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

func (c *PHP) CurrencyValue() float64 {
	return tool.Trunc(c.Value, c.Length)
}
