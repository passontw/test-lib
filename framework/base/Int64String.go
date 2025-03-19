package base

import "strconv"

// 自定义 Int64String 类型
type Int64String int64

// 解析 JSON 时，把 string 转换为 int64
func (i *Int64String) UnmarshalJSON(data []byte) error {
	// 去掉引号
	str := string(data)
	if len(str) > 1 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// 转换成 int64
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*i = Int64String(val)
	return nil
}

// 序列化 JSON 时，把 int64 转换为 string
func (i Int64String) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatInt(int64(i), 10) + `"`), nil
}
