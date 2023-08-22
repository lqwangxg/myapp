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
	//log.Printf("%s Destfile: %s", flags.Action, filePath)
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
	//log.Printf("%s Destfile: %s", flags.Action, filePath)
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

func (rf *RegexFile) Execute() {
	content, err := ReadAll(rf.FromFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	hand := NewRegexTextByParent(rf, content)
	if isMatched, result := hand.Match(); isMatched {
		if flags.Action == "replace" {
			rf.Write(result)
		} else {
			hand.Write(result)
		}
	}
}

// write content to file
func (rs *RegexFile) Write(result *RegexResult) {
	if result == nil || result.MatchCount == 0 {
		return
	}

	if rs.ToFile == "" && rs.FromFile != "" {
		rs.ToFile = rs.FromFile
	}
	//TODO: add check-rule
	content, changed := result.Export(rs.Rule.ReplaceTemplate, false)
	if !changed {
		log.Print("No changed.")
	} else {
		config.Decode(&content)
		WriteAll(rs.ToFile, content)
		log.Printf("%s OK", flags.Action)
	}
}
