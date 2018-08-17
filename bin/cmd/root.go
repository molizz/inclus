package cmd

import (
	"fmt"
	"os"

	"github.com/bangwork/bang-api/app/utils/errors"
	"github.com/spf13/cobra"
)

const (
	ConfigName = "inclus.yaml"
)

var (
	ErrConfigNotFount = errors.New("config file is not found")
)

var RootCmd = &cobra.Command{
	Use:   "inclus version v2.0 " + ConfigName,
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
