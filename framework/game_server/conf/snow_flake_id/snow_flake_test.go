package snowflaker

import (
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"testing"
)

func TestSnowFlakeGetSnowFlakeId(t *testing.T) {
	_ = trace.LoggerInit()
	pDog := tool.NewWatcher("TestSnowFlakeGetSnowFlakeId")
	for i := 0; i < 100; i++ {
		GetSnowFlakeInstance().getSnowId()
	}
	pDog.Stop()
	GetSnowFlakeInstance().StopLoop()
}

// BenchmarkGetSnowFlakeId 是 GetUniqueId 的基准测试
func BenchmarkGetSnowFlakeId(b *testing.B) {
	_ = trace.LoggerInit()
	for i := 0; i < b.N; i++ {
		GetSnowFlakeInstance().getSnowId()
	}
}
