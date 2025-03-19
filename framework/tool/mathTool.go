package tool

import "math"

/**
 * Trunc
 * 浮点数截取指定长度的有效数字
 *
 * @param srcValue float64 - 源数值
 * @param length int - 有效数字长度
 * @return float64 - 按照要求返回的数值
 */

func Trunc(srcValue float64, length int) float64 {
	var offset float64 = 1
	if srcValue < 0 {
		offset = -1
	}
	return math.Trunc((srcValue+offset*(1/math.Pow(10, float64(length+1))))*math.Pow(10, float64(length))) / math.Pow(10, float64(length))
}
