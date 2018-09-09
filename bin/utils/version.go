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
	Ver.versions = make(map[string]map[string]string)
	err := Ver.Prepare()
	if err != nil {
		fmt.Println(err)
	}
}

type Version struct {
	versions map[string]map[string]string
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
		v.versions[fullVer] = temp
	}
	return nil
}

func (v *Version) AddVer(ver map[string]map[string]string) *Version {
	for name, version := range ver {
		v.versions[name] = version
	}
	return v
}

func (v *Version) Save() error {
	cfg, err := GetViper(VerionsFile)
	if err != nil {
		return err
	}

	cfg.SetConfigFile(VerionsFile)

	cfg.Set("versions", v.versions)
	err = cfg.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}
