package topolvm

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)


type LvmdConfig struct {
	SocketName string	`yaml:"socket-name"`
	DeviceClasses []struct {
		Name string			`yaml:"name"`
		VolumeGroup string	`yaml:"volume-group"`
		SpareGb int			`yaml:"spare-gb"`
		Default bool		`yaml:"default"`
	} `yaml:"device-classes"`
}


func LoadLvmdConfig(lvmdConfigFile string) (*LvmdConfig, error ){
	content, err := ioutil.ReadFile(lvmdConfigFile)
	if err != nil {
		return nil, err
	}
	var config LvmdConfig
	err = yaml.UnmarshalStrict(content, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}


