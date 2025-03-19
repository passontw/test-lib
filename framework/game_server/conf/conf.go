package conf

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"regexp"
	"sl.framework.com/trace"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	configMutex sync.Mutex
	ServerConf  *Configuration

	/* 该文件是nacos中配置文件名字 对应nacos中的DataId字段 有框架使用者传入*/
	configurationFileName string
)

// RedisInfo redis集群信息
type (
	RedisInfo struct {
		Host      string `yaml:"host"`
		Port      string `yaml:"port"`
		Db        string `yaml:"defaultDb"`
		UserName  string `yaml:"userName"`
		Password  string `yaml:"password"`
		KeyPrefix string `yaml:"keyPrefix"`
		Mode      string `yaml:"mode"` // redis工作模式 single:单节点模式 cluster:集群模式
	}

	// DbItem mysql数据库信息
	DbItem struct {
		Host      string `yaml:"host"`
		Database  string `yaml:"database"`
		AliasName string `yaml:"aliasName"`
		UserName  string `yaml:"username"`
		Password  string `yaml:"password"`
	}

	// Database mysql数据库信息
	Database struct {
		UidDb  DbItem `yaml:"uid"`
		GameDb DbItem `yaml:"gameDb"`
	}

	// Platform 能力平台信息
	Platform struct {
		Host          string `yaml:"host"`
		Port          int    `yaml:"port"`
		RetryTime     int    `yaml:"retryTime"`
		RetryInterval int    `yaml:"retryInterval"`
	}

	// Common 部分配置
	Common struct {
		HeartbeatInterval int    `yaml:"heartbeatInterval"` //心跳间隔(s)
		HeartbeatExpired  int    `yaml:"heartbeatExpired"`  //心跳redis超时时间(s)
		LoopTaskExpired   int    `yaml:"loopTaskExpired"`   //任务超时时间
		UserLimitSwitch   string `yaml:"userLimitSwitch"`   //个人限红开关 on:开启个人限红校验 off:关闭个人限红校验
		LogLevel          string `yaml:"logLevel"`          //日志登记
		DrawSize          int    `yaml:"drawSize"`          //开奖分片中，每一片的大小
		BetConfirmSIze    int    `yaml:"betConfirmSize"`    //提交注单分片中，每一片的大小
		GameId            int    `yaml:"gameId"`            //游戏Id
	}

	// Rocket 相关配置
	Rocket struct {
		NameServer            string            `yaml:"nameServer"`
		ProducerQueueMaxLen   int               `yaml:"producerQueueMaxLen"`
		JoinMessageRoomTopic  string            `yaml:"joinMessageRoomTopic"`
		LeaveMessageRoomTopic string            `yaml:"leaveMessageRoomTopic"`
		Retries               int               `yaml:"retries"`
		GameDrawTopicsIn      []TopicConfigItem `yaml:"gameDrawTopicsIn"`
		GameDrawTopicsOut     []TopicConfigItem `yaml:"gameDrawTopicsOut"`
		BetConfirmTopicsIn    []TopicConfigItem `yaml:"betConfirmTopicsIn"`
		BetConfirmTopicsOut   []TopicConfigItem `yaml:"betConfirmTopicsOut"`
	}

	TopicConfigItem struct {
		TopicName  string `yaml:"topicName"`
		TopicGroup string `yaml:"topicGroup"`
	}
	BeeGoConfig struct {
		GraceEnable  bool   `yaml:"graceEnable"`
		GracePort    int    `yaml:"gracePort"`
		GraceTimeOUt int    `yaml:"graceTimeOut"`
		AdminEnable  bool   `yaml:"adminEnable"`
		RunMode      string `yaml:"runMode"`
	}
	// Http 相关配置
	Http struct {
		HttpConnectTimeout   int `yaml:"httpConnectTimeout"`   //http连接超时时间 单位秒
		HttpReadWriteTimeout int `yaml:"httpReadWriteTimeout"` //http读写超时时间单位秒
	}

	// Configuration 服务配置信息
	Configuration struct {
		RedisInfo      RedisInfo   `yaml:"redis"`
		Database       Database    `yaml:"database"`
		Platform       Platform    `yaml:"platform"`
		Common         Common      `yaml:"common"`
		Rocketmq       Rocket      `yaml:"rocketmq"`
		BeegoCFG       BeeGoConfig `yaml:"beego"`
		Http           Http        `yaml:"http"`
		ServerId       int64       `yaml:"serverId"`
		GameConfig     GameConfig  `yaml:"gameConfig"`
		ConfigFileName string      //配置文件名字 有具体游戏传入并设置
		AgentToken     string      //跟能力中台交互使用的Token

	}
	// GameConfig 游戏独立配置
	GameConfig struct {
		DynamicOddsEnable bool `yaml:"dynamicOddsEnable"` //是否开启动态赔率是否开启
	}
)

