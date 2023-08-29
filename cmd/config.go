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
		return appContext.AppConfig.Load(configFile)
	case "templates":
		templates := &RegexTemplates{Templates: make([]RegexTemplate, 0)}
		loaded = templates.Load(configFile)
		appContext.RegexTemplates.Templates = append(appContext.RegexTemplates.Templates, templates.Templates...)
	case "regex-rules":
		rules := &RegexRules{Rules: make([]RegexRule, 0)}
		loaded = rules.Load(configFile)
		appContext.RegexRules.Rules = append(appContext.RegexRules.Rules, rules.Rules...)
	}
	return loaded
}
