package conf

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"os"
	"path/filepath"
	"regexp"
	"sl.framework.com/trace"
	"strconv"
	"strings"
	"sync"
)

var (
	nacosInitOnce sync.Once
	configClient  config_client.IConfigClient
)

const (
	nacosRegisterHost = "REGISTER_HOST" //nacos域名
	nacosNameSpace    = "NAMESPACE"     //nacos命名空间

	nacosUserName = "USER_NAME" //登陆nacos账号
	nacosPassword = "PASSWORD"  //登陆nacos密码

	defaultNacosUserName = "nacos" //默认登陆nacos账号
	defaultNacosPassword = "nacos" //默认登陆nacos密码
)

/*
	ConfigLocation 配置文件位置
	OnLocalConf:配置信息在本地 位于conf目录下
	OnServerConf:配置信息在Nacos服务器
*/

type ConfigLocation string

const (
	OnLocalConf  ConfigLocation = "LocalConf"
	OnServerConf ConfigLocation = "ServerConf"
)

/**
 * NacosClientInitOnce
 * 初始化 Nacos
 *
 * @param location ConfigLocation - yml位置
 * 			框架实际使用时都从nacos服务器读取 调试的时候才会用使用本地yml文件
 * 			所以该参数没有从上层具体游戏中传入
 * @return isSuccess bool - 初始化成功或者失败
 */

func NacosClientInitOnce(location ConfigLocation) (isSuccess bool) {
	nacosInitOnce.Do(func() {
		if configClient != nil {
			trace.Error("NacosClientInitOnce configClient init already")
			return
		}
		switch location {
		case OnLocalConf:
			isSuccess = localConfInit()
		case OnServerConf:
			isSuccess = nacosClientInit()
		}
	})

	return
}

// getDetailInHost 解析包换域名和端口字符串中的域名和端口
func getDetailInHost(host string) (string, uint64) {
	//分割成两部分
	re := regexp.MustCompile(`^(.*):(.*)$`)
	matches := re.FindStringSubmatch(host)

	strIp, strPort := "", "0"
	iPort := 0
	if len(matches) >= 3 {
		strIp = strings.TrimPrefix(matches[1], "http://")
		strPort = matches[2]
		iPort, _ = strconv.Atoi(strPort)
	}
	trace.Info("getDetailInHost host=%v, ip=%v, port=%v", host, strIp, iPort)

	return strIp, uint64(iPort)
}

// 本地文件初始化
func localConfInit() (isSuccess bool) {
	var workDir string
	var content []byte
	var err error

	if workDir, err = os.Getwd(); err != nil {
		trace.Error("localConfInit GetConfig, get path error=%v", err.Error())
		return
	}

	fileName := "server.yaml"
	configPath := filepath.Join(workDir, "framework", "game_server", "conf", fileName)
	if content, err = os.ReadFile(configPath); err != nil {
		trace.Error("localConfInit configPath=%v, read file error=%v", configPath, err.Error())
		return
	}

	trace.Info("localConfInit GetConfig content=\n%v", string(content))
	if err = parseContent(string(content)); nil != err {
		trace.Error("localConfInit parseContent, parse error=%v", err.Error())
		return
	}
	isSuccess = true

	return
}

// nacosClientInit 初始化Nacos
func nacosClientInit() (isSuccess bool) {
	var strHost, strNamespace, strDataId string
	var isExist bool

	//REGISTER_HOST为Pod中环境变量
	if isExist, strHost = getEnv(nacosRegisterHost); !isExist {
		trace.Error("nacosClientInit REGISTER_HOST env not exist")
		return
	}
	//NAMESPACE为Pod中环境变量
	if isExist, strNamespace = getEnv(nacosNameSpace); !isExist {
		trace.Error("nacosClientInit NAMESPACE env not exist")
		return
	}
	userName := getEnvWithDefault(nacosUserName, "")
	password := getEnvWithDefault(nacosPassword, "")

	strDataId = GetConfigurationFileName()
	strIp, uiPort := getDetailInHost(strHost)
	trace.Info("nacosClientInit nacos host=%v, ip=%v, port=%v, namespace=%v, dataId=%v",
		strHost, strIp, uiPort, strNamespace, strDataId)

	var err error
	sc := []constant.ServerConfig{{
		IpAddr: strIp,
		Port:   uiPort,
	}}
	// 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
	cc := constant.ClientConfig{
		NamespaceId:         strNamespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/trace",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
		Username:            defaultNacosUserName,
		Password:            defaultNacosPassword,
	}
	if len(userName) != 0 {
		cc.Username = userName
	}
	if len(password) != 0 {
		cc.Password = password
	}
	configClient, err = clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		trace.Error("NacosClientInit error=%v", err.Error())
		return
	}
	trace.Info("NacosClientInit success cc=%+v, sc=%+v", cc, sc)

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: strDataId,
		Group:  "DEFAULT_GROUP",
	})
	if err != nil {
		trace.Info("NacosClientInit GetConfig error=%v", err.Error())
		return
	}
	trace.Info("NacosClientInit GetConfig content=\n%v", content)
	if err = parseContent(content); nil != err {
		trace.Error("NacosClientInit GetConfig, parse error=%v", err.Error())
	}

	err = configClient.ListenConfig(vo.ConfigParam{
		DataId:   strDataId,
		Group:    "DEFAULT_GROUP",
		OnChange: onChangeCallback,
	})

	isSuccess = true
	return
}

/**
 * onChangeCallback
 * 配置变化时的回调
 *
 * @param namespace string - 名字空间
 * @param group string - 组名
 * @param dataId string - 数据Id
 * @param data string - 接收到数据 解析到配置结构中
 * @return
 */

func onChangeCallback(namespace, group, dataId, data string) {
	preLogLevel := ServerConf.Common.LogLevel

	trace.Info("NacosClientInit nacos configuration changed, namespace=%v, group=%v, dataId=%v, data=\n%v",
		namespace, group, dataId, data)
	if err := parseContent(data); nil != err {
		trace.Error("NacosClientInit nacos configuration changed, parse error=%v", err.Error())
		return
	}

	//读取到nacos配置以后 如果日志级别变化则重新设置日志打印级别
	if preLogLevel != ServerConf.Common.LogLevel {
		trace.SetLevel(GetLogLevel())
	}
}