/**
 * SetConfigurationFileName
 * 设置配置文件名字
 *
 * @param fileName string - 配置文件名字
 * @return
 */

func SetConfigurationFileName(fileName string) {
	configurationFileName = fileName
}

/**
 * GetConfigurationFileName
 * 读取配置文件名字
 *
 * @param
 * @return string - 配置文件名字
 */

func GetConfigurationFileName() string {
	return configurationFileName
}

// parseContent 解析xml文件内容到Configuration结构体中
func parseContent(content string) error {
	var err error
	if ServerConf == nil {
		ServerConf = &Configuration{}
	}

	configMutex.Lock()
	defer configMutex.Unlock()
	if err = yaml.Unmarshal([]byte(content), ServerConf); nil != err {
		trace.Error("parseContent Unmarshal failed, error=%v", err.Error())
		return err
	}

	trace.Info("NacosClientInit Unmarshal conf=%+v", ServerConf)
	walkConfMember(ServerConf)
	trace.Info("NacosClientInit Unmarshal after env var replace conf=%+v", ServerConf)

	return err
}

// 遍历配置结构体成员并解析有变量默认值形式的成员
func walkConfMember(conf interface{}) {
	val := reflect.ValueOf(conf)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		trace.Error("walkConfMember conf is not a pointer to a struct")
		return
	}

	// 获取指针指向的实际结构体
	val = val.Elem()
	if val.Kind() == reflect.Ptr {
		trace.Error("walkConfMember val.Elem() is nil pointer")
		return
	}
	// 遍历结构体的所有字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.String && field.CanSet() {
			// 处理string类型字段
			strContent := field.String()
			if strings.Contains(strContent, "$") &&
				strings.Contains(strContent, "{") &&
				strings.Contains(strContent, "}") {
				field.SetString(parseConfLine(field.String()))
			}
		} else if field.Kind() == reflect.Struct {
			// 处理嵌套结构体
			walkConfMember(field.Addr().Interface())
		}
	}
}

