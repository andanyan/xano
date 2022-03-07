package logger

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

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

// 日志颜色分类  绿色、蓝色、白色、黄色、红色、紫红色
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

// 等级设置
type LevelItem struct {
	Color int
	Desc  string
}

var levelConfig = map[int]*LevelItem{
	LoggerLevelAll: {
		Color: 37,
		Desc:  "LOG",
	},
	LoggerLevelTrace: {
		Color: 32,
		Desc:  "TRACE",
	},
	LoggerLevelDebug: {
		Color: 34,
		Desc:  "DEBUG",
	},
	LoggerLevelInfo: {
		Color: 37,
		Desc:  "INFO",
	},
	LoggerLevelWarn: {
		Color: 33,
		Desc:  "WARN",
	},
	LoggerLevelError: {
		Color: 31,
		Desc:  "ERROR",
	},
	LoggerLevelFatal: {
		Color: 35,
		Desc:  "FATAL",
	},
	LoggerLevelOff: {
		Color: 35,
		Desc:  "OFF",
	},
}

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

// 内容自适应
func (l *Logger) write(level int, needTrace bool, msg string) {
	if l.Level > level {
		return
	}
	levelItem := levelConfig[level]

	timeDesc := time.Now().Format("2006-01-02:15:04:05.000000")

	var fileDesc string
	_, file, line, ok := runtime.Caller(2)
	if ok {
		fileIndex := strings.Index(file, "/src/")
		if fileIndex >= 0 {
			file = file[fileIndex+5:]
			fileIndex := strings.Index(file, "/")
			if fileIndex >= 0 {
				file = file[fileIndex+1:]
			}
		}
		fileDesc = fmt.Sprintf("%s:%d", file, line)
	}
	if needTrace {
		msg += "\n" + string(debug.Stack())
	}
	s := fmt.Sprintf("%c[0;40;%dm%s%c[0m [%s]: %s %s\n", 0x1B, levelItem.Color, levelItem.Desc, 0x1B, timeDesc, fileDesc, msg)
	l.Output.WriteString(s)
}

// 追踪
func (l *Logger) Trace(v ...interface{}) {
	l.write(LoggerLevelTrace, true, fmt.Sprint(v...))
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.write(LoggerLevelTrace, true, fmt.Sprintf(format, v...))
}

// 调试
func (l *Logger) Debug(v ...interface{}) {
	l.write(LoggerLevelDebug, false, fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.write(LoggerLevelDebug, false, fmt.Sprintf(format, v...))
}

// 警告
func (l *Logger) Info(v ...interface{}) {
	l.write(LoggerLevelInfo, false, fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.write(LoggerLevelInfo, false, fmt.Sprintf(format, v...))
}

// 警告
func (l *Logger) Warn(v ...interface{}) {
	l.write(LoggerLevelWarn, false, fmt.Sprint(v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.write(LoggerLevelWarn, false, fmt.Sprintf(format, v...))
}

// 报错
func (l *Logger) Error(v ...interface{}) {
	l.write(LoggerLevelError, true, fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.write(LoggerLevelError, true, fmt.Sprintf(format, v...))
}

// 强制退出
func (l *Logger) Fatal(v ...interface{}) {
	l.write(LoggerLevelFatal, true, fmt.Sprint(v...))
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.write(LoggerLevelFatal, true, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// 打印
func (l *Logger) Print(v ...interface{}) {
	l.write(LoggerLevelAll, false, fmt.Sprint(v...))
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.write(LoggerLevelAll, false, fmt.Sprintf(format, v...))
}
