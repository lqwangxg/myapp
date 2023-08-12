package cmd

import (
	"strings"
)

func (rs *Regex) Close() {
	//match restore
	for i, m := range rs.Result.Matches {
		config.Restore(&rs.Result.Matches[i].Value)
		for x := 0; x < len(m.Groups); x++ {
			config.Restore(&rs.Result.Matches[i].Groups[x].Value)
		}
	}
	//range restore
	for x := 0; x < len(rs.Result.Ranges); x++ {
		config.Restore(&rs.Result.Ranges[x].Value)
	}
	//save to cache
	rs.ToCache()
}

func (rs *Regex) ExportMatches(template string) string {
	var sb strings.Builder
	for i := 0; i < len(rs.Result.Matches); i++ {
		ReplaceTemplate(&template, rs.Result.Matches[i].Params)
	}
	sb.WriteString(template)
	exports := sb.String()
	config.Restore(&exports)
	return exports
}
