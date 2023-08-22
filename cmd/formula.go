package cmd

type paramMap map[string]string

func (p *paramMap) Eval(fomula *Formula) {
	if fomula.If == "" {
		return
	}
	tmpIf, changed := NewTemplate(fomula.If).ReplaceByMap(*p)
	if !changed {
		return
	}

	isTrue, err := NewTemplate(tmpIf).EvalTrue()
	if err != nil {
		return
	}
	if !isTrue {
		return
	}
	NewTemplate(fomula.Do).ResetParam(*p)
}
func (rs *RegexResult) Eval(formula *Formula) {
	(*paramMap)(&rs.Params).Eval(formula)
}

func (rs *Capture) Eval(formula *Formula) {
	(*paramMap)(&rs.Params).Eval(formula)
}
