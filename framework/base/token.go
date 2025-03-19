package base


import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type UID_TOKEN struct {
	Val [2]uint64
}

func NewToken(arr []byte) *UID_TOKEN{
	token := &UID_TOKEN{}
	binary.Read(bytes.NewBuffer(arr[:8]), binary.LittleEndian, &token.Val[0])
	binary.Read(bytes.NewBuffer(arr[8:]), binary.LittleEndian, &token.Val[1])
	return token
}

func (p* UID_TOKEN) ToBytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, p.Val)
	return buf.Bytes()
}

func C2H(c byte) byte {
	arr := []byte("0123456789abcdef")
	return arr[c]
}

func H2C(c byte) byte {
	if c >= '0' && c <= '9' {
		return c - '0'
	}

	if c >= 'a' && c <= 'f' {
		return c - 'a' + 10
	}

	if c >= 'A' && c <= 'F' {
		return c - 'A' + 10
	}
	return '0'
}

func Btoxa(arr []byte) string {
	nSrcLen := len(arr)
	nOutLen := nSrcLen<<1 + 1
	i := 0
	hexbuf := make([]byte, nOutLen)
	m := 0
	for nSrcLen > 0 {
		nSrcLen -= 1
		c := arr[i]
		i += 1
		hexbuf[m] = C2H((c >> 4) & 0xF)
		m += 1
		hexbuf[m] = C2H(c & 0xF)
		m += 1
	}
	return string(hexbuf)
}

func Xatob(hexbuf string, x int) ([]byte, error) {
	destbuf := make([]byte, x)
	nOutLen := x
	nHexLen := len(hexbuf) - 1
	if nHexLen%2 != 0 {
		fmt.Println(nHexLen)
		return nil, errors.New("nHexLen % 2 != 0")
	}

	if (nHexLen >> 1) > nOutLen {
		return nil, errors.New("(nHexLen >> 1) > nOutLen")
	}

	d := 0
	h := 0
	for nHexLen > 0 {
		destbuf[d] = (H2C(hexbuf[h]) << 4) & 0xF0
		h++
		destbuf[d] |= H2C(hexbuf[h]) & 0x0F
		d++
		h++
		nHexLen -= 2
	}
	return destbuf, nil
}
