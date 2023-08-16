package cmd

import (
	"log"
)

type RegexFile struct {
	Rule     *RuleConfig
	FromFile string
	ToFile   string
}

func NewRegexFile(ruleName string, filePath string) *RegexFile {
	log.Printf("Regex %s file: %s", flags.Action, filePath)
	rule := appRules.GetRule(ruleName)
	if rule == nil {
		log.Printf("Regex rule not found, Over :<. Rule name: %s", ruleName)
		return nil
	}
	return &RegexFile{
		Rule:     rule,
		FromFile: filePath,
	}
}

func (rf *RegexFile) Match() *RegexResult {
	content, err := ReadAll(rf.FromFile)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	hand := NewRegexText(rf.Rule.Pattern, content)
	return reger.Match(hand)
}

// write content to file
func (rs *RegexFile) Write(result *RegexResult) {
	if result == nil {
		log.Print("No Regex Result, No Write, Over :<.")
		return
	}
	if result.MatchCount == 0 {
		log.Print("No Regex Result, No Write, Over :<.")
		return
	}
	if rs.ToFile == "" && rs.FromFile != "" {
		rs.ToFile = rs.FromFile
	}
	content, changed := result.Export(&rs.Rule.ReplaceTemplate, false)
	if !changed {
		log.Print("No Changed.")
	} else {
		config.Decode(&content)
		WriteAll(rs.ToFile, content)
		log.Printf("%s OK, Written file: %s", flags.Action, rs.ToFile)
	}
}
