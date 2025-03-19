package base

import (
	"container/list"
)

type Keyer interface {
	GetKey() string
}

type OrderMap struct {
	dataMap  map[string]*list.Element
	dataList *list.List
}

func NewOrderMap() *OrderMap {
	return &OrderMap{
		dataMap:  make(map[string]*list.Element),
		dataList: list.New(),
	}
}

func (p *OrderMap) Exists(data Keyer) bool {
	if _, ok := p.dataMap[data.GetKey()]; ok {
		return true
	}
	return false
}

func (p *OrderMap) Push(data Keyer) bool {
	if p.Exists(data) {
		return false
	}
	elem := p.dataList.PushBack(data)
	p.dataMap[data.GetKey()] = elem
	return true
}

func (p *OrderMap) Remove(data Keyer) {
	if !p.Exists(data) {
		return
	}
	p.dataList.Remove(p.dataMap[data.GetKey()])
	delete(p.dataMap, data.GetKey())
}

func (p *OrderMap) Size() int {
	return p.dataList.Len()
}

func (p* OrderMap) Walk(cb func(data Keyer)){
	for elem := p.dataList.Front(); elem != nil; elem = elem.Next(){
		cb(elem.Value.(Keyer))
	}
}
