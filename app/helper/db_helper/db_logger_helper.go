package db_helper

import (
	"context"
	"fmt"
	"gin_base/app/helper/log_helper"
	"gorm.io/gorm/logger"
	"time"
)

// 自定义 Logger
type DbLogger struct {
	logger.Interface
}

func (dl *DbLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *dl
	newLogger.Interface = dl.Interface.(logger.Interface).LogMode(level)
	return &newLogger
}

func (dl *DbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc()
	//异步处理
	go func() {
		errMsg := "无"
		code := 0
		if err != nil {
			code = 500
			errMsg = err.Error()
		}
		sqlInfo := fmt.Sprintf("%s 耗时[%dms] 错误[%s]", sql, time.Since(begin).Milliseconds(), errMsg)
		var logData []map[string]interface{}
		logData = append(logData, map[string]interface{}{
			"code":       code,
			"other_id":   ctx.Value("common_log_other_id"),
			"message":    sqlInfo,
			"waste_time": fmt.Sprintf("%.6f", float64(time.Since(begin).Seconds())),
		})

		log_helper.QueueCommonLog(logData)
	}()
}
