package cmd

import "regexp"

// split input by matches
func SplitMatchIndex(pattern, input string) *[]Bound {
	bounds := make([]Bound, 0)
	// return 0 ~ length when pattern is empty
	if pattern == "" {
		bounds = append(bounds, Bound{Start: 0, End: len(input)})
		return &bounds
	}

	r := regexp.MustCompile(pattern)
	positions := r.FindAllStringSubmatchIndex(input, -1)
	if len(positions) == 0 {
		bounds = append(bounds, Bound{Start: 0, End: len(input)})
	} else {
		for _, pos := range positions {
			// append match.value
			bounds = append(bounds, Bound{Start: pos[0], End: pos[1]})
		}
	}
	return &bounds
}
func MergeRangeStartEnd(rangeStart, rangeEnd, input string) *[]Bound {
	sBound := SplitMatchIndex(rangeStart, input)
	eBound := SplitMatchIndex(rangeEnd, input)

	rBound := make([]Bound, 0)
	if len(*sBound) > 1 {
		//==================================
		// on matches rangeStart, len(*sBound)>1
		sb := *sBound
		for i := 0; i < len(*sBound)-1; i++ {
			s := sb[i]
			n := sb[i+1]
			r := &Bound{Start: s.Start, End: n.Start}
			for _, e := range *eBound {
				if s.Start <= e.Start && e.End <= n.Start {
					r.End = e.End
					break
				}
			}
			rBound = append(rBound, *r)
		}
	} else if len(*eBound) > 1 {
		//==================================
		// on matches rangeEnd, len(*eBound)=1
		eb := *eBound
		for i := 1; i < len(*eBound); i++ {
			eh := eb[i-1]
			ec := eb[i]
			if i == 1 {
				rBound = append(rBound, Bound{Start: 0, End: eh.End})
			} else {
				rBound = append(rBound, Bound{Start: eh.End, End: ec.End})
			}
		}
	} else {
		//rangeStart,rangeEnd
		rBound = *sBound
	}

	return &rBound
}

// func SplitRangeStartEnd(mPattern, rangeStart, rangeEnd, input string) {
// 	mRange := SplitMatches(mPattern, input)
// 	sRange := SplitMatches(rangeStart, input)
// 	eRange := SplitMatches(rangeEnd, input)
// 	rRange := make([]RegexRange, 0)
// 	for _, m := range *mRange {
// 		r := &RegexRange{
// 			Bound: Bound{Start: 0, End: 0},
// 		}
// 		//compare(startRange,matchRange)
// 		for _, s := range *sRange {
// 			if s.Start <= m.Start {
// 				//reset Start on startRange.Start <= matchRange.Start
// 				r.Start = s.Start
// 			} else {
// 				//break on startRange.Start > matchRange.Start
// 				break
// 			}
// 		}
// 		//compare(matchRange, endRange)
// 		for _, e := range *eRange {
// 			if m.End <= e.End {
// 				//break on matchRange.End <= endRange.End
// 				r.End = e.End
// 				break
// 			}
// 		}

// 	}

// }
