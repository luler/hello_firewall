package bin

import (
	"gin_base/app"
	"github.com/spf13/cobra"
)

func MigrateCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "数据库迁移",
		Run: func(cmd *cobra.Command, args []string) {
			// 自动创建表
			app.InitApp("migrate")
		},
	}

	return cmd
}
