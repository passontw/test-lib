package base

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
)

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length < unpadding {
		return []byte("unpadding error")
	}
	return origData[:(length - unpadding)]
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	return PKCS5Padding(ciphertext, blockSize)
}

func PKCS7UnPadding(origData []byte) []byte {
	return PKCS5UnPadding(origData)
}

func TripleDesEncrypt(origData string, key, iv []byte, paddingFunc func([]byte, int) []byte) (string, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	orig := paddingFunc([]byte(origData), block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(orig))
	blockMode.CryptBlocks(crypted, orig)
	return strings.ToUpper(hex.EncodeToString(crypted)), nil
}

func TripleDesDecrypt(encrypted string, key, iv []byte, unPaddingFunc func([]byte) []byte) (string, error) {
	e, err := hex.DecodeString(strings.ToLower(encrypted))
	if err != nil {
		return "", err
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(e))
	blockMode.CryptBlocks(origData, e)
	origData = unPaddingFunc(origData)
	if string(origData) == "unpadding error" {
		return "", errors.New("unpadding error")
	}
	return string(origData), nil
}

func DesEncrypt(src, key string, paddingFunc func([]byte, int) []byte) (string, error) {
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	src = string(paddingFunc([]byte(src), bs))
	if len(src)%bs != 0 {
		return "", errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, []byte(src)[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	//return hex.EncodeToString(out), nil
	sEnc := base64.StdEncoding.EncodeToString(out)
	return sEnc, nil
}

func DesDecrypt(src, key string, unPaddingFunc func([]byte) []byte) (string, error) {
	b, _ := base64.StdEncoding.DecodeString(src)
	src = string(b)
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return "", errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, []byte(src)[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = unPaddingFunc(out)
	return string(out), nil
}

func AesEncrypt(origData, key []byte, iv []byte, paddingFunc func([]byte, int) []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = paddingFunc(origData, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte, iv []byte, unPaddingFunc func([]byte) []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = unPaddingFunc(origData)
	return origData, nil
}

func AesEcbEncrypt(origData, key []byte, paddingFunc func([]byte, int) []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = paddingFunc(origData, blockSize)
	crypted := make([]byte, len(origData))
	for bs, be:=0, blockSize; bs < len(origData); bs,be = bs+blockSize, be+blockSize{
		block.Encrypt(crypted[bs:be],origData[bs:be])
	}
	return crypted, nil
}

func AesEcbDecrypt(crypted, key []byte, unPaddingFunc func([]byte) []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData := make([]byte, len(crypted))

	for bs, be:=0, blockSize; bs < len(origData); bs,be = bs+blockSize, be+blockSize{
		block.Decrypt(origData[bs:be], crypted[bs:be])
	}
	origData = unPaddingFunc(origData)
	return origData, nil
}

type hash_data_order struct {
	loginname [30]byte
	jetton    uint32
	dumb1     uint32
	gmcode    [14]byte
	playtype  uint32
	dumb2     uint32
	remark    [20]byte
}

func GetMd5(str string) string{
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Calhash(loginname string, gmcode string, playtype uint16, account uint32, remark string) string {

	var hash_data hash_data_order
	copy(hash_data.loginname[:], loginname)
	copy(hash_data.gmcode[:], gmcode)
	hash_data.playtype = uint32(playtype)
	hash_data.jetton = account
	copy(hash_data.remark[:], remark)

	hash_data.dumb1 = 0xa0fa60b8
	hash_data.dumb2 = 0x16b0a74b

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, hash_data)

	ctx := md5.New()
	ctx.Write(buf.Bytes())
	return hex.EncodeToString(ctx.Sum(nil))
}