// RpcLogs
package RpcLogs

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/astaxie/beego/logs"
)

const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

const (
	AdapterConsole   = "console"
	AdapterFile      = "file"
	AdapterMultiFile = "multifile"
	AdapterMail      = "smtp"
	AdapterConn      = "conn"
	AdapterEs        = "es"
	AdapterJianLiao  = "jianliao"
	AdapterSlack     = "slack"
)

var Loger *logs.BeeLogger = logs.NewLogger()

type RpcLoger struct {
}

func NewRpcLogs(filename string) *RpcLoger {
	InitLoger(filename)
	return &RpcLoger{}
}

func InitLoger(filename string) {
	Loger.SetLogger("console")
	config := make(map[string]interface{})
	config["filename"] = filename
	config["maxdays"] = 30
	data, _ := json.Marshal(config)
	Loger.SetLogger(logs.AdapterFile, string(data))
	Loger.EnableFuncCallDepth(true)
	depth := Loger.GetLogFuncCallDepth()
	Loger.SetLogFuncCallDepth(depth + 1)
	Loger.SetLevel(logs.LevelInfo)
	//Loger.Async(100)

}

func (l *RpcLoger) SetLogger(adapter string, config ...string) error {
	err := Loger.SetLogger(adapter, config...)
	return err
}

func (l *RpcLoger) EnableFuncCallDepth(b bool) {
	Loger.EnableFuncCallDepth(b)
}

func (l *RpcLoger) SetLogFuncCallDepth(d int) {
	depth := Loger.GetLogFuncCallDepth()
	Loger.SetLogFuncCallDepth(depth + 1 + d)
}

func (l *RpcLoger) SetLevel(level int) {
	Loger.SetLevel(level)
}

func (l *RpcLoger) Printf(format string, v ...interface{}) {
	Loger.Info(fmt.Sprintf(format, v...))
}

func (l *RpcLoger) Print(v ...interface{}) { Loger.Info(fmt.Sprint(v...)) }

func (l *RpcLoger) Println(v ...interface{}) { Loger.Info(fmt.Sprintln(v...)) }

func (l *RpcLoger) Fatal(v ...interface{}) {
	Loger.Error(fmt.Sprint(v...))
	Loger.Flush()
	os.Exit(1)
}

func (l *RpcLoger) Fatalf(format string, v ...interface{}) {
	Loger.Error(fmt.Sprintf(format, v...))
	Loger.Flush()
	os.Exit(1)
}

func (l *RpcLoger) Fatalln(v ...interface{}) {
	Loger.Error(fmt.Sprintln(v...))
	Loger.Flush()
	os.Exit(1)
}

func (l *RpcLoger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	Loger.Error("[Panic]", s)
	Loger.Flush()
	panic(s)
}

func (l *RpcLoger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	Loger.Error("[Panic]", s)
	Loger.Flush()
	panic(s)
}

func (l *RpcLoger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	Loger.Error("[Panic]", s)
	Loger.Flush()
	panic(s)
}

func (l *RpcLoger) Alert(v ...interface{}) {
	Loger.Alert(fmt.Sprintln(v...))
}

func (l *RpcLoger) Critical(v ...interface{}) {
	Loger.Critical(fmt.Sprintln(v...))
}

func (l *RpcLoger) Error(v ...interface{}) {
	Loger.Error(fmt.Sprintln(v...))
}

func (l *RpcLoger) Warning(v ...interface{}) {
	Loger.Warning(fmt.Sprintln(v...))
}

func (l *RpcLoger) Warn(v ...interface{}) {
	Loger.Warn(fmt.Sprintln(v...))
}

func (l *RpcLoger) Notice(v ...interface{}) {
	Loger.Notice(fmt.Sprintln(v...))
}

func (l *RpcLoger) Informational(v ...interface{}) {
	Loger.Informational(fmt.Sprintln(v...))
}

func (l *RpcLoger) Info(v ...interface{}) {
	Loger.Info(fmt.Sprintln(v...))
}

func (l *RpcLoger) Debug(v ...interface{}) {
	Loger.Debug(fmt.Sprintln(v...))
}

func (l *RpcLoger) Trace(v ...interface{}) {
	Loger.Trace(fmt.Sprintln(v...))
}

func SetLogger(adapter string, config ...string) error {
	err := Loger.SetLogger(adapter, config...)
	return err
}

func EnableFuncCallDepth(b bool) {
	Loger.EnableFuncCallDepth(b)
}

func SetLogFuncCallDepth(d int) {
	depth := Loger.GetLogFuncCallDepth()
	Loger.SetLogFuncCallDepth(depth + 1 + d)
}

func SetLevel(l int) {
	Loger.SetLevel(l)
}

func Printf(format string, v ...interface{}) {
	Loger.Info(fmt.Sprintf(format, v...))
}

func Print(v ...interface{}) { Loger.Info(fmt.Sprint(v...)) }

func Println(v ...interface{}) { Loger.Info(fmt.Sprintln(v...)) }

func Fatal(v ...interface{}) {
	Loger.Error(fmt.Sprint(v...))
	Loger.Flush()
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	Loger.Error(fmt.Sprintf(format, v...))
	Loger.Flush()
	os.Exit(1)
}

func Fatalln(v ...interface{}) {
	Loger.Error(fmt.Sprintln(v...))
	Loger.Flush()
	os.Exit(1)
}

func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	Loger.Error("[Panic]", s)
	Loger.Flush()
	panic(s)
}

func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	Loger.Error("[Panic]", s)
	Loger.Flush()
	panic(s)
}

func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	Loger.Error("[Panic]", s)
	Loger.Flush()
	panic(s)
}

func Alert(v ...interface{}) {
	Loger.Alert(fmt.Sprintln(v...))
}

func Critical(v ...interface{}) {
	Loger.Critical(fmt.Sprintln(v...))
}

func Error(v ...interface{}) {
	Loger.Error(fmt.Sprintln(v...))
}

func Warning(v ...interface{}) {
	Loger.Warning(fmt.Sprintln(v...))
}

func Warn(v ...interface{}) {
	Loger.Warn(fmt.Sprintln(v...))
}

func Notice(v ...interface{}) {
	Loger.Notice(fmt.Sprintln(v...))
}

func Informational(v ...interface{}) {
	Loger.Informational(fmt.Sprintln(v...))
}

func Info(v ...interface{}) {
	Loger.Info(fmt.Sprintln(v...))
}

func Debug(v ...interface{}) {
	Loger.Debug(fmt.Sprintln(v...))
}

func Trace(v ...interface{}) {
	Loger.Trace(fmt.Sprintln(v...))
}
