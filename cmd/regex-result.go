package cmd

import (
	"strconv"
	"strings"
)

type RegexResult struct {
	Pattern       string
	GroupNames    []string
	Captures      []Capture
	RangeCaptures []Capture
	Params        map[string]string
	MatchCount    int
	Positions     [][]int
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

		match.Params["match_index"] = strconv.Itoa(x)
		match.Params["match_start"] = strconv.Itoa(match.Start)
		match.Params["match_end"] = strconv.Itoa(match.End)

		for y := 0; y < len(rs.GroupNames); y++ {
			group := &Capture{Start: position[y*2+0], End: position[y*2+1], IsMatch: true}
			group.Value = input[group.Start:group.End]
			if group.Params == nil {
				group.Params = make(map[string]string)
			}
			gname := rs.GroupNames[y]
			if y == 0 {
				gname = "match_value"
			}

			group.Params["group_index"] = strconv.Itoa(y)
			group.Params["group_start"] = strconv.Itoa(group.Start)
			group.Params["group_end"] = strconv.Itoa(group.End)
			group.Params["group_key"] = gname
			group.Params["group_value"] = group.Value

			match.Params[gname] = group.Value
			match.Groups = append(match.Groups, *group)
		}
		x++
	}
	rs.MatchCount = x
	rs.Params["group_count"] = strconv.Itoa(x)
}
func (rs *RegexResult) MergeParams(fromMatch *Capture) {
	if !fromMatch.IsMatch {
		return
	}
	for key, fromValue := range fromMatch.Params {
		// If the key not exists, add key/value
		// if exists, do nothing
		if _, ok := rs.Params[key]; !ok {
			rs.Params[key] = fromValue
		}
	}
}
func (rs *RegexResult) Export(template *RegexTemplate, matchOnly bool) (string, bool) {
	var sb strings.Builder
	hasChanged := false
	if template == nil {
		return "", false
	}
	//---------------------------------------
	// replace Footer: rs.Result.Params
	tmp := template.Header
	tmp = NewTemplate(tmp).append(rs.Params)
	sb.WriteString(tmp)
	//---------------------------------------
	for _, item := range rs.Captures {
		if !item.Skip() {
			if template.Match != "" {
				tMatch := NewTemplate(template.Match)
				tmp, changed := tMatch.ReplaceBy(item)
				hasChanged = hasChanged || changed
				sb.WriteString(tmp)

				if template.Group != "" {
					tmp = template.GroupHeader
					tmp = NewTemplate(tmp).append(rs.Params)
					tmp = NewTemplate(tmp).append(item.Params)
					sb.WriteString(tmp)
					//replace by template.group
					for _, g := range item.Groups {
						tGroup := NewTemplate(template.Group)
						tmp, changed := tGroup.ReplaceBy(g)
						sb.WriteString(tmp)
						if changed {
							hasChanged = true
						}
					}
					tmp = template.GroupFooter
					tmp = NewTemplate(tmp).append(rs.Params)
					tmp = NewTemplate(tmp).append(item.Params)
					sb.WriteString(tmp)
				}
				if template.Param != "" {
					//replace by template.param
					tmp = template.ParamHeader
					tmp = NewTemplate(tmp).append(rs.Params)
					tmp = NewTemplate(tmp).append(item.Params)
					sb.WriteString(tmp)

					x := 0
					for key, val := range item.Params {
						tmp = template.Param
						tmp, keyChanged := NewTemplate(tmp).ReplaceByKeyValue("param_key", key)
						tmp, valChanged := NewTemplate(tmp).ReplaceByKeyValue("param_value", val)
						tmp, idxChanged := NewTemplate(tmp).ReplaceByKeyValue("param_index", strconv.Itoa(x))
						sb.WriteString(tmp)
						if keyChanged || valChanged || idxChanged {
							hasChanged = true
						}
						x++
					}

					tmp = template.ParamFooter
					tmp = NewTemplate(tmp).append(rs.Params)
					tmp = NewTemplate(tmp).append(item.Params)
					sb.WriteString(tmp)
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
	tmp = template.Footer
	tmp = NewTemplate(tmp).append(rs.Params)
	sb.WriteString(tmp)
	//---------------------------------------
	return sb.String(), hasChanged
}
func (rs *RegexResult) RefreshFromRanges() {
	for i := 0; i < len(rs.Captures); i++ {
		if rs.Captures[i].IsMatch {
			rs.Contained(&rs.Captures[i])
		}
	}
}
func (rs *RegexResult) Contained(match *Capture) {
	if rs.RangeCaptures == nil {
		match.Params["match_inrange"] = strconv.FormatBool(true)
		return
	}

	for _, r := range rs.RangeCaptures {
		if r.Start <= match.Start && match.End <= r.End {
			match.Params["match_inrange"] = strconv.FormatBool(true)
			match.Params["range_value"] = r.Value
			match.MergeParams(&r)
			return
		}
	}
	match.Params["match_inrange"] = strconv.FormatBool(false)
}
