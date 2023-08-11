package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"regexp"
)

// var regex regexp.Regexp
func CreateRegexFactory(pattern string) *RegexFactory {
	r := regexp.MustCompile(pattern)
	return &RegexFactory{
		Pattern:    pattern,
		GroupNames: r.SubexpNames(),
		regex:      r,
	}
}

// get hashsum of (pattern + ": "+ input) to string.
func (result *RegexFactory) Hashsum(input string) string {
	h := sha1.New()
	h.Write([]byte(result.Pattern + ": " + input))
	return hex.EncodeToString(h.Sum(nil))
}

// var regex regexp.Regexp
func (rs *RegexFactory) ExecuteMatches(input string) *RegexFactory {
	// reset inputkey from pattern+":"+input
	rs.InputKey = rs.Hashsum(input)
	// get from redis cache
	if rs.FromCache() {
		return rs
	}

	//before match input, transfer special chars
	config.Transform(&input)
	r := rs.regex
	sMatches := r.FindAllStringSubmatch(input, -1)
	positions := r.FindAllStringSubmatchIndex(input, -1)

	for i, smatch := range sMatches {
		position := positions[i]
		groups := make([]RegexGroup, 0)
		for x, groupName := range rs.GroupNames {
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
		rs.Matches = append(rs.Matches, match)
	}
	// split input by matches
	rs.splitByMatches(input)
	return rs
}

// split input by matches
func (rs *RegexFactory) splitByMatches(input string) {
	cpos := 0
	epos := len(input)
	rs.Ranges = rs.Ranges[:cap(rs.Ranges)]
	for _, match := range rs.Matches {
		if cpos < epos && cpos < match.Position.Start {
			//append string before match
			h := &RegexRange{
				Value:   input[cpos:match.Position.Start],
				IsMatch: false,
			}
			rs.Ranges = append(rs.Ranges, *h)
		}
		// append match.value
		m := &RegexRange{
			Value:      input[match.Position.Start:match.Position.End],
			IsMatch:    true,
			MatchIndex: match.Index,
		}
		rs.Ranges = append(rs.Ranges, *m)
		cpos = match.Position.End
	}
	if cpos < epos {
		//append last string
		f := &RegexRange{
			Value:   input[cpos:epos],
			IsMatch: false,
		}
		rs.Ranges = append(rs.Ranges, *f)
	}
}

// func (matches *RegexResult) ReplaceMatches(src string) (string, bool) {
// 	r := matches.r
// 	new_src := r.ReplaceAllStringFunc(src, callback)
// 	return new_src, new_src != src
// }

// func callback(matchValue string) string {

// 	return matchValue
// }

// func (matches *RegexFactory) GetMatches(index int) (bool, *RegexMatch) {
// 	idx := slices.IndexFunc(matches.Matches, func(m RegexMatch) bool { return m.Index == index })
// 	item := &RegexMatch{}
// 	if idx != -1 {
// 		return true, &matches.Matches[idx]
// 	}
// 	return false, item
// }

// func (matches *RegexFactory) log() {
// 	log.Printf("pattern:%s, group.count:%d, group.names:%v",
// 		matches.Pattern, len(matches.GroupNames), matches.GroupNames)
// 	for _, match := range matches.Matches {
// 		match.log()
// 	}
// 	for _, rg := range matches.Ranges {
// 		rg.log()
// 	}
// }

// func (match *RegexMatch) log() {
// 	log.Printf("match[%d].pos(%d,%d), group.count:%d, match.value=%s", match.Index,
// 		match.Position.Start, match.Position.End, len(match.Groups), match.Value)
// 	for _, group := range match.Groups {
// 		group.log()
// 	}
// }

// func (group *RegexGroup) log() {
// 	log.Printf("\tgroup[%d].pos(%d,%d), group.name=%s, group.value=%s", group.Index,
// 		group.Position.Start, group.Position.End, group.Name, group.Value)
// }
// func (rg *RegexRange) log() {
// 	log.Printf("\r range.isMatch=%v, MatchIndex:[%d], range.value=%s", rg.IsMatch, rg.MatchIndex, rg.Value)
// }
