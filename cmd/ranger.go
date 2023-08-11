package cmd

import (
	"strings"
)

func (rs *RegexFactory) Close() {
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
	//save to cache
	rs.ToCache()
}

func (rs *RegexFactory) ToString() string {
	// restore before tostring.
	rs.Close()

	cs := make([]string, 0)
	for _, m := range rs.Ranges {
		cs = append(cs, m.Value)
	}
	return strings.Join(cs, "")
}
