package cmd

import (
	"log"
	"os"
	"path/filepath"
)

type AppConfig struct {
	EChars      map[string]string `yaml:"echars"`
	RuleDir     string            `yaml:"ruledir"`
	Params      map[string]string `yaml:"params"`
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
	content, err := ReadAll(configFilePath)
	if err != nil {
		return false
	}

	result := NewRegexText(kindkey, content).Match()
	if result == nil || result.MatchCount == 0 {
		return false
	}

	kind := result.FirstMatch().Params["key"]
	switch kind {
	case "app":
		return ttManager.Execute(configFilePath, &config)
	case "templates":
		hand := &RegexTemplates{
			Templates: make([]RegexTemplate, 0),
		}
		return ttManager.Execute(configFilePath, hand)
	case "template":
		hand := &RegexTemplate{}
		return ttManager.Execute(configFilePath, hand)
	case "regex-rules":
		hand := &RegexRules{
			Rules: make([]RegexRule, 0),
		}
		return ttManager.Execute(configFilePath, hand)
	case "regex-rule":
		hand := &RegexRule{}
		return ttManager.Execute(configFilePath, hand)
	case "check-rules":
		hand := &CheckRules{
			Rules: make([]CheckRule, 0),
		}
		return ttManager.Execute(configFilePath, hand)
	case "check-rule":
		hand := &CheckRule{}
		return ttManager.Execute(configFilePath, hand)
	}
	return false
}

func LoadAllConfigs(rootPath string) {

	files, err := os.ReadDir(rootPath)
	if err != nil {
		return
	}
	for _, file := range files {
		fullPath := filepath.Join(rootPath, file.Name())
		ok, err := IsDir(fullPath)
		if err == nil {
			if ok {
				handler := NewRegexDirectory(flags.RuleName, fullPath)
				reger.Execute(handler)
			} else {
				handler := NewRegexFile(flags.RuleName, fullPath)
				reger.Execute(handler)
			}
		}
	}
}
