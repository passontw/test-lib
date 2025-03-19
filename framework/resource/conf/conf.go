package conf

/**
 * @Author: M
 * @Date: 2024/7/30 10:05
 * @Desc: 配置文件读取和管理
 *
 * 本文件定义了 `cfg` 结构体，用于处理配置文件的读取和更新。
 * 包含以下功能：
 * 1. 从 Nacos 获取配置并解析。
 * 2. 提供获取配置的接口。
 * 3. 实现 `observer.ConfigObserver` 接口，以便于接收配置更新通知。
 *
 * 文件中主要包含以下函数和方法：
 * - `init()`: 初始化配置，从 Nacos 加载配置并设置观察者。
 * - `Get()`: 返回当前配置对象。
 * - `Section(env, key string) string`: 获取指定环境和键的配置值。
 * - `UpdateConfig(data string) error`: 实现 `observer.ConfigObserver` 接口的更新配置方法。
 * - `reload(data string) error`: 重新加载配置数据。
 */

import (
	"fmt"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/config/yaml"
	"log"
	"os"
	"sl.framework.com/resource/nacos"
	"sl.framework.com/trace"
	"strings"
)

type cfg struct {
	conf        config.Configer
	resetLogger func(int)
}

var c *cfg

// 初始化配置，从 Nacos 加载配置并设置观察者
func init() {
	c = &cfg{}
	if err := loadConfig(); err != nil {
		log.Fatal(err)
	}
	nacos.SetObserver(c)
}

// 加载配置文件并解析
func loadConfig() error {
	envVars := os.Environ()
	trace.Info("[1    ] load os env")
	for i, envVar := range envVars {
		kv := strings.Split(envVar, "=")
		if len(kv) >= 1 && kv[0] != "" {
			trace.Info("[1-%03d] [+] os env [%s] value is [%s]", i+1, kv[0], kv[1])
		}
	}
	trace.Info("[2    ] load os env(%d) finished.", len(envVars))

	// 从 Nacos 获取配置并解析
	yml := yaml.Config{}
	confData := nacos.GetConfFormNacos()
	if confData == "" {
		return fmt.Errorf("no configuration found from Nacos")
	}

	var err error
	if c.conf, err = yml.ParseData([]byte(confData)); err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	section, _ := c.conf.GetSection("beego")
	trace.Info("beego conf map: %v", section)

	return nil
}

func SetResetLogger(fn func(int)) {
	c.resetLogger = fn
}

// 获取当前配置对象
func Get() config.Configer {
	return c.conf
}

// Section retrieves the value of a specified key from a configuration section.
// The configuration section is identified by the 'env' parameter, and the key
// within that section is specified by the 'key' parameter. Both 'env' and 'key'
// are converted to lowercase to ensure case-insensitive lookup.
//
// Parameters:
// - env: A string representing the configuration section identifier (e.g., "development", "production").
// - key: A string representing the key within the specified section whose value needs to be retrieved.
//
// Returns:
//   - A string value associated with the specified key within the given section. If the section or key does not exist,
//     an empty string is returned.
//
// Example usage:
// value := Section("Development", "DatabaseURL")
// This will return the value associated with the "databaseurl" key in the "development" section of the configuration.
func Section(env, key string) string {
	section, _ := c.conf.GetSection(env)
	value := section[key]
	if strings.TrimSpace(value) == "" {
		trace.Error("config: [%s], key: [%s], value is null", env, key)
	} else {
		trace.Info("config: [%s], key: [%s], value: [%s]", env, key, value)
	}
	return value
}

func SectionDefault(env, key, defaultVal string) string {
	section, _ := c.conf.GetSection(env)
	value := section[key]
	if strings.TrimSpace(value) == "" {
		value = defaultVal
		trace.Error("config: [%s], key: [%s], value is null, set to default: [%s]", env, key, defaultVal)
	} else {
		trace.Info("config: [%s], key: [%s], value: [%s]", env, key, value)
	}
	return value
}

func SectionArray(env string) []map[string]interface{} {
	array, _ := c.conf.DIY(env)
	if arr, ok := array.([]interface{}); ok {
		data := make([]map[string]interface{}, 0, len(arr))
		for _, item := range arr {
			// 安全类型断言，确保每个元素都是 map[string]interfaces{}
			if m, ok := item.(map[string]interface{}); ok {
				data = append(data, m)
			}
		}
		return data
	}
	return nil
}

// 实现 observer.ConfigObserver 接口，更新配置
func (c *cfg) UpdateConfig(data string) error {
	return reload(data)
}

// 重新加载配置数据
func reload(data string) error {
	var yml yaml.Config
	rc, err := yml.ParseData([]byte(data))
	if err != nil {
		return fmt.Errorf("failed to reload configuration: %v", err)
	}
	c.conf = rc
	// 重新设置日志级别
	if "prod" == Section("beego", "logMode") {
		c.resetLogger(trace.LevelInfo)
	} else {
		c.resetLogger(trace.LevelDebug)
	}
	trace.Info("Configuration reloaded successfully")
	return nil
}
