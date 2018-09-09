package utils

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

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

func FileExist(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || info.IsDir() {
		return false
	}
	return true
}

func FileCreate(file string) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}
