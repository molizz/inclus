package cmd

import (
	"fmt"
	"os"

	"github.com/bangwork/bang-api/app/utils/errors"
	"github.com/molizz/inclus/bin/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ConfigFile = "inclus.yaml"
)

var (
	ErrConfigNotFount = errors.New("inclus file is not found")
)

var (
	config = func() *viper.Viper {
		c, err := utils.GetViper(ConfigFile)
		if err != nil {
			panic(err)
		}
		return c
	}()
)

var RootCmd = &cobra.Command{
	Use:   "inclus -g web v2.0 app v3.0",
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
