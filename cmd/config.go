package cmd

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	EChars      map[string]string `yaml:"echars"`
	RuleDir     string            `yaml:"ruledir"`
	Params      []string          `yaml:"params"`
	Indent      string            `yaml:"indent"`
	Prefix      string            `yaml:"prefix"`
	RedisOption RedisOption       `mapstructure:"redis"`
}

var config AppConfig

func LoadConfig(configFilePath string, cfg any) bool {
	if !IsExists(configFilePath) {
		log.Printf("Not found configfile:%s", configFilePath)
		return false
	}

	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err == nil {
		viper.Unmarshal(cfg)
		log.Printf("Read config OK. file: %s", viper.ConfigFileUsed())
		return true
	} else {
		log.Fatal(err)
		return false
	}
}
