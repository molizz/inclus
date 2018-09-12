package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/molizz/inclus/bin/utils"
	"github.com/spf13/cobra"
)

type RepoConfig struct {
	CloneDir    string            `yaml:"clone_dir"`
	Definitions map[string]string `yaml:"definitions"`
}

type version = string
type defineName = string

var repoCfg *RepoConfig

func init() {
	RootCmd.AddCommand(versionCmd)

	repoCfg = new(RepoConfig)
	repoCfg.Definitions = make(map[string]string)
	repoCfg.CloneDir = config.GetString("clone_dir")
	defs := config.GetStringMap("definitions")
	for k, v := range defs {
		repoCfg.Definitions[k] = v.(string)
	}
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

		// 校验参数(key是否在配置是否存在)
		definitions := config.Get("definitions")
		err = validParse(definitions.(map[string]interface{}), argsMap)
		if err != nil {
			return err
		}

		// clone 仓库 & 统计commit数
		versionWithCommitTotalMap := make(map[version]uint)
		for name, ver := range argsMap {
			cloner := utils.NewClone(os.Stderr, repoCfg.Definitions[name], repoCfg.CloneDir)
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

		// 生成版本号(v100.12.123)
		newVer := make([]string, 0)
		for _, ver := range argsMap {
			verNum := versionWithCommitTotalMap[ver]
			newVer = append(newVer, fmt.Sprintf("%d", verNum))
		}
		newVerStr := "v" + strings.Join(newVer, ".")

		// 保存
		newVerMap := make(map[string]map[string]string)
		for name, ver := range argsMap {
			if newVerMap[newVerStr] == nil {
				newVerMap[newVerStr] = make(map[string]string)
			}
			newVerMap[newVerStr][name] = ver
		}
		err = utils.Ver.AddVer(newVerMap).Save()
		if err != nil {
			return err
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
