package cmd

import (
	"github.com/spf13/cobra"
)

type Download struct {
	Dir string
}

// 只支持github

func init() {
	RootCmd.AddCommand(downloadCmd)
}

var downloadCmd = &cobra.Command{
	Use:   "d",
	Short: "传入版本号下载",
	Long:  "传入版本号, 下载对应的软件包",
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}
