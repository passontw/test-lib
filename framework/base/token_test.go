package base

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"
)

func Test123(t *testing.T) {
	//str := base.GetMd5("71001gamesvr71001")
	//fmt.Println(str)
	h := md5.New()
	str := "71001gamesvr71001"
	h.Write([]byte(str))
	arr := h.Sum(nil)//
	fmt.Println(arr)

	token := NewToken(arr)
	fmt.Println(token.Val)
	s := Btoxa(arr)
	fmt.Println(hex.EncodeToString(arr))
	fmt.Println(s)

	token1 := NewToken(arr)
	arr1 := token1.ToBytes()
	fmt.Println(arr1)

	fmt.Printf("\nBtoxa====>\n")
	s1 := Btoxa(arr1)
	fmt.Println(s1)


	fmt.Printf("\nXatob====>\n")
	aa, err := Xatob(s1, 16)
	fmt.Println(aa, err)


	//a := 10
	//fmt.Println(a<<1 + 1)
}