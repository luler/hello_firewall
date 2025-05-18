package middleware

import (
	"fmt"
	"gin_base/app/helper/exception_helper"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

// 创建速率限制器
type IpRateLimitStruct struct {
	Limiter   *rate.Limiter
	UpdatedAt time.Time
}

var limiters = make(map[string]*IpRateLimitStruct)

// 限制速度
func IpRateLimit(r float64, b int) gin.HandlerFunc {
	return func(context *gin.Context) {
		ip := context.ClientIP()
		key := fmt.Sprintf("%v_%v_%v", ip, r, b)
		limiter, exist := limiters[key]
		if !exist {
			limiter = &IpRateLimitStruct{
				Limiter:   rate.NewLimiter(rate.Limit(r), b),
				UpdatedAt: time.Now(),
			}
			limiters[key] = limiter
		}
		if limiter.Limiter.Allow() {
			limiter.UpdatedAt = time.Now()
			context.Next()
		} else {
			exception_helper.CommonException("请求过于频繁，请稍后重试", http.StatusTooManyRequests)
		}
	}
}

// 定时清理ip限制缓存
func ClearIpRateLimit() {
	for key, limiter := range limiters {
		if time.Since(limiter.UpdatedAt) > time.Minute {
			delete(limiters, key)
		}
	}
}
