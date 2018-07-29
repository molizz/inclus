package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

const (
	DefConfigName = "includes.yml"
)

var (
	definitions = make(map[string]interface{})
	versions    = make(map[string]interface{})
)

type Command struct {
	version string
	config  string
}

func initConfig(c *Command) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	viper.SetConfigType("yaml")
	body, err := ioutil.ReadFile(c.config)
	if err != nil {
		return err
	}

	err = viper.ReadConfig(bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	definitions = viper.Get("definitions").(map[string]interface{})
	versions = viper.Get("versions").(map[string]interface{})
	return nil
}

func buildDownloads(ver string) ([]string, error) {
	version, ok := versions[ver]
	if !ok {
		return nil, errors.New("not found " + ver)
	}

	// 该版本下的所有的依赖
	versionsMap := version.(map[string]interface{})

	// build url
	downloads := make([]string, 0)

	for defName, verName := range versionsMap {
		defKeyMap, ok := definitions[defName].(map[string]interface{})
		if !ok {
			return nil, errors.New("wrong format")
		}
		url := defKeyMap["url"].(string)
		url = strings.Replace(url, "{{version}}", verName.(string), -1)
		downloads = append(downloads, url)
	}

	return downloads, nil
}

func download(urls []string) error {
	return nil
}

var versionCmd = &cobra.Command{
	Use:   "version VERSION_NAME CONFIG",
	Short: "传入版本号",
	Long:  "传入配置文件, 默认使用当前目录的includes.yml",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := verify(args...)
		if err != nil {
			return err
		}

		c := prepare(args...)
		err = initConfig(c)
		if err != nil {
			return err
		}

		urls, err := buildDownloads(c.version)
		if err != nil {
			return err
		}

		if len(urls) == 0 {
			return errors.New("Not match download url.")
		}

		err = download(urls)
		if err != nil {
			return err
		}

		return nil
	},
}

func prepare(args ...string) *Command {
	c := &Command{}
	c.version = args[0]
	if len(args) >= 2 {
		c.config = args[1]
	} else {
		c.config = DefConfigName
	}
	return c
}

func verify(args ...string) error {
	if len(args) == 0 {
		return errors.New("version is not set.")
	}

	existFunc := func(path string) error {
		configInfo, err := os.Stat(path)
		if os.IsNotExist(err) || configInfo.IsDir() {
			return errors.New("config file is not exist.")
		}
		return nil
	}

	if len(args) < 2 {
		err := existFunc(DefConfigName)
		return err
	}

	err := existFunc(args[1])
	return err
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
