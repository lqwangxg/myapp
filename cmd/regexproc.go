package cmd

type IRegexHandler interface {
	Execute()
}

func Exec(hand IRegexHandler) {
	wg.Add(1)
	defer wg.Done()
	hand.Execute()
}
