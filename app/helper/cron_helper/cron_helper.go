package cron_helper

import (
	"gin_base/app/helper/db_helper"
	"gin_base/app/logic"
	"gin_base/app/middleware"
	"gin_base/app/model"
	"github.com/gogits/cron"
	"time"
)

func InitCron() {
	c := cron.New()
	c.AddFunc("定时清理ip限制缓存", "0 */1 * * * ?", func() {
		middleware.ClearIpRateLimit()
	})

	c.AddFunc("定时关闭过期的ip封禁规则", "0 */1 * * * ?", func() {
		result := db_helper.Db().Model(&model.IPRule{}).
			Where("status = 1").
			Where("expired_at > 0 and expired_at <= ?", time.Now().Unix()).
			Updates(map[string]interface{}{"status": 0, "expired_at": 0})
		// 只有在有规则更新时才重建iptables规则
		if result.RowsAffected > 0 {
			//忽略报错
			logic.GetIPTablesManager().RebuildRules()
		}
	})

	c.Start()
}
