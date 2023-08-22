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
func (rf *RegexDirectory) Execute() {
	log.Printf("DirPath: %s", rf.DirPath)
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
				go Exec(NewRegexDirectory(flags.RuleName, fullPath))
			} else {
				go Exec(NewRegexFile(flags.RuleName, fullPath))
			}
		}
	}
}
