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

//	func LoadRule(fullPath string) (bool, *RuleConfig) {
//		var rule RuleConfig
//		ok := LoadConfig(fullPath, &rule)
//		return ok, &rule
//	}
func (rule *RuleConfig) findByName(ruleName string) bool {
	for _, dirPath := range config.RuleDirs {
		if rule.findRule(dirPath, ruleName) {
			return true
		}
	}
	return false
}
func (rule *RuleConfig) findRule(dirPath string, ruleName string) bool {
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
		return LoadConfig(fullPath, &rule)
	}
	return false
}
