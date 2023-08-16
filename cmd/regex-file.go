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
	rule := appRules.GetRule(ruleName)
	if rule == nil {
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
		log.Print("No RegexResult, call Match firstly.")
		return
	}
	if rs.ToFile == "" && rs.FromFile != "" {
		rs.ToFile = rs.FromFile
	}
	content := result.Export(&rs.Rule.ReplaceTemplate, false)
	log.Printf("Writing To: %s", rs.ToFile)
	WriteAll(rs.ToFile, content)
	log.Printf("Written Completed: %s", rs.ToFile)
}
