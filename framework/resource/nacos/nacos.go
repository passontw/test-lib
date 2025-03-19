package nacos

/*
提供了与 Nacos 配置管理服务交互的实现。
这个包包含：
1. **Client 结构体**：
   - 包含 Nacos 配置客户端的结构体。
2. **初始化**：
   - `init` 函数使用环境变量或默认值初始化 Nacos 客户端，并获取和处理初始配置。
3. **函数**：
   - `GetConfFormNacos`：返回从 Nacos 获取的配置内容。
   - `GetEnvMap`：返回配置中找到的环境变量映射。
   - `NewNacosClient`：使用给定的服务器和客户端配置创建一个新的 `Client` 实例。
   - `GetConfig`：从 Nacos 获取配置内容。
   - `PublishConfig`：将新的配置内容发布到 Nacos。
   - `DeleteConfig`：从 Nacos 删除配置。
   - `setConfEnvKeys`：解析配置内容以提取环境变量键。
   - `replacePlaceholders`：用相应的环境变量值替换配置内容中的占位符。
4. **常量**：
   - 包含 Nacos 配置的常量，如数据 ID、组 ID、命名空间 ID 以及服务器地址/端口。
*/

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sl.framework.com/resource/conf/observer"
	"sl.framework.com/trace"
	"strconv"
	"strings"
)

type Client struct {
	ConfigClient config_client.IConfigClient
	Observer     observer.ConfigObserver // 引入中介接口
}

var client *Client
var envMap map[string]string // 环境变量
var content string           // 配置文件内容

const (
	dataId                     = "DATAID"                // nacos dataId
	dataIdDefault              = "roulette-resource.yml" // nacos dataId windows下默认值
	defaultGroupDefault        = "DEFAULT_GROUP"         // nacos defaultGroup windows下默认值
	namespaceId                = "NAMESPACE"             // nacos namespaceId
	namespaceIdDefault         = "g32-dev"               // nacos namespaceId windows下默认值
	nacosServerAddrPort        = "REGISTER_HOST"         // nacos nacosServerAddr
	nacosServerAddrPortDefault = "10.146.40.110:8848"    // nacos nacosServerAddr windows下默认值
)

func init() {

	// 初始化windows下环境变量
	InitEnv()

	var err error
	envMap = make(map[string]string)
	var svrAddr string
	var svrPort uint64
	_dataId := GetEnvValueOrDefault(dataId, dataIdDefault)
	_namespaceId := GetEnvValueOrDefault(namespaceId, namespaceIdDefault)
	trace.Notice("========================================================================================================================")

	trace.Info("     dataId: [%s]", _dataId)
	trace.Info("    groupId: [%s]", defaultGroupDefault)
	trace.Info("namespaceId: [%s]", _namespaceId)

	nc := GetEnvValueOrDefault(nacosServerAddrPort, nacosServerAddrPortDefault)
	_host, _ := extractHostAndPort(nc)
	host := strings.Split(_host, ":")
	if len(host) < 2 {
		log.Fatalf("unable to obtain the correct nacos configuration, please check, nacos host: %s", host)
	}
	value, err := strconv.ParseUint(host[1], 10, 64)
	if err != nil {
		trace.Warn("error converting nacos port [%s] string to uint64: %v", host[1], err)
		trace.Warn("set the default value of nacos port to 8848")
		svrPort = 8848
	} else {
		svrPort = value
	}
	svrAddr = host[0]
	trace.Notice("Nacos连接: [%s:%d]", svrAddr, svrPort)

	sc := []constant.ServerConfig{{
		IpAddr:      svrAddr,
		Port:        svrPort,
		ContextPath: "/nacos",
		Scheme:      "http",
	}}
	cc := constant.ClientConfig{
		NamespaceId:         _namespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos_log",
		CacheDir:            "tmp/nacos_cache",
		LogLevel:            "debug",
		Username:            "nacos",
		Password:            "nacos",
	}
	client, err = NewNacosClient(sc, cc)
	if err != nil {
		return
	}

	if content, err = client.GetConfig(_dataId, defaultGroupDefault); err != nil || strings.TrimSpace(content) == "" {
		trace.Error("failed to get configuration dataId: [%s] groupId: [%s], content: %s, err: %v",
			_dataId, defaultGroupDefault, content, err)
		path := filePath(localFilePath("config.yaml"), filepath.Join(".", "conf", "config.yaml"))
		byteContent, err := os.ReadFile(path)
		content = string(byteContent)
		if err != nil {
			trace.Error("配置文件加载异常: %s", err.Error())
			os.Exit(1)
		} else {
			trace.Notice("[1] 最终加载路径: %s", path)
			trace.Notice("[1] 本地配置文件: \n%s\n", content)
		}
	} else {
		trace.Notice("Nacos连接: [%s:%d] connection successful", svrAddr, svrPort)
		trace.Notice("[1] 远程配置文件: \n%s\n", content)
	}

	content = replacePlaceholders(content)
	trace.Notice("[2] 最终配置文件: \n%s\n", content)

	client.ConfigClient.ListenConfig(vo.ConfigParam{
		DataId: _dataId,
		Group:  defaultGroupDefault,
		OnChange: func(namespace, group, dataId, data string) {
			trace.Info("configuration changed: namespace=[%s], group=[%s], dataId=[%s], newData=\n%s\n", namespace, group, dataId, data)
			content = replacePlaceholders(data)
			if client.Observer != nil {
				if err := client.Observer.UpdateConfig(content); err != nil {
					trace.Error("failed to update configuration: %v", err)
				}
			}
		},
	})
	trace.Notice("========================================================================================================================")

}

