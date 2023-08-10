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
type RegexFactory struct {
	Pattern    string
	GroupNames []string
	Matches    []RegexMatch
	r          *regexp.Regexp
	content    string //match input string
	Ranges     []RegexRange
}

type RegexRange struct {
	Value      string
	IsMatch    bool
	MatchIndex int
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
	FullPatterns         CheckPatternConfig `yaml:"full_patterns"`
	Range_pattern        string             `yaml:"range_pattern"`
	Range_params_pattern []string           `yaml:"range_params_pattern"`
	RangePatterns        CheckPatternConfig `yaml:"range_patterns"`
	Match_pattern        string             `yaml:"match_pattern"`
	MatchPatterns        CheckPatternConfig `yaml:"match_patterns"`
	Match_formulas       []string           `yaml:"match_formulas"`
	Match_evals          CheckPatternConfig `yaml:"match_evals"`
	Key                  string
}
type CheckPatternConfig struct {
	SkipIfany []string `yaml:"skip_ifany"`
	SkipWhen  []string `yaml:"skip_when"`
	DoIfany   []string `yaml:"do_ifany"`
	DoWhen    []string `yaml:"do_when"`
}
