package db

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"sl.framework.com/resource/conf"
	"sl.framework.com/resource/httpserver/beego/log"
	"sl.framework.com/trace"
	"strings"
	"time"
)

/**
 * @Author: M
 * @Date: 2024/8/2 15:18
 * @Desc:
 */

const section = "database"

// 参数1        数据库的别名，用来在 ORM 中切换数据库使用
// 参数2        driverName
// 参数3        对应的链接字符串
// 参数4(可选)  设置最大空闲连接
// 参数5(可选)  设置最大数据库连接 (go >= 1.2)
func init() {

	username := conf.Section(section, "user")
	password := conf.Section(section, "password")
	host := conf.Section(section, "host")
	dbName := conf.Section(section, "dbname")
	alias := conf.Section(section, "alias")
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", username, password, host, dbName)

	trace.Info("[%s] 数据库连接: %s", alias, dataSource)
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	maxIdle := 30
	maxConn := 30
	// 设置为 UTC 时间
	orm.DefaultTimeLoc = time.UTC
	// 开启SQL日志
	orm.Debug = true
	orm.DebugLog = orm.NewLog(&log.NoOpWriter{}) // 不输出ORM日志
	orm.LogFunc = dbLog                          // 输出格式化日志
	err := orm.RegisterDataBase(alias, "mysql", dataSource, orm.MaxIdleConnections(maxIdle), orm.MaxOpenConnections(maxConn))
	if err != nil {
		trace.Error("[%s] 数据库初始化失败: %s", alias, err.Error())
	} else {
		trace.Notice("[%s] 数据库初始化成功", alias)
	}
}

func dbLog(query map[string]interface{}) {
	sql := parseSQL(query["sql"].(string))
	if sql == "COMMIT" || sql == "START TRANSACTION" || sql == "ROLLBACK" {
		trace.Notice("[%s/SQL/%.2fms]: %s", query["flag"], query["cost_time"], sql)
	} else {
		trace.Info("[%s/SQL/%.2fms]: %s", query["flag"], query["cost_time"], sql)
	}
	if value, exists := query["err"]; exists {
		trace.Error("[执行有误]: %s", value.(error).Error())
	}
}

func parseSQL(input string) string {

	// 解析字符串为 SQL 部分和参数部分
	parts := strings.SplitN(input, "-", 2)
	if len(parts) != 2 {
		fmt.Println("错误：输入格式不正确，缺少'-'分隔符")
		return input
	}
	sqlPart := parts[0]
	paramsPart := parts[1]

	// 正则表达式用于匹配可能的 JSON 参数
	re := regexp.MustCompile(`\{.*?\}`)
	matches := re.FindAllString(paramsPart, -1)

	// 替换 JSON 字符串以便安全分割其他参数
	tempParamsPart := paramsPart
	placeholders := []string{}
	for i, match := range matches {
		placeholder := fmt.Sprintf("__JSON_PLACEHOLDER_%d__", i)
		placeholders = append(placeholders, placeholder)
		tempParamsPart = strings.Replace(tempParamsPart, match, placeholder, 1)
	}

	// 定义正则表达式，匹配反引号中的内容
	re = regexp.MustCompile("`([^`]*)`")

	// 查找所有匹配的内容
	matchess := re.FindAllStringSubmatch(tempParamsPart, -1)
	var params []string
	for _, match := range matchess {
		params = append(params, fmt.Sprintf("`%s`", match[1]))
	}
	if len(matchess) <= 0 {
		// 分割参数
		params = strings.Split(tempParamsPart, ",")
	}

	// 恢复 JSON 字符串到参数列表
	for i, placeholder := range placeholders {
		for j, param := range params {
			if strings.Contains(param, placeholder) {
				// 只能用``包裹 不然需要处理不定量的单双引号
				params[j] = fmt.Sprintf("`%s`", matches[i])
				break
			}
		}
	}

	// 替换 SQL 语句中的问号占位符为实际参数
	paramCount := strings.Count(sqlPart, "?")
	if len(params) != paramCount && paramsPart != "``" {
		trace.Warning("错误：参数数量不匹配，需要 %d 个参数，但提供了 %d 个，返回 ORM原生SQL", paramCount, len(params))
		return input
	}

	for _, param := range params {
		// 替换 SQL 语句中的第一个出现的问号为参数
		sqlPart = strings.Replace(sqlPart, "?", param, 1)
	}

	trace.Debug("最终的 SQL 语句: %s", sqlPart)
	return sqlPart
}
