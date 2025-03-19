package log

import (
	"fmt"
	blog "github.com/beego/beego/v2/core/logs"
	"runtime"
	"sl.framework.com/resource/conf"
	"sl.framework.com/trace"
	"time"
)

/**
 * @Author: M
 * @Date: 2024/7/31 16:47
 * @Desc:
 */

func init() {

	//initBLogs()
	f := &trace.PatternLogFormatter{
		Pattern:    "%w %t %f:%n => %m" + fmt.Sprintf(`%s`, newline()),
		WhenFormat: time.DateTime + ".000000 -0700",
	}
	logger(trace.LevelInfo)
	trace.RegisterFormatter("pattern", f)
	_ = trace.SetGlobalFormatter("pattern")
	trace.EnableFuncCallDepth(true)
	conf.SetResetLogger(logger)
}

func logger(level int) {
	trace.SetLevel(level)
	//setLogger(
	//	trace.AdapterFile,
	//	fmt.Sprintf(`{"filename":"%s.log","level":%d,"maxlines":%d,"maxsize":%d,"daily":%t,"maxdays":%d}`,
	//		conf.SectionDefault("beego", "appName", "dragon-tiger-resource"), level,
	//		1000000, 1<<28, true, 30))
	setLogger(trace.AdapterConsole, fmt.Sprintf(`{"level":%d,"color":true}`, level))
}

func initBLogs() {
	f := &blog.PatternLogFormatter{
		Pattern:    "%w %t %f:%n => %m" + fmt.Sprintf(`%s`, newline()),
		WhenFormat: time.DateTime + ".000000 -0700",
	}
	blog.RegisterFormatter("pattern", f)
	setLogger(
		blog.AdapterFile,
		fmt.Sprintf(`{"filename":"%s.log","level":%d,"maxlines":%d,"maxsize":%d,"daily":%t,"maxdays":%d,"color":true}`,
			conf.SectionDefault("beego", "appName", "roulette-resource"), trace.LevelDebug,
			1000000, 1<<28, true, 30))
	setLogger(trace.AdapterConsole, `{"level":7,"color":true}`)
	blog.EnableFuncCallDepth(true)
}

func setLogger(adapter string, config ...string) {
	_ = blog.SetLogger(adapter, config...)
}

// newline获取当前操作系统的换行符
func newline() string {
	switch runtime.GOOS {
	case "windows":
		return "\r\n"
	default:
		return "\n"
	}
}
