package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/viper"
)

func IsExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// isDirectory determines if a file represented
// by `path` is a directory or not
func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func ReadAll(filePath string) (string, error) {
	f, err := os.Open(filePath)
	check(err)
	defer f.Close()

	content, err := io.ReadAll(f) // 全部読み込んでくれる
	check(err)

	return string(content), nil
}
func WriteAll(filePath, content string) {
	f, err := os.Create(filePath)
	check(err)
	defer f.Close()

	_, err = f.WriteString(content)
	check(err)
	f.Sync()
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

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
