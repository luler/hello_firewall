package bin

import (
	"github.com/spf13/cobra"
)

func DebugCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "debug",
		Short: "调试专用",
		Run: func(cmd *cobra.Command, args []string) {
			handle()
		},
	}

	return cmd
}

func handle() {
	//
}
