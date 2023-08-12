/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rules = make(map[string]ReplaceRuleConfig)

// replaceCmd represents the replaceTF command
var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "replace text file or text files under directory by pattern name",
	Long: ` replace text file or text files under directory by pattern name. 
	      pattern name connects to a json or configMap which includes rules of pattern/replacement/skipRules.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("regex replace called")
		fmt.Printf("name=%s,  textfile=%s, directory=%s, suffix=%s \n", flags.RuleName, flags.DestFile, flags.DestDir, flags.IncludeSuffix)
		// if flags.RuleName == "" {
		// 	fmt.Printf("--name=%s is required.", flags.RuleName)
		// }
		if flags.DestFile == "" && flags.DestDir == "" {
			fmt.Printf("--textfile=%s and --directory=%s can't be empty neither.", flags.DestFile, flags.DestDir)
		}

		pattern := flags.Pattern
		if pattern == "" {
			return
		}

		for _, dir := range config.RuleDirs {
			loadRules(dir, rules)
		}
		rslog(rules)

		rs := NewRegex(pattern)
		rs.Action = ReplaceAction
		if flags.Content != "" {
			rs.ReplaceText(flags.Content)
		}
		if flags.DestFile != "" {
			rs.ProcFile(flags.DestFile)
		}
		if flags.DestDir != "" {
			rs.ProcDir(flags.DestDir)
		}

	},
}

func rslog(rs map[string]ReplaceRuleConfig) {
	for key, rule := range rs {
		log.Printf("key:%s, rule.name: %s , rule.pattern:%s", key, rule.Name, rule.MatchPattern)
	}
}
func loadRules(dirPath string, rules map[string]ReplaceRuleConfig) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
		return
	}
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
		if rule, err := loadRule(fullPath); err == nil {
			rule.Key = fullPath
			rules[rule.Key] = rule
		}
	}
}
func loadRule(ruleFile string) (ReplaceRuleConfig, error) {
	viper.SetConfigFile(ruleFile)

	var rule ReplaceRuleConfig
	err := viper.ReadInConfig()
	if err == nil {
		viper.Unmarshal(&rule)
		//log.Printf("Using ruleFile:%s, content:%s", viper.ConfigFileUsed(), rule)
	} else {
		log.Fatal(err)
	}
	return rule, err
}

func init() {
	regexCmd.AddCommand(replaceCmd)
}
