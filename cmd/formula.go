package cmd

type paramMap map[string]string

func (p *paramMap) Eval(formula *Formula) string {
	if isTrue, errMsg := p.evalIf(formula); !isTrue {
		return "evalIf:false, errMsg:" + errMsg
	}
	return NewTemplate(formula.Do).ResetParam(*p)
}
func (p *paramMap) evalIf(formula *Formula) (bool, string) {
	if formula.If == "" {
		return true, ""
	}

	tmpIf, _ := NewTemplate(formula.If).ReplaceByMap(*p)

	isTrue, err := NewTemplate(tmpIf).EvalTrue(*p)
	if err != nil {
		return false, err.Error()
	}
	return isTrue, ""
}

// func (rs *RegexResult) Eval(formula *Formula) {
// 	(*paramMap)(&rs.Params).Eval(formula)
// }

func (rs *Capture) Eval(formula *Formula) string {
	return (*paramMap)(&rs.Params).Eval(formula)
}
