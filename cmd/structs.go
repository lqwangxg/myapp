package cmd

import "regexp"

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
	r          *regexp.Regexp
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
}

type AppConfig struct {
	EChars   map[string]string `yaml:"echars"`
	RuleDirs []string          `yaml:"ruledirs"`
	Params   []string          `yaml:"params"`
}
type ReplaceRuleConfig struct {
	Name                 string             `yaml:"name"`
	Group                string             `yaml:"group"`
	Full_patterns        CheckPatternConfig `yaml:"full_patterns"`
	Range_pattern        string             `yaml:"range_pattern"`
	Range_params_pattern []string           `yaml:"range_params_pattern"`
	Range_patterns       CheckPatternConfig `yaml:"range_patterns"`
	Match_pattern        string             `yaml:"match_pattern"`
	Match_patterns       CheckPatternConfig `yaml:"match_patterns"`
	Match_formulas       []string           `yaml:"match_formulas"`
	Match_evals          CheckPatternConfig `yaml:"match_evals"`
	Key                  string
}
type CheckPatternConfig struct {
	Skip_ifany []string `yaml:"skip_ifany"`
	Skip_when  []string `yaml:"skip_when"`
	Do_ifany   []string `yaml:"do_ifany"`
	Do_when    []string `yaml:"do_when"`
}
