package route

import (
	"gin_base/app/controller"
	"gin_base/app/controller/common"
	"gin_base/app/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func InitRouter(e *gin.Engine) {
	//e.NoRoute(func(context *gin.Context) {
	//	response_helper.Common(context, 404, "路由不存在")
	//})
	//前端路由
	e.Static("/helloFirewall", "./web/dist/helloFirewall")
	e.NoRoute(func(context *gin.Context) {
		context.File("./web/dist/helloFirewall/index.html")
	})
	e.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/helloFirewall/ipRule")
	})
	e.GET("/api/README.md", func(context *gin.Context) {
		context.File("./README.md")
	})
	//接口
	api := e.Group("/api")
	api.GET("/test", common.Test)
	//swagger页面
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//登录相关
	api.POST("/login", middleware.IpRateLimit(0.1, 5), controller.Login)
	api.GET("/casLogin", controller.CasLogin)
	auth := api.Group("", middleware.Auth())
	auth.POST("/test_auth", common.Test)
	//ip管理
	auth.POST("/banIp", common.BanIp)
	auth.POST("/unBanIp", common.UnBanIp)
	auth.POST("/changeStatus", common.ChangeStatus)
	auth.GET("/getBanIpList", common.GetBanIpList)
	//用户各相关
	auth.POST("/resetPassword", controller.ResetPassword)
	auth.GET("/getUserInfo", controller.GetUserInfo)

}
