package impl

import "sl.framework.com/tool"

type HKD struct {
	BaseCurrency
}

/**
 * NewHKD
 * 创建对象
 *
 * @param src  float64- 原始数值
 * @param length  int- 有效长度
 * @return *HKD - 货币对象
 */

func NewHKD(src float64, length int) *HKD {
	if length == 0 {
		length = 2
	}
	return &HKD{
		BaseCurrency{
			Name:   "HKD",
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

func (c *HKD) CurrencyValue() float64 {
	return tool.Trunc(c.Value, c.Length)
}
