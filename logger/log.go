package logger

import "os"

// 支持级别
// OFF、FATAL、ERROR、WARN、INFO、DEBUG、TRACE、 ALL
// color fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "testPrintColor", 0x1B)
// 前景(字体颜色) 背景 颜色
// ---------------------------------------
// 30 40 黑色
// 31 41 红色
// 32 42 绿色
// 33 43 黄色
// 34 44 蓝色
// 35 45 紫红色
// 36 46 青蓝色
// 37 47 白色
//
// 代码 意义
// -------------------------
// 0 终端默认设置
// 1 高亮显示
// 4 使用下划线
// 5 闪烁
// 7 反白显示
// 8 不可见

// 日志颜色分类  绿色、蓝色、白色、黄色、红色、青蓝色
const (
	LoggerLevelAll int = iota
	LoggerLevelTrace
	LoggerLevelDebug
	LoggerLevelInfo
	LoggerLevelWarn
	LoggerLevelError
	LoggerLevelFatal
	LoggerLevelOff
)

type Logger struct {
	Output *os.File
	Level  int
}

func NewLogger() *Logger {
	l := new(Logger)
	l.Output = os.Stdout
	return l
}

// 设置日志等级
func (l *Logger) SetLoggerLevel(level int) {
	if level >= LoggerLevelAll && level <= LoggerLevelOff {
		l.Level = level
	}
}

// 设置输出目标
func (l *Logger) SetOutput(output *os.File) {
	if output != nil {
		l.Output = output
	}
}

// 强制退出
func (*Logger) Fatal(v ...interface{}) {

}
