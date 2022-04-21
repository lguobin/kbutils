package alarm

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/lguobin/kbutils/log"
)

type erralarm uint

type errString struct{}

type errorInfo struct {
	STime    string      `json:"startTime"`
	Alarm    string      `json:"alarm"`
	Message  interface{} `json:"message"`
	Filename string      `json:"filename"`
	Funcname string      `json:"funcname"`
	Line     int         `json:"line"`
}

// 定义具体的日志级别常量
const (
	WX erralarm = iota + 1
	SMS
	EMAIL
	LOGGER
)

func NewAlarm() *errString {
	return &errString{}
}

// 写一个根据传进来的Level 获取对应的字符串
func geterrLevelStr(level erralarm) string {
	switch level {
	case WX:
		return "WX"
	case SMS:
		return "SMS"
	case EMAIL:
		return "EMAIL"
	case LOGGER:
		return "Logger"
	default:
		return "Logger"
	}
}

func (e *errString) process(el erralarm, args ...interface{}) {
	currentTime := time.Now()                   // 获取当前时间
	fileName, line, functionName := "?", 0, "?" // 定义\ 文件名、行号、方法名
	pc, fileName, line, ok := runtime.Caller(2)

	if ok {
		functionName = runtime.FuncForPC(pc).Name()
		functionName = filepath.Ext(functionName)
		functionName = strings.TrimPrefix(functionName, ".")
	}

	elevel := geterrLevelStr(el)
	msg := errorInfo{
		STime:    currentTime.String(),
		Alarm:    elevel,
		Message:  args,
		Filename: fileName,
		Funcname: functionName,
		Line:     line,
	}

	jsons, errs := json.Marshal(msg)
	if errs != nil {
		fmt.Println("json marshal error:", errs)
	}
	errorJsonInfo := string(jsons)
	// fmt.Println("模拟控制台打印: ", errorJsonInfo)

	if elevel == "EMAIL" {
		// 执行发邮件

	} else if elevel == "SMS" {
		// 执行发短信

	} else if elevel == "WX" {
		// 执行发微信

	} else if elevel == "LOGGER" {
		// 执行记日志
		log.Debug("实时打印 >> ", errorJsonInfo)

	} else {
		// 缺省值 - 执行记日志
		log.Debug("实时打印 >> ", errorJsonInfo)

	}
}

// 保存到日志
func (e *errString) Logger(text ...interface{}) { e.process(LOGGER, text) }

// 发邮件
func (e *errString) Email(text ...interface{}) { e.process(SMS, text) }

// 发短信
func (e *errString) Sms(text ...interface{}) { e.process(SMS, text) }

// 发微信
func (e *errString) WeChat(text ...interface{}) { e.process(WX, text) }
