package cmd

import "strings"

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
	for _, m := range rs.Matches {
		m.Restore()
	}
	//range restore
	for _, m := range rs.Ranges {
		m.Restore()
	}
}
func (m *RegexMatch) Restore() {
	m.Value = config.Restore(&m.Value)
	for _, g := range m.Groups {
		g.Restore()
	}
}
func (m *RegexGroup) Restore() {
	m.Value = config.Restore(&m.Value)
}
func (m *RegexRange) Restore() {
	m.Value = config.Restore(&m.Value)
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
