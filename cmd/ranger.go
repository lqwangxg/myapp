package cmd

// split input by pattern matches
// if matchOnly=true, will get matches Capture[Start:End] only.
func GetMatchCaptures(pattern, input string, matchOnly bool) (bool, *RegexResult) {
	// return 0 ~ length when pattern is empty
	result := &RegexResult{
		Pattern:   pattern,
		Positions: make([][]int, 0),
		Captures:  make([]Capture, 0),
	}
	// return full text to capture.
	if pattern == "" {
		c := &Capture{Start: 0, End: len(input)}
		c.Value = input[c.Start:c.End]
		result.Captures = append(result.Captures, *c)
		return false, result
	}

	config.Encode(&pattern)
	config.EncodePattern(&pattern)

	r := NewRegexText(pattern, input)
	isMatched, rsMatched := r.GetMatchResult(matchOnly, false)
	if isMatched {
		return true, rsMatched
	} else {
		return false, result
	}
}

func (rule *RegexRule) MergeRangeStartEnd(input string) *[]Capture {
	if rule.RangeStart == "" && rule.RangeEnd == "" {
		return nil
	}
	startMatched, rsStart := GetMatchCaptures(rule.RangeStart, input, true)
	endMatched, rsEnd := GetMatchCaptures(rule.RangeEnd, input, true)
	rBounds := make([]Capture, 0)
	if !startMatched && !endMatched {
		return &rBounds
	}
	sBounds := &rsStart.Captures
	eBounds := &rsEnd.Captures

	if startMatched {
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

	} else if endMatched {
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
