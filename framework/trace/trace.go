package trace

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

/*
	1.对beego中的log做个封装，外部直接使用uclog即可，不用每次使用再创建一个log对象
	2.todo:日志中统一加入goroutineid
*/

var (
	logOnce sync.Once
	trace   *BeeLogger
)

// getGoRoutineID 获取协程ID
func getGoRoutineID() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	_, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return idField
}

// RFC5424 trace message levels.
//onst (
//	LevelEmergency = iota
//	LevelAlert
//	LevelCritical
//	LevelError
//	LevelWarning
//	LevelNotice
//	LevelInformational
//	LevelDebug
//)

/*
LoggerInit 初始化log，使用once保证只初始化一次
*/
func LoggerInit() (err error) {
	err = errors.New("trace initialized already")
	logOnce.Do(func() {
		if nil == trace {
			trace = NewLogger()

			//设置Console输出
			if err = trace.SetLogger(AdapterConsole); err != nil {
				fmt.Printf("LoggerInit AdapterConsole error=%v", err.Error())
				return
			}

			EnableFuncCallDepth(true)
			//设置文件输出
			jsonConf := fmt.Sprintf(`{
					"filename": "server.log",
					"daily": false,
					"hourly": true,
					"maxhours":%d,
					"level":%d
				}`, 24*15, 6)
			trace.Info("[%v] LoggerInit file config json=%v", getGoRoutineID(), jsonConf)
			if err = trace.SetLogger(AdapterFile, jsonConf); err != nil {
				trace.Info("LoggerInit AdapterFile error=%v", err.Error())
				return
			}
		}
	})
	return
}
