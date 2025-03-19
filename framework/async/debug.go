package async

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	maxStack  = 20
	separator = "---------------------------------------\n"
)

// GetCallStack 获取调用栈信息并将其存储为二维数组。
//
// 参数:
//   - skip: int, 从调用栈中跳过的层级数。skip=0 表示从当前函数的调用者开始，skip=1 表示再上一级调用者，以此类推。
//   - max: int, 获取调用栈的最大层数。指定要获取的最大调用栈深度。
//
// 返回值:
//
//	[][]string: 返回一个二维数组，每一行包含三个元素，分别是：
//	  - 函数名 (string)
//	  - 文件名 (string)
//	  - 行号 (string)
//
// 示例:
//   - GetCallStack(0, 10) // 获取从当前调用函数开始的最多 10 层的调用栈信息。
//
// 返回的二维数组格式示例:
//
//	[
//	  {"main.main", "/path/to/file/main.go", "20"},
//	  {"runtime.main", "/usr/local/go/src/runtime/proc.go", "255"},
//	  {"runtime.goexit", "/usr/local/go/src/runtime/asm_amd64.s", "1581"}
//	]
func GetCallStack(skip, max int) [][]string {
	stackTrace := make([]uintptr, max)       // 获取调用栈的最大层数
	n := runtime.Callers(2+skip, stackTrace) // 跳过前两层（当前函数和 runtime.Callers 本身）

	// 定义二维数组，用于存储调用栈信息，每行包含函数名、文件名和行号
	var callStack [][]string

	for i := 0; i < n; i++ {
		pc := stackTrace[i]
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		// 将函数名、文件名和行号存入切片
		callInfo := []string{
			filepath.Base(fn.Name()), // 函数名
			file,                     // 文件名
			fmt.Sprintf("%d", line),  // 行号转换为字符串
		}
		// 将调用栈信息切片添加到二维数组中
		callStack = append(callStack, callInfo)
	}

	return callStack
}

func stack() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s\ntraceback:\n", separator))

	i := 1
	for {
		pc, file, line, ok := runtime.Caller(i)
		if !ok || i > maxStack {
			break
		}
		funcName := "<unknown>"
		if f := runtime.FuncForPC(pc); f != nil {
			funcName = f.Name()
		}
		sb.WriteString(fmt.Sprintf("    stack: %d %v [file: %s] [func: %s] [line: %d]\n", i, ok, file, funcName, line))
		i++
	}

	sb.WriteString(separator)
	return sb.String()
}

func TryException() {
	errs := recover()
	if errs == nil {
		return
	}

	exeName := os.Args[0]                                           // 获取程序名称
	pid := os.Getpid()                                              // 获取进程ID
	now := time.Now().UTC().Format("20060102150405")                // 设定时间格式
	fileName := fmt.Sprintf("%s_%d_%s-dump.log", exeName, pid, now) // 保存错误信息文件名:程序名-进程ID-当前时间（年月日时分秒）
	fmt.Printf("dump to file: %s\n", fileName)

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("failed to create dump trace file: %v\n", err)
		return
	}
	defer file.Close()
	_stack := stack()

	if _, writeErr := file.WriteString(fmt.Sprintf("%v\r\n", errs)); writeErr != nil { //输出panic信息
		fmt.Printf("Failed to write error message to trace file: %v\n", writeErr)
	}
	if _, writeErr := file.WriteString("========\r\n"); writeErr != nil {
		fmt.Printf("Failed to write separator to trace file: %v\n", writeErr)
	}
	if _, writeErr := file.WriteString(_stack); writeErr != nil { //输出堆栈信息
		fmt.Printf("Failed to write stack info to trace file: %v\n", writeErr)
	}
	// 输出到控制台
	fmt.Println(_stack)
}
