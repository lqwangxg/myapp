package cmd

import (
	"os"
	"path/filepath"
)

type AppContext struct {
	*AppConfig
	*RegexTemplates
	*RegexRules
	*CheckRules
}

func NewContext(config *AppConfig) *AppContext {
	return &AppContext{
		AppConfig:      config,
		RegexTemplates: &RegexTemplates{Templates: make([]RegexTemplate, 0)},
		RegexRules:     &RegexRules{Rules: make([]RegexRule, 0)},
		CheckRules:     &CheckRules{Rules: make([]CheckRule, 0)},
	}
}

func (ctx *AppContext) LoadDirectory(rootPath string) {

	files, err := os.ReadDir(rootPath)
	if err != nil {
		return
	}
	for _, file := range files {
		fullPath := filepath.Join(rootPath, file.Name())
		ok, err := IsDir(fullPath)
		if err == nil {
			if ok {
				ctx.LoadDirectory(fullPath)
			} else {
				ctx.LoadFile(fullPath)
			}
		}
	}
}
