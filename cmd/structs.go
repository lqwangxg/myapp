package cmd

type FlagConfig struct {
	name       string
	configFile string
	pattern    string
	origin     string
	destFile   string
	destDir    string
	suffix     string
}

type RegexMatchResult struct {
	Pattern    string
	GroupNames []string
	Matches    []RegexMatch
}
type RegexMatch struct {
	Index    int
	Value    string
	Groups   []RegexGroup
	Position RegexMatchIndex
}
type RegexGroup struct {
	Index    int
	Name     string
	Value    string
	Position RegexMatchIndex
}
type RegexMatchIndex struct {
	Start int
	End   int
	//Length int //=EndIndex-StartIndex
}

type AppConfig struct {
	echars  map[string]string
	pparams []string
}
type RuleConfig struct {
	name                 string
	group                string
	full_patterns        CheckPatternConfig
	range_pattern        string
	range_patterns       CheckPatternConfig
	range_params_pattern []string
	match_pattern        string
	match_patterns       CheckPatternConfig
	match_formulas       []string
	match_evals          CheckPatternConfig
}
type CheckPatternConfig struct {
	skip_ifany []string
	skip_when  []string
	do_ifany   []string
	do_when    []string
}
