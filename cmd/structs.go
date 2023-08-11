package cmd

import "regexp"

type FlagConfig struct {
	Name       string
	ConfigFile string
	Pattern    string
	Content    string
	DestFile   string
	DestDir    string
	Suffix     string
}
type Context struct {
	Params map[string]string
}
type RegexFactory struct {
	Pattern    string
	GroupNames []string
	Matches    []RegexMatch
	Regex      *regexp.Regexp
	Content    string //match input string
	Ranges     []RegexRange
	Params     map[string]string
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
	Name               string             `yaml:"name"`
	Group              string             `yaml:"group"`
	FullCheckPatterns  CheckPatternConfig `yaml:"full_patterns"`
	RangePattern       string             `yaml:"range_pattern"`
	RangeParamsPattern []string           `yaml:"range_params_pattern"`
	RangeCheckPatterns CheckPatternConfig `yaml:"range_patterns"`
	MatchPattern       string             `yaml:"match_pattern"`
	MatchCheckPatterns CheckPatternConfig `yaml:"match_patterns"`
	MatchFormulas      []string           `yaml:"match_formulas"`
	MatchCheckEvals    CheckPatternConfig `yaml:"match_evals"`
	Key                string
}
type CheckPatternConfig struct {
	SkipIfany []string `yaml:"skip_ifany"`
	SkipWhen  []string `yaml:"skip_when"`
	DoIfany   []string `yaml:"do_ifany"`
	DoWhen    []string `yaml:"do_when"`
}
