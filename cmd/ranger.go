package cmd

import "regexp"

// // split input by matches
// func (rs *Regex) SplitMatches(input string) {
// 	cpos := 0
// 	epos := len(input)
// 	rs.Result.Ranges = rs.Result.Ranges[:cap(rs.Result.Ranges)]
// 	for i, match := range rs.Result.Matches {
// 		if cpos < epos && cpos < match.Start {
// 			//append string before match
// 			h := &RegexRange{
// 				Capture: &Capture{
// 					Start: cpos,
// 					End:   match.Start,
// 					Value: input[cpos:match.Start],
// 					RType: UnMatchType,
// 				},
// 			}
// 			rs.Result.Ranges = append(rs.Result.Ranges, *h)
// 		}
// 		// append match.value
// 		m := &RegexRange{
// 			Capture: &Capture{
// 				Start: match.Start,
// 				End:   match.End,
// 				Value: input[match.Start:match.End],
// 				RType: MatchType,
// 			},
// 			MatchIndex: i,
// 		}
// 		rs.Result.Ranges = append(rs.Result.Ranges, *m)
// 		cpos = match.End
// 	}
// 	if cpos < epos {
// 		//append last string
// 		f := &RegexRange{
// 			Capture: &Capture{
// 				Start: cpos,
// 				End:   epos,
// 				Value: input[cpos:epos],
// 				RType: UnMatchType,
// 			},
// 		}
// 		rs.Result.Ranges = append(rs.Result.Ranges, *f)
// 	}
// }

// split input by pattern matches
// if matchOnly=true, will get matches Capture[Start:End] only.
func SplitMatchIndex(pattern, input string, matchOnly bool) *[]Capture {
	// return 0 ~ length when pattern is empty
	captures := make([]Capture, 0)
	if pattern == "" {
		c := &Capture{Start: 0, End: len(input)}
		c.Value = input[c.Start:c.End]
		captures = append(captures, *c)
		return &captures
	}
	config.Encode(&pattern)
	rs := regexp.MustCompile(pattern)
	positions := rs.FindAllStringSubmatchIndex(input, -1)
	return SplitBy(&positions, input, matchOnly, captures)
}

// split positions
// if matchOnly=true, will get matches Capture[Start:End] only.
func SplitBy(positions *[][]int, input string, matchOnly bool, captures []Capture) *[]Capture {
	//captures := make([]Capture, 0)
	if captures == nil {
		captures = make([]Capture, 0)
	}
	if len(*positions) == 0 {
		captures = append(captures, Capture{Start: 0, End: len(input)})
	} else {
		cpos := 0
		epos := len(input)
		for _, pos := range *positions {
			// match Capture
			match := Capture{Start: pos[0], End: pos[1], IsMatch: true}
			// append the ahead of match
			if !matchOnly && cpos < epos && cpos < match.Start {
				captures = append(captures, Capture{Start: cpos, End: match.Start})
			}
			// append match.value
			captures = append(captures, match)
			cpos = match.End
		}
		// append last string
		if !matchOnly && cpos < epos {
			captures = append(captures, Capture{Start: cpos, End: epos})
		}
	}
	for i := 0; i < len(captures); i++ {
		c := &captures[i]
		c.Value = input[c.Start:c.End]
	}
	return &captures
}
func (rule *RuleConfig) MergeRangeStartEnd(input string) *[]Capture {
	sBounds := SplitMatchIndex(rule.RangeStart, input, true)
	eBounds := SplitMatchIndex(rule.RangeEnd, input, true)

	rBounds := make([]Capture, 0)
	if len(*sBounds) > 1 {
		//==================================
		// on matches rangeStart, len(*sBound)>1
		sb := *sBounds
		for i := 0; i < len(*sBounds); i++ {
			s := sb[i]
			c := &Capture{Start: s.Start, End: len(input)}
			if i < len(*sBounds)-1 {
				c.End = sb[i+1].Start
			}
			for _, e := range *eBounds {
				if s.Start <= e.Start && e.End <= c.End {
					c.End = e.End
					break
				}
			}
			c.Value = input[c.Start:c.End]
			rBounds = append(rBounds, *c)
		}
		//last capture

	} else if len(*eBounds) > 1 {
		//==================================
		// on matches rangeEnd, len(*eBound)=1
		eb := *eBounds
		for i := 0; i < len(*eBounds); i++ {
			c := &Capture{Start: 0, End: eb[i].End}
			if i > 1 {
				c.Start = eb[i-1].End
				c.End = eb[i].End
			}
			c.Value = input[c.Start:c.End]
			rBounds = append(rBounds, *c)
		}
	} else {
		//rangeStart,rangeEnd
		rBounds = *sBounds
	}

	return &rBounds
}
