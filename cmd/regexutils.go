package cmd

import (
	"log"
	"regexp"

	"golang.org/x/exp/slices"
)

// var regex regexp.Regexp
func CreateRegexFactory(pattern string, content string) *RegexFactory {
	r := regexp.MustCompile(pattern)
	//c := content
	//config.Transform(&content)
	return &RegexFactory{
		Pattern:    pattern,
		GroupNames: r.SubexpNames(),
		r:          r,
		content:    config.Transform(&content),
	}
}

// var regex regexp.Regexp
func (result *RegexFactory) ExecuteMatches() *RegexFactory {
	r := result.r
	sMatches := r.FindAllStringSubmatch(result.content, -1)
	positions := r.FindAllStringSubmatchIndex(result.content, -1)

	for i, smatch := range sMatches {
		position := positions[i]
		groups := make([]RegexGroup, 0)
		for x, groupName := range result.GroupNames {
			//from match 0 ~ count - 1
			groupIndex := r.SubexpIndex(groupName)
			if groupName == "" {
				groupIndex = x
			}

			group := RegexGroup{
				Index: x,
				Name:  groupName,
				Value: smatch[groupIndex],
				Position: RegexMatchIndex{
					Start: position[x*2+0],
					End:   position[x*2+1],
				},
			}
			groups = append(groups, group)
		}

		match := RegexMatch{
			Index: i,
			Position: RegexMatchIndex{
				Start: position[0],
				End:   position[1],
			},
			Groups: groups,
			Value:  smatch[0],
		}
		result.Matches = append(result.Matches, match)
	}
	return result
}

// func (matches *RegexResult) ReplaceMatches(src string) (string, bool) {
// 	r := matches.r
// 	new_src := r.ReplaceAllStringFunc(src, callback)
// 	return new_src, new_src != src
// }

// func callback(matchValue string) string {

// 	return matchValue
// }

func (matches *RegexFactory) GetMatches(index int) (bool, *RegexMatch) {
	idx := slices.IndexFunc(matches.Matches, func(m RegexMatch) bool { return m.Index == index })
	item := &RegexMatch{}
	if idx != -1 {
		return true, &matches.Matches[idx]
	}
	return false, item
}

func (matches *RegexFactory) log() {
	log.Printf("pattern:%s, group.count:%d, group.names:%v",
		matches.Pattern, len(matches.GroupNames), matches.GroupNames)
	for _, match := range matches.Matches {
		match.log()
	}
	for _, rg := range matches.Ranges {
		rg.log()
	}
}

func (match *RegexMatch) log() {
	log.Printf("match[%d].pos(%d,%d), group.count:%d, match.value=%s", match.Index,
		match.Position.Start, match.Position.End, len(match.Groups), match.Value)
	for _, group := range match.Groups {
		group.log()
	}
}

func (group *RegexGroup) log() {
	log.Printf("\tgroup[%d].pos(%d,%d), group.name=%s, group.value=%s", group.Index,
		group.Position.Start, group.Position.End, group.Name, group.Value)
}
func (rg *RegexRange) log() {
	log.Printf("\r range.isMatch=%v, MatchIndex:[%d], range.value=%s", rg.IsMatch, rg.MatchIndex, rg.Value)
}
