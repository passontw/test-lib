package base

import (
	"fmt"
	"testing"
	"time"
)

type Task01 struct {
	i int
	fn func(int)
}

func (t*Task01) RunTask()  {
	t.fn(t.i)
}


func TestGoPool_Post(t *testing.T) {
	p := NewGoPool(10)
	p.Run()

	go func() {
		for i := 0; i < 100; i++{
			p.AsyncRunTask(&Task01{
				fn: func(i int) {
					fmt.Println(i)
				},
				i:i,
			})
		}
	}()

	time.Sleep(time.Second*3)
}