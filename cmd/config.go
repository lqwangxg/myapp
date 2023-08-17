package cmd

import (
	"log"
)

func (appContext *AppContext) LoadFile(configFile string) bool {
	if !IsExists(configFile) {
		log.Printf("Not found configfile:%s", configFile)
		return false
	}
	content, err := ReadAll(configFile)
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
		return ttManager.Execute(configFile, appContext.AppConfig)
	case "templates":
		loaded = ttManager.Execute(configFile, appContext.RegexTemplates)
	case "template":
		hand := &RegexTemplate{}
		appContext.RegexTemplates.Templates = append(appContext.RegexTemplates.Templates, *hand)
		loaded = ttManager.Execute(configFile, hand)
	case "regex-rules":
		loaded = ttManager.Execute(configFile, appContext.RegexRules)
	case "regex-rule":
		hand := &RegexRule{}
		appContext.RegexRules.Rules = append(appContext.RegexRules.Rules, *hand)
		loaded = ttManager.Execute(configFile, hand)
	case "check-rules":
		loaded = ttManager.Execute(configFile, appContext.CheckRules)
	case "check-rule":
		hand := &CheckRule{}
		appContext.CheckRules.Rules = append(appContext.CheckRules.Rules, *hand)
		loaded = ttManager.Execute(configFile, hand)
	}
	return loaded
}
