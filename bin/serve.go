package bin

import (
	"gin_base/app/middleware"
	"gin_base/route"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
)

func ServeCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "启动Gin服务",
		Run: func(cmd *cobra.Command, args []string) {
			//开启gin服务
			StartServer()
		},
	}

	return cmd
}

// 开启gin服务
func StartServer() {
	gin.SetMode(os.Getenv(gin.EnvGinMode))
	//不输出请求日志
	//gin.DefaultWriter = ioutil.Discard

	engine := gin.Default()
	//初始化中间件
	middleware.InitMiddleware(engine)
	//初始化路由
	route.InitRouter(engine)

	engine.Run() // listen and serve on 0.0.0.0:8080
}
