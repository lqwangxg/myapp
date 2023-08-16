package cmd

type IRegexHandler interface {
	Match()
	Write()
}

type RegexManager struct {
	RegexResult
}

var reger RegexManager

func (manager *RegexManager) Execute(hand IRegexHandler) {
	if manager == nil {
		manager = &RegexManager{}
	}
	hand.Match()
	hand.Write()
}
