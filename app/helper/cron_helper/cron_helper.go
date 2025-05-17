package cron_helper

import (
	"gin_base/app/middleware"
	"github.com/gogits/cron"
)

func InitCron() {
	c := cron.New()
	c.AddFunc("定时清理ip限制缓存", "0 */1 * * * ?", func() {
		middleware.ClearIpRateLimit()
	})

	c.Start()
}
