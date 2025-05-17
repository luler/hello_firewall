package log_helper

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logHelper *logrus.Logger

// 初始化日志助手
func InitlogHelper() {
	logHelper = logrus.New()
	// 设置日志级别为 Info
	logHelper.SetLevel(logrus.InfoLevel)
	//设置日志格式
	logHelper.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true, // 显示完整的时间戳
	})
	// 创建一个新的 lumberjack.Logger 实例
	logFilePath := "./runtime/logs/app.log"
	hook := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    50,  // 单位：MB
		MaxAge:     365, // 保留时间：天
		MaxBackups: 100, // 最大备份数量
	}

	// 设置日志输出到 hook
	logHelper.SetOutput(hook)
}

// 写日志
func Info(args ...interface{}) {
	logHelper.Info(args...)
}
func Error(args ...interface{}) {
	logHelper.Error(args...)
}
func Warning(args ...interface{}) {
	logHelper.Warning(args...)
}
func Debug(args ...interface{}) {
	logHelper.Debug(args...)
}
func Fatal(args ...interface{}) {
	logHelper.Fatal(args...)
}
