package cmd

import "strings"

func getControlTemplate(template string) string {
	rstart := NewCacheRegex(templateCtl.LoopStart, false)
	rsend := NewCacheRegex(templateCtl.LoopEnd, false)
	if rstart.IsMatch(template) && rsend.IsMatch(template) {
		var sb strings.Builder
		sb.WriteString(templateCtl.LoopStart)
		sb.WriteString("(?P<process>.*)")
		sb.WriteString(templateCtl.LoopEnd)
		rs := NewCacheRegex(flags.Pattern, false)
		rs.ScanMatches(template)
		return rs.Result.Params["process"]
		//rs.replaceText()
	}
	return template
}

func (rs *Regex) ReplaceLoop(template *string, repFunc ConvertFunc) string {
	newTmp := getControlTemplate(*template)
	repFunc(&newTmp, rs.Result.Params)
	for _, m := range rs.Result.Matches {
		repFunc(&newTmp, m.Params)
		for _, g := range m.Groups {
			gmap := make(map[string]string)
			gmap["group.key"] = g.Name
			gmap["group.value"] = g.Value
			repFunc(&newTmp, gmap)
		}
	}
	return newTmp
}
