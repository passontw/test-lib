package impl

import (
	"math"
	"sl.framework.com/tool"
)

type VND struct {
	BaseCurrency
}

/**
 * NewVND
 * 创建对象
 *
 * @param src  float64- 原始数值
 * @param length  int- 有效长度
 * @return *VND - 货币对象
 */

func NewVND(src float64, length int) *VND {
	if length == 0 {
		length = 0
	}
	return &VND{
		BaseCurrency{
			Name:   "VND",
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

func (c *VND) CurrencyValue() float64 {
	if c.Length == 0 {
		return math.Trunc(c.Value)
	}
	return tool.Trunc(c.Value, c.Length)
}
