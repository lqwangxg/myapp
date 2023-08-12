package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"regexp"
)

// var regex regexp.Regexp
func NewRegex(pattern string) *Regex {
	return NewCacheRegex(pattern, config.RedisOption.Enable)
}

// var regex regexp.Regexp
func NewCacheRegex(pattern string, cache bool) *Regex {
	r := regexp.MustCompile(pattern)
	return &Regex{
		R: r,
		Result: RegexResult{
			Pattern:    pattern,
			GroupNames: r.SubexpNames(),
			Params:     make(map[string]string),
		},
		Cache: cache,
	}
}

// get hashsum of (pattern + ": "+ input) to string.
func (rs *Regex) Hashsum(input string) string {
	h := sha1.New()
	h.Write([]byte(rs.Result.Pattern + ": " + input))
	return hex.EncodeToString(h.Sum(nil))
}
func (rs *Regex) IsMatch(input string) bool {
	return rs.R.MatchString(input)
}

// var regex regexp.Regexp
func (rs *Regex) ScanMatches(input string) {
	// get from redis cache
	if rs.FromCache(input) {
		return
	}

	//before match input, transfer special chars
	config.Transform(&input)
	r := rs.R
	if !r.MatchString(input) {
		return
	}
	subMatches := r.FindAllStringSubmatch(input, -1)
	positions := r.FindAllStringSubmatchIndex(input, -1)

	for i, smatch := range subMatches {
		position := positions[i]
		groups := make([]RegexGroup, 0)
		match := RegexMatch{
			Index: i,
			Position: RegexMatchIndex{
				Start: position[0],
				End:   position[1],
			},
			Groups: groups,
			Value:  smatch[0],
			Params: make(map[string]string),
		}
		for x, groupName := range rs.Result.GroupNames {
			//from match 0 ~ count - 1
			groupIndex := r.SubexpIndex(groupName)
			if groupName == "" {
				groupIndex = x
			}
			gname := groupName
			if x == 0 {
				gname = "match.value"
			}
			group := RegexGroup{
				Index: x,
				Name:  gname,
				Value: smatch[groupIndex],
				Position: RegexMatchIndex{
					Start: position[x*2+0],
					End:   position[x*2+1],
				},
			}
			match.Params[gname] = group.Value
			groups = append(groups, group)
		}

		rs.Result.Matches = append(rs.Result.Matches, match)
	}
	// split input by matches
	rs.SplitMatches(input)
}

// split input by matches
func (rs *Regex) SplitMatches(input string) []RegexRange {
	cpos := 0
	epos := len(input)
	rs.Result.Ranges = rs.Result.Ranges[:cap(rs.Result.Ranges)]
	for _, match := range rs.Result.Matches {
		if cpos < epos && cpos < match.Position.Start {
			//append string before match
			h := &RegexRange{
				Value:   input[cpos:match.Position.Start],
				IsMatch: false,
			}
			rs.Result.Ranges = append(rs.Result.Ranges, *h)
		}
		// append match.value
		m := &RegexRange{
			Value:      input[match.Position.Start:match.Position.End],
			IsMatch:    true,
			MatchIndex: match.Index,
		}
		rs.Result.Ranges = append(rs.Result.Ranges, *m)
		cpos = match.Position.End
	}
	if cpos < epos {
		//append last string
		f := &RegexRange{
			Value:   input[cpos:epos],
			IsMatch: false,
		}
		rs.Result.Ranges = append(rs.Result.Ranges, *f)
	}
	return rs.Result.Ranges
}

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
