package tool

import (
	"testing"
	"time"
)

func TestTrunc(t *testing.T) {
	re := Trunc(329, 4)
	t.Logf("结果：%v", re)

	currentTimeStamp := time.Now().UnixMilli()

	ti := time.UnixMilli(currentTimeStamp)
	str := ti.Format("2006-01-02 15:04:05.000")
	t.Logf("当前时间为：%v", str)
}
