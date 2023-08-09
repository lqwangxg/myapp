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
