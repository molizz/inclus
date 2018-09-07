package cmd

import (
	"errors"
	"fmt"
	"github.com/molizz/inclus/bin/utils"
	"github.com/spf13/cobra"
	"os"
)

type version = string
type defineName = string

var inclusCfg *Inclus

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "g web v2.0 app v3.0",
	Short: "传入版本号",
	Long:  "传入配置文件, 默认使用当前目录的 " + ConfigFile,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 格式化参数
		argsMap, err := parseArgs(args)
		if err != nil {
			return err
		}

		// 加载inclus.yaml
		config, err := GetViper(ConfigFile)
		if err != nil {
			return err
		}

		// 校验参数(key是否在配置是否存在)
		definitions := config.Get("definitions")
		err = validParse(definitions.(map[string]interface{}), argsMap)
		if err != nil {
			return err
		}

		// clone 仓库 & 统计commit数
		versionWithCommitTotalMap := make(map[version]uint)
		for name, ver := range argsMap {
			cloner := utils.NewClone(os.Stderr, inclusCfg.Definitions[name], inclusCfg.CloneDir)
			err = cloner.Clone()
			if err != nil {
				return err
			}

			t, err := cloner.CommitCountByTagName(ver)
			if err != nil {
				return err
			}
			versionWithCommitTotalMap[ver] = t
		}

		return nil
	},
}

// 验证参数
func validParse(definitions map[string]interface{}, argsMap map[defineName]version) error {
	for k, _ := range argsMap {
		if _, ok := definitions[k]; !ok {
			return fmt.Errorf("参数 %s 没有在inclus.yaml中定义. 请定义后再试.", k)
		}
	}
	return nil
}

// 将参数格式化
func parseArgs(args []string) (map[defineName]version, error) {
	if len(args)%2 != 0 {
		return nil, errors.New("args参数数量不对")
	}

	argsMap := make(map[defineName]version)

	k := ""
	for i, a := range args {
		if i%2 == 0 {
			k = a
		} else {
			argsMap[k] = a
		}
	}
	return argsMap, nil
}
