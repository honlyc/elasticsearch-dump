package config

import (
	"fmt"
)
import "github.com/spf13/viper"

const defaultConfigFile = "config.yaml"

var (
	CONFIG Config
)

type Config struct {
	Es Es `json:"es"`
}

func init() {
	v := viper.New()
	v.SetConfigFile(defaultConfigFile)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if err := v.Unmarshal(&CONFIG); err != nil {
		fmt.Println(err)
	}
	fmt.Println(CONFIG)
}

type Es struct {
	Cluster string `mapstructure:"cluster" json:"cluster" yaml:"cluster"`
	Name    string `json:"name"`
	Port    int    `json:"port"`

	IndexName string `json:"index_name"`

	Size  int    `json:"size"`
	Query string `json:"query"`
}
