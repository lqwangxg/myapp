package cmd

func (rule *RegexRule) Eval(rs *RegexResult) {
	if rule.Formulas == nil || len(rule.Formulas) == 0 {
		return
	}
	for _, fomula := range rule.Formulas {
		fomula.EvalTrue()
	}
}
func (f *Formula) EvalTrue() bool {
	return true
}
