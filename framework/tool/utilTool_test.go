package tool

import "testing"

func TestSubString(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{"123abc", "3a"},
		{"1c_？", "_？"},
		{"中国人ab123", "人a"},
	}

	for _, v := range tests {
		num := SubString(v.in, 2, 4)
		if num != v.expected {
			t.Logf("in:%v,expected:%v,result:%v", v.in, v.expected, num)
			t.Fail()
		}
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected int64
	}{
		{"123", 123},
		{123, 123},
		{123.4, 123},
		{123.6, 123},
		{-123, -123},
		{int8(123), 123},
		{int16(123), 123},
		{int32(123), 123},
		{int64(123), 123},
		{uint8(123), 123},
		{uint16(123), 123},
		{uint32(123), 123},
		{uint64(123), 123},
	}

	for _, v := range tests {
		num, err := ToInt64(v.in)
		if err != nil || num != v.expected {
			t.Logf("in:%v,expected:%v,result:%v,err:%v", v.in, v.expected, num, err)
			t.Fail()
		}
	}
}

func TestToInt32(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected int32
	}{
		{"123", 123},
		{123, 123},
		{123.4, 123},
		{123.6, 123},
		{-123, -123},
		{int8(123), 123},
		{int16(123), 123},
		{int32(123), 123},
		{int64(123), 123},
		{uint8(123), 123},
		{uint16(123), 123},
		{uint32(123), 123},
		{uint64(123), 123},
	}

	for _, v := range tests {
		num, err := ToInt32(v.in)
		if err != nil || num != v.expected {
			t.Logf("in:%v,expected:%v,result:%v,err:%v", v.in, v.expected, num, err)
			t.Fail()
		}
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected int
	}{
		{"123", 123},
		{123, 123},
		{123.4, 123},
		{123.6, 123},
		{-123, -123},
		{int8(123), 123},
		{int16(123), 123},
		{int32(123), 123},
		{int64(123), 123},
		{uint8(123), 123},
		{uint16(123), 123},
		{uint32(123), 123},
		{uint64(123), 123},
	}

	for _, v := range tests {
		num, err := ToInt(v.in)
		if err != nil || num != v.expected {
			t.Logf("in:%v,expected:%v,result:%v,err:%v", v.in, v.expected, num, err)
			t.Fail()
		}
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected float64
	}{
		{"123", 123},
		{123, 123},
		{123.4, 123.4},
		{123.6, 123.6},
		{-123, -123},
		{int8(123), 123},
		{int16(123), 123},
		{int32(123), 123},
		{int64(123), 123},
		{uint8(123), 123},
		{uint16(123), 123},
		{uint32(123), 123},
		{uint64(123), 123},
	}

	for _, v := range tests {
		num, err := ToFloat(v.in)
		if err != nil || num != v.expected {
			t.Logf("in:%v,expected:%v,result:%v,err:%v", v.in, v.expected, num, err)
			t.Fail()
		}
	}
}

func TestStartsWith(t *testing.T) {
	if !StartsWith("a123", "a1") {
		t.Logf("in:%v,expected:%v,result:%v", "a123", true, false)
		t.Fail()
	}
	if !StartsWith("？a123", "？a1") {
		t.Logf("in:%v,expected:%v,result:%v", "？a123", true, false)
		t.Fail()
	}
	if !StartsWith("？a123", "？") {
		t.Logf("in:%v,expected:%v,result:%v", "？a123", true, false)
		t.Fail()
	}
	if !StartsWith("中国a123", "中国a1") {
		t.Logf("in:%v,expected:%v,result:%v", "中国a123", true, false)
		t.Fail()
	}
	if !StartsWith("中国a123", "中") {
		t.Logf("in:%v,expected:%v,result:%v", "中国a123", true, false)
		t.Fail()
	}
}

func TestEndsWith(t *testing.T) {
	if !EndsWith("a123", "23") {
		t.Logf("in:%v,expected:%v,result:%v", "a123", true, false)
		t.Fail()
	}
	if !EndsWith("？a123", "a123") {
		t.Logf("in:%v,expected:%v,result:%v", "？a123", true, false)
		t.Fail()
	}
	if !EndsWith("？a123", "？a123") {
		t.Logf("in:%v,expected:%v,result:%v", "？a123", true, false)
		t.Fail()
	}
	if !EndsWith("中国a123", "国a123") {
		t.Logf("in:%v,expected:%v,result:%v", "中国a123", true, false)
		t.Fail()
	}
}
