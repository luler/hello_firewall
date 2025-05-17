package middleware

import (
	"encoding/json"
	"fmt"
	"gin_base/app/helper/helper"
	"gin_base/app/helper/log_helper"
	"gin_base/app/helper/request_helper"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"github.com/syyongx/php2go"
	"os"
	"strings"
	"time"
)

// 限制速度
func CommonLog() gin.HandlerFunc {
	return func(context *gin.Context) {
		//判断是否开启
		if os.Getenv("COMMON_LOG_ENABLE") != "true" {
			// 处理请求
			context.Next()
			return
		}
		other_id := ksuid.New().String()
		context.Set("common_log_other_id", other_id)
		// 开始时间
		startTime := time.Now()
		// 处理请求
		context.Next()
		// 结束时间
		endTime := time.Now()
		go func() {
			//日志等级
			level := "info"
			// 执行时间
			latencyTime := endTime.Sub(startTime)
			// 计算执行时间（秒），精确到小数点后四位
			if latencyTime.Seconds() >= 10 {
				level = "warning"
			}
			wasteTime := fmt.Sprintf("%.4f", float64(latencyTime.Seconds()))
			// 请求路由
			reqUri := context.Request.URL.Path
			// 请求IP
			clientIP := context.ClientIP()
			// 服务器IP
			serverIP := helper.GetFirstServerIP()
			// 状态码
			statusCode := context.Writer.Status()
			if php2go.InArray(statusCode, []int{500}) {
				level = "error"
			}
			//返回json数据
			var response_message string
			if data, ok := context.Get("response_data"); ok {
				switch data.(type) {
				case map[string]interface{}:
					data := data.(map[string]interface{})
					if _, hasCode := data["code"]; hasCode {
						statusCode = data["code"].(int)
					}
					jsonBytes, _ := json.Marshal(data)
					response_message = string(jsonBytes)
				default:
					response_message = data.(string)
				}
			}
			//请求参数
			request_param := request_helper.Input(context)
			//请求头
			headers := make(map[string]string)
			for k, v := range context.Request.Header {
				headers[k] = strings.Join(v, ", ")
			}
			message, _ := json.Marshal(gin.H{
				"header": headers,
				"param":  request_param,
			})
			// 日志数据
			var logData []map[string]interface{}
			logData = append(logData, map[string]interface{}{
				"level":       level, // 你可以根据需要设置不同的级别
				"code":        statusCode,
				"url":         reqUri,
				"waste_time":  wasteTime,
				"message":     string(message),
				"other":       response_message,
				"other_id":    other_id,
				"create_time": startTime.UnixMilli(),
				"client_ip":   clientIP,
				"server_ip":   serverIP,
			})

			log_helper.QueueCommonLog(logData)
			log_helper.PushCommonLog()
		}()
	}
}
