package tool

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"sl.framework.com/trace"
	"strconv"
	"strings"
	"time"
)

func I2Int(num interface{}) (numInt int64) {
	if num == nil {
		return
	}
	var err error
	switch num.(type) {
	case string:
		numInt, err = strconv.ParseInt(num.(string), 10, 64)
		if err != nil {
			trace.Error("I2Int string to int failed err=%v", err)
		}
	case int:
		numInt = num.(int64)
	case int64:
		numInt = num.(int64)
	case float64:
		numInt = int64(num.(float64))
	}
	return
}

func I2Float(num interface{}) (numFloat float64) {
	if num == nil {
		return
	}
	var err error
	switch num.(type) {
	case string:
		numFloat, err = strconv.ParseFloat(num.(string), 64)
		if err != nil {
			trace.Error("I2Float string to int failed err=%v", err)
		}
	case float64:
		numFloat = num.(float64)
	case float32:
		numFloat = num.(float64)
	case int64:
		numFloat = float64(num.(int64))
	}
	return
}

func I2TimeUnix(timeI interface{}) int64 {
	switch timeI.(type) {
	case time.Time:
		return timeI.(time.Time).Unix()
	default:
		trace.Error("I2TimeUnix not time")
		return 0
	}
}

// 四舍五入, pos是保留小数点后几位
func Floor(x float64, pos int) float64 {
	strFormat := "%." + strconv.Itoa(pos) + "f" //%.(pos)f
	strRet := fmt.Sprintf(strFormat, x)
	ret, _ := strconv.ParseFloat(strRet, 64)
	return ret
}

// to do: test
func ObjectMapToSlice(dataMap map[string]map[string]interface{}) []map[string]interface{} {
	dataList := make([]map[string]interface{}, 0)
	for _, data := range dataMap {
		dataList = append(dataList, data)
	}
	return dataList
}

// 截取字符串
func SubString(str string, start int, end int) (result string) {
	return string([]rune(str)[start:end])
}

// 转换字符串类型
func ToString(param interface{}) string {
	return fmt.Sprintf("%v", param)
}

// interfaces to int64T
func ToInt64(num interface{}) (numInt int64, err error) {
	if num == nil {
		return
	}
	switch num.(type) {
	case string:
		numInt, err = strconv.ParseInt(num.(string), 10, 64)
		if err != nil {
			trace.Error("ToInt string to int failed err=%v", err)
			return 0, err
		}
	case int:
		numInt = int64(num.(int))
	case int8:
		numInt = int64(num.(int8))
	case int16:
		numInt = int64(num.(int16))
	case int32:
		numInt = int64(num.(int32))
	case int64:
		numInt = num.(int64)
	case uint:
		numInt = int64(num.(uint))
	case uint8:
		numInt = int64(num.(uint8))
	case uint16:
		numInt = int64(num.(uint16))
	case uint32:
		numInt = int64(num.(uint32))
	case uint64:
		numInt = int64(num.(uint64))
	case float32, float64:
		numInt = int64(num.(float64))
	}
	return numInt, err
}

// interfaces to int32
func ToInt32(num interface{}) (numInt int32, err error) {
	temp, err := ToInt64(num)
	return int32(temp), err
}

// interfaces to int
func ToInt(num interface{}) (numInt int, err error) {
	temp, err := ToInt64(num)
	return int(temp), err
}

// interfaces to float64
func ToFloat(num interface{}) (numFloat float64, err error) {
	if num == nil {
		return
	}
	switch num.(type) {
	case string:
		numFloat, err = strconv.ParseFloat(num.(string), 64)
		if err != nil {
			trace.Error("ToFloat string to int failed err=%v", err)
			return 0, err
		}
	case float32:
		numFloat = float64(num.(float32))
	case float64:
		numFloat = num.(float64)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if temp, err := ToInt64(num); err != nil {
			return 0, err
		} else {
			numFloat = float64(temp)
		}
	}
	return numFloat, err
}

// 判断字符串是否以给定字符串开始
func StartsWith(str string, start string) bool {
	tempStr := []rune(str)
	tempStart := []rune(start)
	return string(tempStr[:len(tempStart)]) == start
}

// 判断字符串是否以给定字符串结尾
func EndsWith(str string, end string) bool {
	tempStr := []rune(str)
	tempEnd := []rune(end)
	s := string(tempStr[len(tempStr)-len(tempEnd):])
	return s == end
}

// ShuffleSliceAny 打乱Slice的顺序
func ShuffleSliceAny[T any](data []T) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})
}

func GetLocalIp() string {
	ipLocal := "127.0.0.1"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ipLocal
	}

	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipLocal = ipNet.IP.String()
			}
		}
	}

	return ipLocal
}

// ExtractPathFromIndex 按系统分隔符提取路径中从指定索引开始的内容
func ExtractPathFromIndex(path string, partIndex int) string {

	// 修正路径
	path = filepath.FromSlash(path)

	separator := string(os.PathSeparator) // 获取系统分隔符

	// 将路径根据分隔符分割为片段
	parts := strings.Split(path, separator)

	// 处理倒数索引
	if partIndex < 0 {
		partIndex = len(parts) + partIndex
	}

	// 检查索引是否越界
	if partIndex < 0 || partIndex >= len(parts) {
		// 返回整个路径
		return path
	}

	// 拼接从索引开始的部分
	return strings.Join(parts[partIndex:], separator)
}

// 判断接口是否为空
func IsEmpty(inte interface{}) bool {
	return inte == nil || reflect.ValueOf(inte).IsZero()
}

// 获取当前时间，毫秒级

func Current() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

/**
 * SplitList
 * 把list切割为size大小的列表
 *
 * @param srcList  []T - 原始list
 * @param size  int - 子list的大小
 * @return [][]T -
 */

func SplitList[T any](srcList []T, size int) [][]T {
	result := make([][]T, 0)
	if len(srcList) == 0 {
		return result
	}
	// 遍历切片，并按照给定的大小切割
	for size < len(srcList) {
		// 切割前 size 个元素，添加到结果中
		srcList, result = srcList[size:], append(result, srcList[0:size:size])
	}
	// 处理剩余部分
	return append(result, srcList)
}

/**
 * GenerateRandomString
 * 随机生成给定长度的字符串
 *
 * @param length  int32 - 原始list
 * @return string - 生成的随机字符串
 */

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result []rune
	// 创建随机数生成器，使用当前 Unix 时间作为种子
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := 0; i < length; i++ {
		result = append(result, rune(charset[r.Intn(len(charset))]))
	}

	return string(result)
}

/**
 * GenerateRandomRange
 * 随机生成min~max的随机数
 *
 * @param min int - 最小值
 * @param max int - 最大值
 * @return int64 - 返回值说明
 */

func GenerateRandomRange(min, max int) int {
	if min >= max {
		panic("min must be less than max")
	}
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	return r.Intn(max-min) + min
}
func Int64SliceToStringSlice(intList []int64) []string {
	dstList := make([]string, len(intList))
	for i, val := range intList {
		dstList[i] = strconv.FormatInt(val, 10) // 将 int64 转换为字符串
	}
	return dstList
}

/**
 * FormatTime
 * 把时间戳（毫秒级）格式化为2006-01-02 15:04:05.000的格式
 *
 * @param timestamp int64 - 时间戳（毫秒级）
 * @return string - 格式化时间
 */

func FormatTime(timestamp int64) string {
	ti := time.UnixMilli(timestamp)
	str := ti.Format("2006-01-02 15:04:05.000")
	return str
}
