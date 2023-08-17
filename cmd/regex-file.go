package cmd

import (
	"log"
)

type RegexFile struct {
	Rule     *RegexRule
	FromFile string
	ToFile   string
}

func NewRegexFile(ruleName string, filePath string) *RegexFile {
	log.Printf("Regex %s file: %s", flags.Action, filePath)
	rule := appContext.RegexRules.GetRule(ruleName)
	if rule == nil {
		log.Printf("Regex rule not found, Over :<. Rule name: %s", ruleName)
		return nil
	}
	return &RegexFile{
		Rule:     rule,
		FromFile: filePath,
	}
}
func NewRegexFileByPattern(pattern, ruleName string, filePath string) *RegexFile {
	log.Printf("Regex %s file: %s", flags.Action, filePath)
	if ruleName == "" {
		ruleName = "default"
	}
	rule := appContext.RegexRules.GetRule(ruleName)
	if rule == nil {
		log.Printf("Regex rule not found, Over :<. Rule name: %s", ruleName)
		return nil
	}
	rule.Pattern = pattern
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
	result := reger.Match(hand)
	if flags.Action != "replace" {
		reger.Write(result, hand)
		return nil
	} else {
		return result
	}
}

// write content to file
func (rs *RegexFile) Write(result *RegexResult) {
	if result == nil {
		log.Print("No Result.")
		return
	}
	log.Printf("Match.Count=%d", result.MatchCount)
	if result.MatchCount == 0 {
		return
	}

	if rs.ToFile == "" && rs.FromFile != "" {
		rs.ToFile = rs.FromFile
	}
	content, changed := result.Export(&rs.Rule.ReplaceTemplate, false)
	if !changed {
		log.Print("No changed.")
	} else {
		config.Decode(&content)
		WriteAll(rs.ToFile, content)
		log.Printf("%s OK", flags.Action)
	}
}
