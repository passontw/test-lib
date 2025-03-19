package base

import (
	"fmt"
	"testing"
)

type video struct {
	vid string
	name string
}

func (p *video) GetKey() string {
	return p.vid
}

func TestNewOrderMap(t *testing.T) {
	dataMap := NewOrderMap()
	for i := 0; i < 10; i++ {
		vid := fmt.Sprintf("%03d", i)
		dataMap.Push(&video{vid: vid, name: vid})
	}
	dataMap.Walk(func(data Keyer) {
		fmt.Println(data.(*video).name)
	})
}
