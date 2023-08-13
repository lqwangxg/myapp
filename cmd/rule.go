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

func LoadRules(dirPath string, rules map[string]RuleConfig) {
	files, err := os.ReadDir(dirPath)
	check(err)

	r := regexp.MustCompile(`\w+\.(yaml|yml)$`)
	for _, file := range files {
		fullPath := filepath.Join(dirPath, file.Name())
		isdir, err := IsDir(fullPath)
		if err != nil {
			log.Fatal(err)
		}
		if isdir {
			LoadRules(fullPath, rules)
			continue
		}
		//skip ifnot yaml or yml
		if !r.MatchString(fullPath) {
			continue
		}

		//=============================
		if ok, rule := LoadRule(fullPath); ok {
			rules[rule.Name] = rule
		}
	}
}
func LoadRule(fullPath string) (bool, RuleConfig) {
	var rule RuleConfig
	ok := LoadConfig(fullPath, &rule)
	return ok, rule
}
