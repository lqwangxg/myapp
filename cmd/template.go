package cmd

import "log"

func (rs *Regex) ReplaceLoop(template *string, repFunc ConvertFunc) {
	rsLoop := NewNoCacheRegex(templateCtl.Loop)
	rsLoop.ScanMatches(*template)
	if len(rsLoop.Result.Matches) == 0 {
		return
	}

	*template = rsLoop.replaceText("${process}")
	repFunc(template, rs.Result.Params)
	isParamLoop := false
	for _, m := range rsLoop.Result.Matches {
		v, found := m.Params["items"]
		if found && v == "params" {
			isParamLoop = true
			break
		}
	}

	if isParamLoop {
		for _, m := range rs.Result.Matches {
			for key, value := range m.Params {
				ReplaceTemplateByKeyValue(template, "param.key", key)
				ReplaceTemplateByKeyValue(template, "param.value", value)
			}
		}
	}

	log.Print(template)
}
