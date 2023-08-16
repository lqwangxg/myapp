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
	Result  *RegexResult
}

func NewRegexText(pattern, content string) *RegexText {
	rs := &RegexText{pattern, content, nil}
	return rs
}

func (rs *RegexText) Match() {
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
		MatchCount: len(positions) / 2,
	}
	result.Captures = *SplitBy(&positions, input, false, result.Captures)
	result.Params["matches.count"] = strconv.Itoa(result.MatchCount)
	result.Params["groups.count"] = strconv.Itoa(len(result.GroupNames))
	result.Params["groups.keys"] = strings.Join(result.GroupNames, ",")
	rs.Result = result
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
			group.SetValue(input)
			match.Groups = append(match.Groups, *group)
			match.Params[gname] = group.Value
		}
		x++
	}
}

// write content to file
func (rs *RegexText) Write() {
	if rs.Result == nil {
		log.Print("No RegexResult, call Match firstly.")
		return
	}
	log.Printf("Match.Count=%d", rs.Result.MatchCount)
	if rs.Result.MatchCount == 0 {
		return
	}

	rule := appRules.GetDefaultRule()
	content := rs.Result.Export(&rule.ExportTemplate, false)
	log.Println(content)
}
