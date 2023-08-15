package cmd

import "regexp"

type FlagConfig struct {
	RuleName        string
	ConfigFile      string
	Pattern         string
	Content         string
	DestFile        string
	DestDir         string
	IncludeSuffix   string
	ExcludeSuffix   string
	ExportTemplate  string
	ReplaceTemplate string
}
type RegexAction int

type Regex struct {
	R         *regexp.Regexp
	Result    RegexResult
	Cache     bool
	CacheKey  string
	content   string
	fromFile  string
	IsReplace bool
	template  string
	Rule      *RuleConfig
}
type RegexResult struct {
	Pattern    string
	GroupNames []string
	Matches    []RegexMatch
	Ranges     []RegexRange
	Params     map[string]string
}

type RegexRange struct {
	*Capture
	MatchIndex int
}
type RegexMatch struct {
	*Capture
	Groups []RegexGroup
	Params map[string]string
}
type RegexGroup struct {
	*Capture
	Name string
}

// define RegexResultType as int.
type RegexResultType int

// RegexType of UnMatchType, MatchType and GroupType
const (
	UnMatchType RegexResultType = iota
	MatchType
	GroupType
)

type Capture struct {
	Start int
	End   int
	Value string
	RType RegexResultType
}

type Bound struct {
	Start int
	End   int
}

type AppConfig struct {
	EChars      map[string]string `yaml:"echars"`
	RuleDirs    []string          `yaml:"ruledirs"`
	Params      []string          `yaml:"params"`
	Indent      string            `yaml:"indent"`
	Prefix      string            `yaml:"prefix"`
	RedisOption RedisOption       `mapstructure:"redis"`
}
type RedisOption struct {
	Enable   bool   `yaml:"enable"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type RuleConfig struct {
	Name            string        `mapstructure:"name"`
	Group           string        `mapstructure:"group"`
	IncludeSuffix   string        `mapstructure:"include-suffix"`
	ExcludeSuffix   string        `mapstructure:"exclude-suffix"`
	RangeStart      string        `mapstructure:"range_start"`
	RangeEnd        string        `mapstructure:"range_end"`
	Pattern         string        `mapstructure:"pattern"`
	ExportTemplate  RegexTemplate `mapstructure:"export-template"`
	ReplaceTemplate RegexTemplate `mapstructure:"replace-template"`

	FullPatterns  CheckPatternConfig `mapstructure:"full_patterns"`
	RangePatterns CheckPatternConfig `mapstructure:"range_patterns"`
	MatchPatterns CheckPatternConfig `mapstructure:"match_patterns"`
	MatchFormulas []string           `mapstructure:"match_formulas"`
	MatchEvals    CheckPatternConfig `mapstructure:"match_evals"`
}
type RegexTemplate struct {
	Name   string `mapstructure:"name"`
	Header string `mapstructure:"header"`
	Match  string `mapstructure:"match"`
	Group  string `mapstructure:"group"`
	Footer string `mapstructure:"footer"`
}
type CheckPatternConfig struct {
	SkipIfany []string `mapstructure:"skip_ifany"`
	SkipWhen  []string `mapstructure:"skip_when"`
	DoIfany   []string `mapstructure:"do_ifany"`
	DoWhen    []string `mapstructure:"do_when"`
}
type TemplateControlConfig struct {
	Loop    string `mapstructure:"loop"`
	Process string `mapstructure:"process"`
}

type ConvertFunc func(*string, map[string]string) *string
