package nacos

import (
	"os"
	"reflect"
	"runtime"
	"sl.framework.com/trace"
	"strconv"
)

/**
 * @Author: M
 * @Date: 2024/8/5 15:34
 * @Desc: 解析nacos配置
 */

// 定义一个自定义类型约束，只允许 int 和 string 和 uint64 类型
type IntOrString interface {
	~int | ~string | ~uint64
}

// IsZeroValue 判断泛型类型 T 是否为零值
func IsZeroValue[T IntOrString](value T) bool {
	return reflect.ValueOf(value).IsZero() // 是否为零值
}

// GetEnvValueOrDefault 默认值仅支持 windows系统 泛型函数，支持 int 和 string 类型
func GetEnvValueOrDefault[T IntOrString](envKey string, defaultValue T) T {
	var zeroValue T // 处理不可返回 nil 的情况
	envValue := os.Getenv(envKey)

	if envValue == "" {
		if runtime.GOOS == "windows" {
			return defaultValue
		}
		return zeroValue
	}

	switch any(defaultValue).(type) {
	case string:
		return any(envValue).(T)
	case int:
		if intValue, err := strconv.Atoi(envValue); err == nil {
			return any(intValue).(T)
		}
		trace.Warn("Error converting environment variable %s to int, using default value: %d", envKey, defaultValue)
	default:
		trace.Error("Unsupported type %T for GetEnvValueOrDefault", defaultValue)
	}

	return defaultValue
}
