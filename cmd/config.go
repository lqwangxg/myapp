package cmd

import (
	"log"
)

func (appContext *AppContext) LoadConfig(configFilePath string) bool {
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
	loaded := false
	switch kind {
	case "app":
		return ttManager.Execute(configFilePath, appContext.AppConfig)
	case "templates":
		loaded = ttManager.Execute(configFilePath, appContext.RegexTemplates)
	case "template":
		hand := &RegexTemplate{}
		appContext.RegexTemplates.Templates = append(appContext.RegexTemplates.Templates, *hand)
		loaded = ttManager.Execute(configFilePath, hand)
	case "regex-rules":
		loaded = ttManager.Execute(configFilePath, appContext.RegexRules)
	case "regex-rule":
		hand := &RegexRule{}
		appContext.RegexRules.Rules = append(appContext.RegexRules.Rules, *hand)
		loaded = ttManager.Execute(configFilePath, hand)
	case "check-rules":
		loaded = ttManager.Execute(configFilePath, appContext.CheckRules)
	case "check-rule":
		hand := &CheckRule{}
		appContext.CheckRules.Rules = append(appContext.CheckRules.Rules, *hand)
		loaded = ttManager.Execute(configFilePath, hand)
	}
	return loaded
}
