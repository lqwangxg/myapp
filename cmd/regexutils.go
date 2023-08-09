package cmd

import (
	"log"
	"regexp"
)

// var regex regexp.Regexp
func Matches(pattern, content string) *RegexMatchResult {

	r := regexp.MustCompile(pattern)
	var matchResult = &RegexMatchResult{
		Pattern:    pattern,
		GroupNames: r.SubexpNames(),
		r:          r,
	}

	sMatches := r.FindAllStringSubmatch(content, -1)
	positions := r.FindAllStringSubmatchIndex(content, -1)

	for i, smatch := range sMatches {
		position := positions[i]
		groups := make([]RegexGroup, 0)
		for x, groupName := range matchResult.GroupNames {
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
		matchResult.Matches = append(matchResult.Matches, match)

	}
	return matchResult
}

func (matches *RegexMatchResult) ReplaceMatches(src string) (string, bool) {
	r := matches.r
	new_src := r.ReplaceAllStringFunc(src, callback)
	return new_src, new_src != src
}

func callback(matchValue string) string {

	return matchValue
}

func (matches *RegexMatchResult) restore() {
	for _, match := range matches.Matches {
		match.Value = afterMatch(match.Value)
		for _, group := range match.Groups {
			group.Value = afterMatch(group.Value)
		}
	}
}

func (matches *RegexMatchResult) log() {
	log.Printf("pattern:%s, group.count:%d, group.names:%v",
		matches.Pattern, len(matches.GroupNames), matches.GroupNames)
	for _, match := range matches.Matches {
		match.log()
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
