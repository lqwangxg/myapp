package cmd

import (
	"log"
	"strconv"
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

	result := rs.GetMatchResult()
	// match.Index.
	x := 0
	for i, c := range result.Captures {
		// skip if it's not match
		if !c.IsMatch {
			continue
		}

		match := &result.Captures[i]
		match.Groups = make([]Capture, 0)
		match.Params = make(map[string]string)
		position := result.Positions[x]
		match.Params["index"] = strconv.Itoa(x)
		match.Params["match.start"] = strconv.Itoa(match.Start)
		match.Params["match.end"] = strconv.Itoa(match.End)

		for y := 0; y < len(result.GroupNames); y++ {
			gname := result.GroupNames[y]
			if y == 0 {
				gname = "match.value"
			}
			group := &Capture{Start: position[y*2+0], End: position[y*2+1], IsMatch: true}
			group.Value = rs.Content[group.Start:group.End]
			if group.Params == nil {
				group.Params = make(map[string]string)
			}
			group.Params["index"] = strconv.Itoa(y)
			group.Params["group.start"] = strconv.Itoa(group.Start)
			group.Params["group.end"] = strconv.Itoa(group.End)
			group.Params["group.key"] = gname
			group.Params["group.value"] = group.Value
			match.Groups = append(match.Groups, *group)
			match.Params[gname] = group.Value
		}
		x++
	}
	result.MatchCount = x
	result.Params["matches.count"] = strconv.Itoa(x)
	return result
}

// write content to file
func (rs *RegexText) Write(result *RegexResult) {
	if result == nil {
		return
	}
	log.Printf("Match.Count=%d", result.MatchCount)
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
