package impl

import (
	"sl.framework.com/game_server/currency/intefaces"
	"sl.framework.com/trace"
)

type BaseCurrency struct {
	Value  float64
	Name   string
	Length int
}

/**
 * FunctionName
 * 函数功能描述
 *
 * @param PARAM - 参数说明
 * @return RETURN - 返回值说明
 */

func NewCurrency(currency string, src float64) intefaces.ICurrency {
	if len(currency) == 0 {
		trace.Error("创建货币类型  参数非法 NewCurrency: currency is empty")
		return nil
	}
	switch currency {
	case "CNY":
		return NewCNY(src, 0)
	case "USD":
		return NewUSD(src, 0)
	case "EUR":
		return NewEUR(src, 0)
	case "GBP":
		return NewGBP(src, 0)
	case "JPY":
		return NewJPY(src, 0)
	case "BRL":
		return NewBRL(src, 0)
	case "CHIPS":
		return NewCHIPS(src, 0)
	case "HKD":
		return NewHKD(src, 0)
	case "INR":
		return NewINR(src, 0)
	case "PHP":
		return NewPHP(src, 0)
	case "USDT":
		return NewUSDT(src, 0)
	case "VND":
		return NewVND(src, 0)
	default:
		trace.Error("创建货币类型，未知的货币类型 currency:%v", currency)
		return nil

	}
}
