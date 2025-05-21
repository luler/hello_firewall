package app

import (
	"gin_base/app/helper/cron_helper"
	"gin_base/app/helper/db_helper"
	"gin_base/app/helper/log_helper"
	"gin_base/app/model"
	"github.com/joho/godotenv"
)

const (
	InitTypeBase    string = "base"
	InitTypeCron    string = "cron"
	InitTypeMigrate string = "migrate"
)

// 项目启动初始化
func InitApp(initTypes ...string) {
	for _, s := range initTypes {
		switch s {
		case InitTypeBase:
			//加载.env配置
			godotenv.Load()
			//初始化日志记录方式
			log_helper.InitlogHelper()
		case InitTypeCron:
			//初始化定时任务
			cron_helper.InitCron()
		case InitTypeMigrate:
			// 自动创建表
			db_helper.Db().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci").AutoMigrate(
				&model.User{},
				&model.IPRule{},
			)
		}
	}

}
