package cmd

import (
	"strings"
)

/*
 * is skip check
 */
func (rs *RegexFactory) SplitMatch() {
	cpos := 0
	epos := len(rs.content)
	rs.Ranges = rs.Ranges[:cap(rs.Ranges)]
	//ranges := make([]RegexRange, 0)
	for _, match := range rs.Matches {
		if cpos < epos && cpos < match.Position.Start {
			//append string before match
			h := &RegexRange{
				Value:   rs.content[cpos:match.Position.Start],
				IsMatch: false,
			}
			rs.Ranges = append(rs.Ranges, *h)
		}
		// append match.value
		m := &RegexRange{
			Value:      rs.content[match.Position.Start:match.Position.End],
			IsMatch:    true,
			MatchIndex: match.Index,
		}
		rs.Ranges = append(rs.Ranges, *m)
		cpos = match.Position.End
	}
	if cpos < epos {
		//append last string
		f := &RegexRange{
			Value:   rs.content[cpos:epos],
			IsMatch: false,
		}
		rs.Ranges = append(rs.Ranges, *f)
	}
}

func (rs *RegexFactory) Restore() {
	config.Restore(&rs.content)
	//match restore
	for i, m := range rs.Matches {
		config.Restore(&rs.Matches[i].Value)
		for x := 0; x < len(m.Groups); x++ {
			config.Restore(&rs.Matches[i].Groups[x].Value)
		}
	}

	//range restore
	for x := 0; x < len(rs.Ranges); x++ {
		config.Restore(&rs.Ranges[x].Value)
	}
}

func (rs *RegexFactory) ToString() string {
	// restore before tostring.
	rs.Restore()

	cs := make([]string, 0)
	for _, m := range rs.Ranges {
		cs = append(cs, m.Value)
	}
	return strings.Join(cs, "")
}
