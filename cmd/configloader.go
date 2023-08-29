package cmd

import (
	"log"

	"github.com/spf13/viper"
)

type IConfigLoader interface {
	Load(configFilePath string) bool
}
type RegexTemplates struct {
	Templates []RegexTemplate `mapstructure:"templates"`
}
type TemplateManager struct{}

// ===============================================

// ===============================================
var ttManager TemplateManager

// ===============================================

// load delegate
func (manager *TemplateManager) Execute(configFilePath string, hand IConfigLoader) bool {
	if manager == nil {
		manager = &TemplateManager{}
	}
	return hand.Load(configFilePath)
}

// Load kind=templates
func (ts *AppConfig) Load(configFilePath string) bool {
	return loadAny(ts, configFilePath)
}

// Load kind=templates
func (ts *RegexTemplates) Load(configFilePath string) bool {
	return loadAny(ts, configFilePath)
}

// Load kind=template
func (ts *RegexTemplate) Load(configFilePath string) bool {
	return loadAny(ts, configFilePath)
}

// Load kind=regex-rules
func (ts *RegexRules) Load(configFilePath string) bool {
	return loadAny(ts, configFilePath)
}

// Load kind=regex-rule
func (ts *RegexRule) Load(configFilePath string) bool {
	return loadAny(ts, configFilePath)
}

// Load kind=regex-rule
func (ts *CheckRules) Load(configFilePath string) bool {
	return loadAny(ts, configFilePath)
}

// Load kind=regex-rules
func (ts *CheckRule) Load(configFilePath string) bool {
	return loadAny(ts, configFilePath)
}

// loadAny kind=templates
func loadAny(ts any, configFilePath string) bool {
	viper.SetConfigFile(configFilePath)

	if err := viper.ReadInConfig(); err == nil {
		viper.Unmarshal(&ts)
		log.Printf("Read config OK. file: %s", viper.ConfigFileUsed())
		return true
	} else {
		log.Fatal(err)
		return false
	}
}
