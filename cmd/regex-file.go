package cmd

import (
	"log"
)

type RegexFile struct {
	Rule     *RuleConfig
	Result   *RegexResult
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

func (rf *RegexFile) Match() {
	content, err := ReadAll(rf.FromFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	hand := NewRegexText(rf.Rule.Pattern, content)
	reger.Execute(hand)
}

// write content to file
func (rs *RegexFile) Write() {
	if rs.Result == nil {
		log.Print("No RegexResult, call Match firstly.")
		return
	}
	if rs.ToFile == "" && rs.FromFile != "" {
		rs.ToFile = rs.FromFile
	}
	content := rs.Result.Export(&rs.Rule.ReplaceTemplate, false)
	log.Printf("Writing To: %s", rs.ToFile)
	WriteAll(rs.ToFile, content)
	log.Printf("Written Completed: %s", rs.ToFile)
}
