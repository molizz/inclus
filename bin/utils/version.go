package utils

import (
	"fmt"
)

const (
	VerionsFile = "versions.yaml"
)

var Ver *Version

func init() {
	Ver = new(Version)
	Ver.Versions = make(map[string]map[string]string)
	err := Ver.Prepare()
	if err != nil {
		fmt.Println(err)
	}
}

type Version struct {
	Versions map[string]map[string]string
}

func (v *Version) Prepare() error {
	if !FileExist(VerionsFile) {
		err := FileCreate(VerionsFile)
		if err != nil {
			return err
		}
	}

	cfg, err := GetViper(VerionsFile)
	if err != nil {
		return err
	}

	configMap := cfg.GetStringMap("versions")

	for fullVer, nameWithVer := range configMap {
		temp := make(map[string]string)
		for name, ver := range nameWithVer.(map[string]interface{}) {
			temp[name] = fmt.Sprintf("%v", ver)
		}
		v.Versions[fullVer] = temp
	}
	return nil
}

func (v *Version) AddVer(ver map[string]map[string]string) *Version {
	for name, version := range ver {
		v.Versions[name] = version
	}
	return v
}

func (v *Version) Save() error {
	cfg, err := GetViper(VerionsFile)
	if err != nil {
		return err
	}

	cfg.SetConfigFile(VerionsFile)

	cfg.Set("versions", v.Versions)
	err = cfg.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}