// parseConfLine解析配置并从env中获取对应的值否则使用默认值 line格式为:${Var:DefaultValue} ${Var}
// 解析此种形式的配置 ${MQ_NAMESERVER:http://rocketmq-nameserver-dev.mq:9876} 或者${MQ_NAMESERVER}
func parseConfLine(line string) string {
	re := regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)?(:?([^}]+))?\}`)
	matches := re.FindStringSubmatch(line)
	trace.Info("parseConfLine line=%v matches=%v", line, matches)

	if len(matches) > 1 {
		envVar := matches[1] // var
		defaultValue := ""   // default
		if len(matches) > 3 && matches[3] != "" {
			defaultValue = matches[3]
		}
		realVal := getEnvWithDefault(envVar, defaultValue)
		trace.Info("parseConfLine line=%v, var=%v, defaultVal=%v, realVal=%v", line, envVar, defaultValue, realVal)
		return realVal
	}
	trace.Error("parseConfLine line=%v, no var and default val, please check line.", line)

	return line
}

// getEnvWithDefault 获取环境变量，如果不存在则返回默认值
func getEnvWithDefault(envKey, defaultVal string) string {
	strEnv := defaultVal
	if v, bIsExit := os.LookupEnv(envKey); bIsExit {
		strEnv = v
	}

	return strEnv
}

// getEnv 获取环境变量
func getEnv(envKey string) (bIsExit bool, val string) {
	val = ""

	val, bIsExit = os.LookupEnv(envKey)
	return bIsExit, val
}

// GetRedisAddr 获取redis集群信息
func GetRedisAddr() (addr []string, err error) {
	if ServerConf == nil {
		err = errors.New("ServerConf is nul")
		trace.Error("GetRedisAddr ServerConf == nil")
		return
	}

	addr = []string{ServerConf.RedisInfo.Host + ":" + ServerConf.RedisInfo.Port}
	trace.Info("GetRedisAddr addr=%v", addr)
	return
}

// RedisMode 工作模式
type RedisMode string

const (
	RedisModeSingle  RedisMode = "single"  //单节点工作模式
	RedisModeCluster RedisMode = "cluster" //集群工作模式
)

/**
 * GetRedisMode
 * redis工作模式 ClusterNodes 有值则是多节点工作模式 无值则为单节点工作模式
 *
 * @param
 * @return string - single:单节点工作模式 cluster:集群工作模式
 */

func GetRedisMode() RedisMode {
	if ServerConf == nil {
		trace.Error("GetRedisMode ServerConf == nil")
		return ""
	}

	mode := RedisModeSingle
	if RedisMode(ServerConf.RedisInfo.Mode) == RedisModeCluster {
		mode = RedisModeCluster
	}
	trace.Info("GetRedisMode ClusterNodes=%v, mode=%v", ServerConf.RedisInfo.Mode, mode)
	return mode
}

// GetRedisDb 获取redis集群信息
func GetRedisDb() int {
	if ServerConf == nil {
		trace.Error("GetRedisDb ServerConf == nil")
		return 0
	}

	db, _ := strconv.Atoi(ServerConf.RedisInfo.Db)
	return db
}

// GetKeyPrefix redis key的公用前缀
func GetKeyPrefix() string {
	if ServerConf == nil {
		trace.Error("GetKeyPrefix ServerConf == nil")
		return ""
	}

	return ServerConf.RedisInfo.KeyPrefix
}

// GetRedisUserInfo 获取redis用户
func GetRedisUserInfo() (string, string) {
	if ServerConf == nil {
		trace.Error("GetRedisUserInfo ServerConf == nil")
		return "", ""
	}

	return ServerConf.RedisInfo.UserName, ServerConf.RedisInfo.Password
}

// GetUidDbAliasName 以数据库名作为别名
func GetUidDbAliasName() string {
	if ServerConf == nil {
		trace.Error("GetMySqlDsnUidDb ServerConf == nil")
		return ""
	}

	return ServerConf.Database.UidDb.AliasName
}

// GetMySqlDsnGameDb 获取MySql信息
func GetMySqlDsnGameDb() string {
	if ServerConf == nil {
		trace.Error("GetMySqlDsnUidDb ServerConf == nil")
		return ""
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
		ServerConf.Database.UidDb.UserName, ServerConf.Database.UidDb.Password,
		ServerConf.Database.UidDb.Host, ServerConf.Database.UidDb.Database)

	return dsn
}

// GetMySqlGameDb 获取MySql信息
func GetMySqlGameDb() string {
	if ServerConf == nil {
		trace.Error("GetMySqlGameDb ServerConf == nil")
		return ""
	}
	gameDb := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
		ServerConf.Database.GameDb.UserName, ServerConf.Database.GameDb.Password,
		ServerConf.Database.GameDb.Host, ServerConf.Database.GameDb.Database)

	return gameDb
}

// GetGameDbAliasName 以数据库名作为别名
func GetGameDbAliasName() string {
	if ServerConf == nil {
		trace.Error("GetMySqlDsnUidDb ServerConf == nil")
		return ""
	}

	return ServerConf.Database.GameDb.AliasName
}

/**
 * GetGameId
 * 查询游戏Id
 *
 * @param
 * @return int64 - 游戏Id
 */

func GetGameId() int64 {
	if ServerConf == nil {
		trace.Error("GetGameId ServerConf == nil")
		return 0
	}

	return int64(ServerConf.Common.GameId)
}

// GetServerId 获取MySql信息
func GetServerId() int64 {
	if ServerConf == nil {
		trace.Error("GetServerId ServerConf == nil")
		return 0
	}

	configMutex.Lock()
	defer configMutex.Unlock()
	return ServerConf.ServerId
}

// SetServerId 获取MySql信息
func SetServerId(serverId int64) {
	if ServerConf == nil {
		trace.Error("SetServerId ServerConf == nil")
		return
	}
	trace.Info("SetServerId server id=%v", serverId)

	configMutex.Lock()
	ServerConf.ServerId = serverId
	configMutex.Unlock()
}

const suffix = "v1/"

// GetPlatformInfoUrl 获取平台host port信息
func GetPlatformInfoUrl() string {
	if ServerConf == nil {
		trace.Error("GetPlatformInfo ServerConf == nil")
		return ""
	}

	//host需要以 http:// 或者 https:// 开头
	host := ServerConf.Platform.Host
	if !strings.HasPrefix(host, "http://") &&
		!strings.HasPrefix(host, "https://") {
		host = "http://" + host
	}

	strUrl := ""
	if ServerConf.Platform.Port == 0 {
		strUrl = host
	} else {
		strUrl = fmt.Sprintf("%v:%v", host, ServerConf.Platform.Port)
	}
	//if !strings.HasSuffix(strUrl, "/") {
	//	strUrl = strUrl + "/"
	//}

	return strUrl
}

// GetHeartbeatInterval 获取心跳配置信息
func GetHeartbeatInterval() (interval int) {
	if ServerConf == nil {
		trace.Error("GetPlatformInfo ServerConf == nil")
		return -1
	}

	interval = ServerConf.Common.HeartbeatInterval
	if interval <= 0 {
		interval = 1
	}

	return
}

// GetMySqlDsnUidDb 获取MySql信息
func GetMySqlDsnUidDb() string {
	if ServerConf == nil {
		trace.Error("GetMySqlDsnUidDb ServerConf == nil")
		return ""
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
		ServerConf.Database.UidDb.UserName, ServerConf.Database.UidDb.Password,
		ServerConf.Database.UidDb.Host, ServerConf.Database.UidDb.Database)

	return dsn
}

// GetLoopTaskExpired 获取心跳配置信息
func GetLoopTaskExpired() (expired int) {
	if ServerConf == nil {
		trace.Error("GetPlatformInfo ServerConf == nil")
		return -1
	}
	expired = ServerConf.Common.LoopTaskExpired
	if expired <= 0 {
		expired = 10
	}
	return
}

// GetHeartbeatExpired 获取心跳配置信息
func GetHeartbeatExpired() (expired int) {
	if ServerConf == nil {
		trace.Error("GetPlatformInfo ServerConf == nil")
		return -1
	}

	expired = ServerConf.Common.HeartbeatExpired
	if expired <= 0 {
		expired = 6 //单位s秒
	}

	return
}

// GetHttpRetryInfo 获取http重试次数和重试间隔
func GetHttpRetryInfo() (retryTimes int, retryInterval int) {
	if ServerConf == nil {
		trace.Error("GetPlatformInfo ServerConf == nil")
		return -1, -1
	}

	retryTimes = ServerConf.Platform.RetryTime
	retryInterval = ServerConf.Platform.RetryInterval
	if retryTimes <= 0 {
		retryTimes = 5
	}
	if retryInterval <= 0 {
		retryInterval = 200 //单位ms
	}

	return
}

// GetRocketMQNameServer 获取mq地址
func GetRocketMQNameServer() []string {
	if ServerConf == nil {
		trace.Error("GetRocketMQNameServer ServerConf == nil")
		return nil
	}
	nameSvr := ServerConf.Rocketmq.NameServer

	return []string{nameSvr}
}

/**
 * GetGameServerEventGroup
 * 获取玩家进入房间mq topic组名
 *
 * @return string - 组名
 */

func GetJoinMessageRoomGroup() string {
	if ServerConf == nil {
		trace.Error("GetJoinMessageRoomGroup ServerConf == nil")
		return ""
	}

	return ServerConf.Rocketmq.JoinMessageRoomTopic
}

/**
 * GetLeaveMessageRoomGroup
 * 获取玩家离开房间mq topic组名
 *
 * @return string - 组名
 */

func GetLeaveMessageRoomGroup() string {
	if ServerConf == nil {
		trace.Error("GetLeaveMessageRoomGroup ServerConf == nil")
		return ""
	}

	return ServerConf.Rocketmq.LeaveMessageRoomTopic
}

// GetRocketMQRetries 获取mq地址
func GetRocketMQRetries() int {
	if ServerConf == nil {
		trace.Error("GetRocketMQRetries ServerConf == nil")
		return 2
	}
	return ServerConf.Rocketmq.Retries
}

func GetRocketMQQueueMaxLen() int {
	if ServerConf == nil {
		trace.Error("GetRocketMQQueueMaxLen ServerConf == nil")
		return 500
	}

	maxLen := ServerConf.Rocketmq.ProducerQueueMaxLen
	if maxLen <= 0 {
		maxLen = 500
	}

	return maxLen
}

// GetHttpConnectTimeout 获取http连接超时时间单位秒
func GetHttpConnectTimeout() time.Duration {
	if ServerConf == nil {
		trace.Error("GetHttpConnectTimeout ServerConf == nil")
		return 5
	}

	connectTimeout := ServerConf.Http.HttpConnectTimeout
	if connectTimeout <= 0 {
		connectTimeout = 5
	}

	return time.Duration(connectTimeout) * time.Second
}

// GetHttpReadWriteTimeout 获取http连接超时时间单位秒
func GetHttpReadWriteTimeout() time.Duration {
	if ServerConf == nil {
		trace.Error("GetHttpReadWriteTimeout ServerConf == nil")
		return 5
	}

	readWriteTimeout := ServerConf.Http.HttpReadWriteTimeout
	if readWriteTimeout <= 0 {
		readWriteTimeout = 5
	}

	return time.Duration(readWriteTimeout) * time.Second
}

type Switch string

const (
	switchOn  Switch = "on"
	switchOff Switch = "off"
)

// GetUserLimitSwitch 获取个人限红校验开关 返回是否校验个人限红
func GetUserLimitSwitch() bool {
	if ServerConf == nil {
		trace.Error("GetUserLimitSwitch ServerConf == nil")
		return false
	}

	bIsOn := false
	if switchOn == Switch(ServerConf.Common.UserLimitSwitch) {
		bIsOn = true
	}

	return bIsOn
}

// LogLevel 日志等级类型
type LogLevel string

const (
	LogLevelEmergency     LogLevel = "Emergency"
	LogLevelAlert         LogLevel = "Alert"
	LogLevelCritical      LogLevel = "Critical"
	LogLevelError         LogLevel = "Error"
	LogLevelWarning       LogLevel = "Warning"
	LogLevelNotice        LogLevel = "Notice"
	LogLevelInformational LogLevel = "Informational"
	LogLevelDebug         LogLevel = "Debug"
)

/**
 * GetLogLevel
 * 获取日志等级
 *
 * @param
 * @return int - 日志等级
 */

func GetLogLevel() int {
	level := trace.LevelInformational
	if ServerConf == nil {
		trace.Error("GetLogLevel ServerConf == nil")
		return level
	}
	trace.Info("GetLogLevel log level=%v", ServerConf.Common.LogLevel)

	switch LogLevel(ServerConf.Common.LogLevel) {
	case LogLevelEmergency:
		level = trace.LevelEmergency
	case LogLevelAlert:
		level = trace.LevelAlert
	case LogLevelCritical:
		level = trace.LevelCritical
	case LogLevelError:
		level = trace.LevelError
	case LogLevelWarning:
		level = trace.LevelWarning
	case LogLevelNotice:
		level = trace.LevelNotice
	case LogLevelInformational:
		level = trace.LevelInformational
	case LogLevelDebug:
		level = trace.LevelDebug
	default:
		trace.Error("GetLogLevel no log level=%v", ServerConf.Common.LogLevel)
	}

	return level
}
