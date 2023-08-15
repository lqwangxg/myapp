package cmd

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig(configFilePath string, cfg any) bool {
	if !IsExists(configFilePath) {
		log.Printf("not found configfile:%s", configFilePath)
		return false
	}

	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err == nil {
		viper.Unmarshal(cfg)
		log.Printf("Read config file OK:%s\n=↓↓↓======\n%v\n=↑↑↑=======", viper.ConfigFileUsed(), cfg)
		return true
	} else {
		log.Fatal(err)
		return false
	}
}
