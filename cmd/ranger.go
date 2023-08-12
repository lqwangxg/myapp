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

func (rs *Regex) ReplaceMatch(template string, index int) {
	rs.Result.Matches[index].Value = ReplaceTemplate(template, rs.Result.Matches[index].Params)
}

func (rs *Regex) ExportMatches(template string) string {
	cs := make([]string, 0)
	for i := 0; i < len(rs.Result.Matches); i++ {
		cs = append(cs, ReplaceTemplate(template, rs.Result.Matches[i].Params))
	}
	exports := strings.Join(cs, "")
	config.Restore(&exports)
	return exports
}
func (rs *Regex) ToString() string {
	rs.Close()

	cs := make([]string, 0)
	for _, m := range rs.Result.Ranges {
		cs = append(cs, m.Value)
	}
	return strings.Join(cs, "")
}

// func (rs *RegexFactory) export(template string) string {
// 	//r := rs.regex
// 	//r.ReplaceAllString(template,)
// 	cs := make([]string, 0)
// 	for _, m := range rs.Ranges {
// 		cs = append(cs, m.Value)
// 	}
// 	return strings.Join(cs, "")
// }
