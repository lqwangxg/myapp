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

	isMatched, result := NewRegexText(PATTERN_KIND_KEY, content).GetMatchResult(true, false)
	if !isMatched {
		return false
	}

	kind := result.FirstMatch().Params["key"]
	loaded := false
	switch kind {
	case "app":
		return ttManager.Execute(configFile, appContext.AppConfig)
	case "templates":
		loaded = ttManager.Execute(configFile, appContext.RegexTemplates)
	case "regex-rules":
		loaded = ttManager.Execute(configFile, appContext.RegexRules)
	case "check-rules":
		loaded = ttManager.Execute(configFile, appContext.CheckRules)
	}
	return loaded
}
