package cmd

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type RegexText struct {
	Pattern string
	Content string
	Parent  *RegexFile
	*RegexRule
}

func NewRegexText(pattern, content string) *RegexText {
	rs := &RegexText{pattern, content, nil, nil}
	return rs
}
func NewRegexTextByParent(parent *RegexFile, content string) *RegexText {
	var rule *RegexRule
	if parent != nil {
		rule = parent.Rule
	}
	rs := &RegexText{
		Pattern:   parent.Rule.Pattern,
		Content:   content,
		Parent:    parent,
		RegexRule: rule}
	return rs
}
func (rs *RegexText) refreshRule(ruleName string) {

	if ruleName == "" {
		ruleName = "default"
	}
	rule := appContext.RegexRules.GetRule(ruleName)
	if rule == nil {
		log.Printf("Regex rule not found, Over :<. Rule name: %s", ruleName)
		return
	}
	rs.RegexRule = rule
}
func (rs *RegexText) Match() *RegexResult {
	// before match
	config.Encode(&rs.Content)
	config.EncodePattern(&rs.Pattern)

	isMatched, result := rs.GetMatchResult(false)
	if !isMatched {
		return nil
	}
	return result
}

func (rs *RegexText) GetMatchResult(matchOnly bool) (bool, *RegexResult) {
	r := regexp.MustCompile(rs.Pattern)
	positions := r.FindAllStringSubmatchIndex(rs.Content, -1)
	if len(positions) == 0 {
		log.Printf("No Matched. pattern: %s", rs.Pattern)
		return false, nil
	}
	result := &RegexResult{
		Pattern:    rs.Pattern,
		GroupNames: r.SubexpNames(),
		Params:     make(map[string]string),
		Positions:  positions,
	}

	result.Params["pattern"] = rs.Pattern
	result.Params["matches.count"] = strconv.Itoa(0)
	result.Params["groups.count"] = strconv.Itoa(len(result.GroupNames))
	result.Params["groups.keys"] = strings.Join(result.GroupNames, ",")
	result.SplitBy(rs.Content, matchOnly)
	result.FillParams(rs.Content)
	return true, result
}

// write content to file
func (rs *RegexText) Write(result *RegexResult) {
	if result == nil {
		return
	}
	if result.MatchCount == 0 {
		return
	}

	var rule RegexRule
	if rs.Parent != nil {
		rule = *rs.Parent.Rule
	} else {
		rule = *appContext.RegexRules.GetDefaultRule()
	}
	rule.ResetTemplate()
	//TODO: add check-rule
	content, changed := result.Export(rule.ExportTemplate, true)
	if !changed {
		log.Print("No changed.")
	} else {
		config.Decode(&content)
		log.Println(content)
	}
}
func (rs *RegexResult) FirstMatch() *Capture {
	if rs == nil || rs.MatchCount == 0 {
		return nil
	}
	for _, c := range rs.Captures {
		if c.IsMatch {
			return &c
		}
	}
	return nil
}
