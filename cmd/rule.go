/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/viper"
)

func loadRules(dirPath string, rules map[string]RuleConfig) {
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
			loadRules(fullPath, rules)
			continue
		}
		//skip ifnot yaml or yml
		if !r.MatchString(fullPath) {
			continue
		}
		if rule := loadRule(fullPath); rule.Name != "" {
			rules[rule.Name] = rule
		}
	}
}

func loadRule(ruleFile string) RuleConfig {
	var rule RuleConfig
	if IsExists(ruleFile) {
		viper.SetConfigFile(ruleFile)
		err := viper.ReadInConfig()
		if err == nil {
			viper.Unmarshal(&rule)
		} else {
			check(err)
		}

	}
	return rule
}
