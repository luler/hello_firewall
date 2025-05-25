package main

import (
	"fmt"
	"gin_base/app"
	"gin_base/bin"
	_ "gin_base/docs"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	//项目初始化
	app.InitApp(app.InitTypeBase, app.InitTypeCron, app.InitTypeMigrate)
}

// @title 接口文档
// @version 1.0
// @description 当前页面用于展示项目一些开放的接口
// @termsOfService http://swagger.io/terms/
// @contact.name 开发人员
// @contact.url https://cas.luler.top/
// @contact.email 1207032539@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host
// @BasePath
func main() {
	cmd := &cobra.Command{
		Use:   "myapp",
		Short: "主程序入口",
		Long:  "主程序入口，启动程序或者执行自定义命令",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("请使用子命令，或添加 --help 查看帮助")
		},
	}
	///////////////////
	//自定义命令开始
	///////////////////
	cmd.AddCommand(bin.ServeCommand())   //启动Gin服务命令
	cmd.AddCommand(bin.DebugCommand())   //调试专用
	cmd.AddCommand(bin.MigrateCommand()) //数据库迁移

	///////////////////
	//自定义命令结束
	///////////////////

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
