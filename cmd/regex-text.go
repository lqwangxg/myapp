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
}

func NewRegexText(pattern, content string) *RegexText {
	rs := &RegexText{pattern, content}
	return rs
}

func (rs *RegexText) Match() *RegexResult {
	// before match
	input := rs.Content
	config.Encode(&input)
	config.EncodePattern(&rs.Pattern)

	r := regexp.MustCompile(rs.Pattern)
	positions := r.FindAllStringSubmatchIndex(input, -1)
	result := &RegexResult{
		Pattern:    rs.Pattern,
		GroupNames: r.SubexpNames(),
		Captures:   make([]Capture, 0),
		Params:     make(map[string]string),
	}
	result.Captures = *SplitBy(&positions, input, false, result.Captures)
	result.Params["groups.count"] = strconv.Itoa(len(result.GroupNames))
	result.Params["groups.keys"] = strings.Join(result.GroupNames, ",")

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
		position := positions[x]

		for x := 0; x < len(result.GroupNames); x++ {
			gname := result.GroupNames[x]
			if x == 0 {
				gname = "match.value"
			}
			group := &Capture{Start: position[x*2+0], End: position[x*2+1]}
			group.Value = input[group.Start:group.End]
			if group.Params == nil {
				group.Params = make(map[string]string)
			}
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
		log.Print("No RegexResult, call Match firstly.")
		return
	}
	log.Printf("Match.Count=%d", result.MatchCount)
	if result.MatchCount == 0 {
		return
	}

	rule := appRules.GetDefaultRule()
	content := result.Export(&rule.ExportTemplate, false)
	config.Decode(&content)
	log.Println(content)
}
