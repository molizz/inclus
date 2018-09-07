package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

type Upload struct {
	PushURL string `yaml:"push_url"`
	Branch  string `yaml:"branch"`
	Path    string `yaml:"push_to"`
}
type Inclus struct {
	Upload      Upload            `yaml:"upload"`
	CloneDir    string            `yaml:"clone_dir"`
	Definitions map[string]string `yaml:"definitions"`
}

func GetViper(cfgPath string) (*viper.Viper, error) {
	versionCfg := viper.New()
	versionCfg.SetConfigType("yaml")

	body, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	err = versionCfg.ReadConfig(bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return versionCfg, nil
}

func fileExist(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || info.IsDir() {
		return false
	}
	return true
}
