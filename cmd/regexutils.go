package cmd

import (
	"log"
	"regexp"
)

func Matches(pattern, content string) *RegexMatchResult {

	r := regexp.MustCompile(pattern)
	var matchResult = &RegexMatchResult{
		Pattern:    pattern,
		GroupNames: r.SubexpNames(),
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

func (matches *RegexMatchResult) restore() {
	for _, match := range matches.Matches {
		match.Value = afterMatch(match.Value)
		for _, group := range match.Groups {
			group.Value = afterMatch(group.Value)
		}
	}
}

func (matches *RegexMatchResult) Log() {
	log.Printf("pattern:%s, group.count:%d, group.names:%v",
		matches.Pattern, len(matches.GroupNames), matches.GroupNames)
	for _, match := range matches.Matches {
		match.Log()
	}
}

func (match *RegexMatch) Log() {
	log.Printf("match[%d].pos(%d,%d), group.count:%d, match.value=%s", match.Index,
		match.Position.Start, match.Position.End, len(match.Groups), match.Value)
	for _, group := range match.Groups {
		group.Log()
	}
}

func (group *RegexGroup) Log() {
	log.Printf("\tgroup[%d].pos(%d,%d), group.name=%s, group.value=%s", group.Index,
		group.Position.Start, group.Position.End, group.Name, group.Value)
}
