package main

import (
	"fmt"
	"gin_base/app"
	"gin_base/bin"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	//项目初始化
	app.InitApp("base")
}

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
