package cmd

import (
	"strconv"
	"strings"
)

type RegexResult struct {
	Pattern    string
	GroupNames []string
	Captures   []Capture
	Params     map[string]string
	MatchCount int
	Positions  [][]int
}

// split positions
// if matchOnly=true, will get matches Capture[Start:End] only.
func (rs *RegexResult) SplitBy(input string, matchOnly bool) *[]Capture {
	// no alias, because no update, reference only.
	positions := rs.Positions
	//alias result.Captures for updating itself
	captures := &rs.Captures
	if len(positions) == 0 {
		*captures = append(*captures, Capture{Start: 0, End: len(input)})
	} else {
		cpos := 0
		epos := len(input)
		for _, pos := range positions {
			// match Capture
			match := Capture{Start: pos[0], End: pos[1], IsMatch: true}
			// append the ahead of match
			if !matchOnly && cpos < epos && cpos < match.Start {
				*captures = append(*captures, Capture{Start: cpos, End: match.Start})
			}
			// append match.value
			*captures = append(*captures, match)
			cpos = match.End
		}
		// append last string
		if !matchOnly && cpos < epos {
			*captures = append(*captures, Capture{Start: cpos, End: epos})
		}
	}
	//refresh result.Captures.value
	for i := 0; i < len(rs.Captures); i++ {
		c := &rs.Captures[i]
		c.Value = input[c.Start:c.End]
	}
	return captures
}

func (rs *RegexResult) FillParams(input string) {

	// match.Index.
	x := 0
	for i, c := range rs.Captures {
		// skip if it's not match
		if !c.IsMatch {
			continue
		}

		match := &rs.Captures[i]
		match.Groups = make([]Capture, 0)
		match.Params = make(map[string]string)
		position := rs.Positions[x]
		match.Params["index"] = strconv.Itoa(x)
		match.Params["match.start"] = strconv.Itoa(match.Start)
		match.Params["match.end"] = strconv.Itoa(match.End)

		for y := 0; y < len(rs.GroupNames); y++ {
			gname := rs.GroupNames[y]
			if y == 0 {
				gname = "match.value"
			}
			group := &Capture{Start: position[y*2+0], End: position[y*2+1], IsMatch: true}
			group.Value = input[group.Start:group.End]
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
	rs.MatchCount = x
	rs.Params["matches.count"] = strconv.Itoa(x)
}

func (rs *RegexResult) Export(template *RegexTemplate, matchOnly bool) (string, bool) {
	var sb strings.Builder
	hasChanged := false
	//---------------------------------------
	// replace Footer: rs.Result.Params
	if template != nil && template.Header != "" {
		tHeader := NewTemplate(template.Header)
		header, changed := tHeader.ReplaceByMap(rs.Params)
		sb.WriteString(header)
		if changed {
			hasChanged = true
		}
	}

	//---------------------------------------
	for _, item := range rs.Captures {
		if item.IsMatch {
			if template != nil && template.Match != "" {
				tMatch := NewTemplate(template.Match)
				tmp, changed := tMatch.ReplaceBy(item)
				sb.WriteString(tmp)
				if changed {
					hasChanged = true
				}
				if template.Group != "" {
					//replace by template.group
					for _, g := range item.Groups {
						tGroup := NewTemplate(template.Group)
						tmp, changed := tGroup.ReplaceBy(g)
						sb.WriteString(tmp)
						if changed {
							hasChanged = true
						}
					}
				}
			} else {
				//when template is empty, export match.value
				sb.WriteString(item.Value)
			}
		} else if !matchOnly {
			// export not only match
			sb.WriteString(item.Value)
		}
	}
	//---------------------------------------
	if template != nil && template.Footer != "" {
		tFooter := NewTemplate(template.Footer)
		footer, changed := tFooter.ReplaceByMap(rs.Params)
		sb.WriteString(footer)
		if changed {
			hasChanged = true
		}
	}
	//---------------------------------------
	return sb.String(), hasChanged
}
