package cmd

import (
	"log"
	"os"
	"path/filepath"
)

type RegexDirectory struct {
	Rule      *RuleConfig
	Result    *RegexResult
	DirPath   string
	ToDirPath string
}

func NewRegexDirectory(ruleName, dirPath string) *RegexDirectory {
	rule := appRules.GetRule(ruleName)
	if rule == nil {
		return nil
	}
	return &RegexDirectory{
		Rule:    rule,
		DirPath: dirPath,
	}
}
func (rf *RegexDirectory) Match() {
	files, err := os.ReadDir(rf.DirPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, file := range files {
		fullPath := filepath.Join(rf.DirPath, file.Name())
		ok, err := IsDir(fullPath)
		if err == nil {
			if ok {
				handler := NewRegexDirectory(flags.RuleName, fullPath)
				reger.Execute(handler)
			} else {
				handler := NewRegexFile(flags.RuleName, fullPath)
				reger.Execute(handler)
			}
		}
	}
}

// write content to file
func (rs *RegexDirectory) Write() {
	log.Printf("Written Completed: %s", rs.DirPath)
}
