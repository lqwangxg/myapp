/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func (rs *RegexRules) GetDefaultRule() *RegexRule {
	for _, r := range rs.Rules {
		if r.Name == "default" {
			return &r
		}
	}
	return nil
}
func (rs *RegexRules) GetRule(name string) *RegexRule {
	for _, r := range rs.Rules {
		if r.Name == name {
			return &r
		}
	}
	return nil
}

//	func LoadRule(fullPath string) (bool, *RuleConfig) {
//		var rule RuleConfig
//		ok := LoadConfig(fullPath, &rule)
//		return ok, &rule
//	}
func (rule *RegexRule) findByName(ruleName string) bool {
	return rule.findRule(config.RuleDir, ruleName)
}
func (rule *RegexRule) findRule(dirPath string, ruleName string) bool {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return false
	}

	r := regexp.MustCompile(ruleName + `\.(yaml|yml)$`)
	for _, file := range files {
		fullPath := filepath.Join(dirPath, file.Name())
		isdir, err := IsDir(fullPath)
		if err != nil {
			log.Fatal(err)
		}
		if isdir {
			if rule.findRule(fullPath, ruleName) {
				return true
			}
			continue
		}
		//skip ifnot destfile of yaml or yml
		if !r.MatchString(file.Name()) {
			continue
		}
		//=============================
		//TODO
		//return LoadConfig(fullPath, &rule)
	}
	return false
}
