package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "inclus version v2.0 includes.yml",
	Short: "版本依赖管理工具",
	Long:  `版本依赖管理工具`,
}

func Execute() error {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return nil
}
