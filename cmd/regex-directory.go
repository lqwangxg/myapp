package cmd

import (
	"log"
	"os"
	"path/filepath"
)

type RegexDirectory struct {
	Rule      *RegexRule
	Result    *RegexResult
	DirPath   string
	ToDirPath string
}

func NewRegexDirectory(ruleName, dirPath string) *RegexDirectory {
	rule := appContext.RegexRules.GetRule(ruleName)
	if rule == nil {
		log.Printf("RegexRule not found by ruleName: %s", ruleName)
		return nil
	}
	return &RegexDirectory{
		Rule:    rule,
		DirPath: dirPath,
	}
}
func (rf *RegexDirectory) Match() *RegexResult {
	log.Printf("DirPath: %s", rf.DirPath)
	files, err := os.ReadDir(rf.DirPath)
	if err != nil {
		log.Fatal(err)
		return nil
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
	return nil
}

// write content to file
func (rs *RegexDirectory) Write(result *RegexResult) {
	//log.Printf("Written Completed: %s", rs.DirPath)
}
