/*
 * @Author: xwu
 * @Date: 2021-12-26 18:45:38
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-21 14:30:03
 */
package Logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var BaseLogger *zap.Logger
var TerminalSlientLogger *zap.Logger

func init() {
	InitBaseLogger()
	InitTerminalSlientLogger()
}

func InitBaseLogger() {
	var err error
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		// EncodeTime:    zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
	}

	// 设置日志级别
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)

	config := zap.Config{
		Level:         atom,          // 日志级别
		Development:   true,          // 开发模式，堆栈跟踪
		Encoding:      "console",     // 输出格式 console 或 json
		EncoderConfig: encoderConfig, // 编码器配置
		// InitialFields:    map[string]interface{}{"serviceName": "trade"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout", "./normal.log"}, // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}

	// 构建日志
	BaseLogger, err = config.Build()
	if err != nil {
		panic(fmt.Sprintf("Uber logger initial failed: %v", err))
	}
}

func InitTerminalSlientLogger() {
	var err error
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		// EncodeTime:    zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
	}

	// 设置日志级别
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)

	config := zap.Config{
		Level:         atom,          // 日志级别
		Development:   true,          // 开发模式，堆栈跟踪
		Encoding:      "console",     // 输出格式 console 或 json
		EncoderConfig: encoderConfig, // 编码器配置
		// InitialFields:    map[string]interface{}{"serviceName": "trade"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"./trade.log"}, // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}

	// 构建日志
	TerminalSlientLogger, err = config.Build()
	if err != nil {
		panic(fmt.Sprintf("Uber logger initial failed: %v", err))
	}
}
