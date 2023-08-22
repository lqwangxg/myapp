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
func (rs *RegexText) Execute() {
	isMatched, result := rs.Match()
	if !isMatched {
		return
	}
	rs.Write(result)
}

func (rs *RegexText) Match() (bool, *RegexResult) {
	// before match
	config.Encode(&rs.Content)
	config.EncodePattern(&rs.Pattern)
	return rs.GetMatchResult(false)
}

func (rs *RegexText) GetMatchResult(matchOnly bool) (bool, *RegexResult) {
	r := regexp.MustCompile(rs.Pattern)
	positions := r.FindAllStringSubmatchIndex(rs.Content, -1)
	if len(positions) == 0 {
		if rs.Parent != nil {
			log.Printf("No Matched. file:%s", rs.Parent.FromFile)
		} else {
			log.Print("No Matched.")
		}
		return false, nil
	}
	result := &RegexResult{
		Pattern:    rs.Pattern,
		GroupNames: r.SubexpNames(),
		Params:     make(map[string]string),
		Positions:  positions,
	}
	if rs.Parent != nil {
		result.Params["from.file"] = rs.Parent.FromFile
		result.Params["to.file"] = rs.Parent.ToFile
	}
	if rs.RegexRule != nil {
		result.Params["rule.name"] = rs.RegexRule.Name
	}
	result.Params["pattern"] = rs.Pattern
	result.Params["matches.count"] = strconv.Itoa(0)
	result.Params["groups.count"] = strconv.Itoa(len(result.GroupNames))
	result.Params["groups.keys"] = strings.Join(result.GroupNames, ",")
	result.SplitBy(rs.Content, matchOnly)
	result.FillParams(rs.Content)
	rs.FillParams(result)
	return true, result
}

func (rs *RegexText) FillParams(result *RegexResult) {
	if rs.RegexRule == nil {
		return
	}
	// match full content
	for _, pattern := range rs.RegexRule.ParamPatterns.Fulls {
		tmpRT := NewRegexText(pattern, rs.Content)
		if tmpOK, tmpRS := tmpRT.GetMatchResult(true); tmpOK {
			for _, m := range tmpRS.Captures {
				result.MergeParams(&m)
			}
		}
	}

	//TODO: match range.value

	// match match.value
	for _, pattern := range rs.RegexRule.ParamPatterns.Matches {
		for x := 0; x < len(result.Captures); x++ {
			if !result.Captures[x].IsMatch {
				continue
			}
			tmpRT := NewRegexText(pattern, result.Captures[x].Value)
			if tmpOK, tmpRS := tmpRT.GetMatchResult(true); tmpOK {
				for _, tm := range tmpRS.Captures {
					result.Captures[x].MergeParams(&tm)
				}
			}
		}
	}
}
func (rs *RegexText) EvalFormulas() {

}

// write content to file
func (rs *RegexText) Write(result *RegexResult) {
	if result == nil || result.MatchCount == 0 {
		log.Printf("result:%v", result)
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