func InitEnv() {
	trace.Notice("[√] Environment variable initialization starts")

	if runtime.GOOS == "windows" {
		setEnv("LOG_MODE", "dev")
		setEnv("STABLE_MODE", "dev")
		// Redis
		setEnv("REDIS_HOST", "10.146.40.80")
		setEnv("REDIS_PORT", "7001-7006")
		setEnv("REDIS_PASSWORD", "")
		setEnv("REDIS_DATABASE", "0")
		// DB
		setEnv("TIDB_HOST", "10.146.40.240:30032")
		setEnv("TIDB_USERNAME", "root")
		setEnv("TIDB_PASSWORD", "123456")
		setEnv("TIDB_DATABASE", "g32_game_dss_dev")
		// UIDDB
		setEnv("UID_HOST", "10.146.40.240:30032")
		setEnv("UID_USERNAME", "root")
		setEnv("UID_PASSWORD", "123456")
		setEnv("UID_DATABASE", "g32_uid_dev")
		// Nacos
		setEnv("REGISTER_HOST", nacosServerAddrPortDefault)
		setEnv("NAMESPACE", namespaceIdDefault)
		setEnv("DATAID", dataIdDefault)
	}
	trace.Notice("[√] Environment variable initialization completed")
}

func setEnv(k, v string) {
	if err := os.Setenv(k, v); err != nil {
		trace.Error("[x] set env error: %s", err)
	}
}

func SetObserver(ob observer.ConfigObserver) {
	client.Observer = ob
}

func GetConfFormNacos() string {
	return content
}

func GetEnvMap() map[string]string {
	return envMap
}

// NewNacosClient 创建一个新的 Client 实例
func NewNacosClient(serverConfigs []constant.ServerConfig, clientConfig constant.ClientConfig) (*Client, error) {

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("create config client failed: %v", err)
	}

	return &Client{ConfigClient: configClient}, nil
}

// GetConfig 获取配置
func (n *Client) GetConfig(dataId, group string) (string, error) {
	return n.ConfigClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
}

// PublishConfig 发布配置
func (n *Client) PublishConfig(dataId, group, content string) error {
	_, err := n.ConfigClient.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: content,
	})
	if err != nil {
		return fmt.Errorf("publish config failed: %v", err)
	}
	return nil
}

// DeleteConfig 删除配置
func (n *Client) DeleteConfig(dataId, group string) error {
	_, err := n.ConfigClient.DeleteConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return fmt.Errorf("delete config failed: %v", err)
	}
	return nil
}

// setConfEnvKeys 解析配置文件
func setConfEnvKeys(conf string) {
	re := regexp.MustCompile(`\$\{(\w+)\}`)
	matches := re.FindAllStringSubmatch(conf, -1)
	for _, match := range matches {
		key := match[1]
		if _, exists := envMap[key]; !exists {
			envMap[key] = ""
		}
	}
}

// replacePlaceholders 替换配置文件中的环境变量占位符
func replacePlaceholders(configContent string) string {
	setConfEnvKeys(configContent)

	// 正则匹配所有形如 ${abc} 的占位符
	pattern := regexp.MustCompile(`\${(\w+)}`)
	lines := strings.Split(configContent, "\n")

	for i, line := range lines {
		matches := pattern.FindAllStringSubmatch(line, -1)
		if len(matches) > 0 {
			for _, match := range matches {
				key := match[1]
				envValue := os.Getenv(key)

				if envValue != "" {
					// 替换占位符为环境变量值
					lines[i] = strings.ReplaceAll(lines[i], match[0], envValue)
					envMap[key] = envValue
					trace.Info("[+] os env [%s] value is [%s]", key, envValue)
				} else {
					// 将该行内容清空，记录错误日志
					lines[i] = ""
					trace.Error("[-] os env [%s] is undefined, line [%03d] already removed", key, i)
				}
			}
		}
	}

	content := strings.Join(lines, "\n")
	return content
}

const (
	http  = `http://`
	https = `https://`
)

// Extract host and port from the given URL string
func extractHostAndPort(rawURL string) (string, error) {
	// Ensure the URL contains a scheme, add http:// if missing
	if !strings.HasPrefix(rawURL, http) && !strings.HasPrefix(rawURL, https) {
		rawURL = http + rawURL
	}

	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// If the Host is empty, try extracting directly from the URL string
	if parsedURL.Host == "" {
		// Extract host and port part from the URL string
		hostPort := strings.SplitN(rawURL, "/", 2)[0]
		return hostPort, nil
	}

	return parsedURL.Host, nil
}

func filePath(path, defaultPath string) string {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return defaultPath
	}
	return path
}

// 本地文件初始化
func localFilePath(fileName string) string {
	var workDir string
	var err error

	if workDir, err = os.Getwd(); err != nil {
		trace.Error("localConfInit GetConfig, get path error=%v", err.Error())
		return ""
	}

	return filepath.Join(workDir, "conf", fileName)
}
