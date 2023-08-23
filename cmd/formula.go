package cmd

type paramMap map[string]string

func (p *paramMap) Eval(formula *Formula) {
	if !p.evalIf(formula) {
		return
	}
	NewTemplate(formula.Do).ResetParam(*p)
}
func (p *paramMap) evalIf(formula *Formula) bool {
	if formula.If == "" {
		return true
	}

	tmpIf, _ := NewTemplate(formula.If).ReplaceByMap(*p)

	isTrue, err := NewTemplate(tmpIf).EvalTrue(*p)
	if err != nil {
		return false
	}
	return isTrue
}

func (rs *RegexResult) Eval(formula *Formula) {
	(*paramMap)(&rs.Params).Eval(formula)
}

func (rs *Capture) Eval(formula *Formula) {
	(*paramMap)(&rs.Params).Eval(formula)
}
