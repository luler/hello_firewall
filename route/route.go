package route

import (
	"gin_base/app/controller/common"
	"gin_base/app/helper/response_helper"
	"gin_base/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(e *gin.Engine) {
	e.NoRoute(func(context *gin.Context) {
		response_helper.Common(context, 404, "路由不存在")
	})
	api := e.Group("/api")
	api.GET("/test", common.Test)
	//ip管理
	api.POST("/banIp", common.BanIp)
	api.POST("/unBanIp", common.UnBanIp)
	api.POST("/getBanIpList", common.GetBanIpList)

	//登录相关
	auth := api.Group("", middleware.Auth())
	auth.POST("/test_auth", common.Test)
}
