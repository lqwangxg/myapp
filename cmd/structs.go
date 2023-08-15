package cmd

import "regexp"

type FlagConfig struct {
	RuleName        string
	ConfigFile      string
	Pattern         string
	ExportFlag      bool
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
	R           *regexp.Regexp
	Result      RegexResult
	Cache       bool
	CacheKey    string
	FromFile    string
	ToFile      string
	ExportFlag  bool
	ReplaceFlag bool
	Rule        *RuleConfig
}
type RegexResult struct {
	Pattern    string
	GroupNames []string
	Matches    []RegexMatch
	Captures   []Capture
	Params     map[string]string
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

type AppConfig struct {
	EChars      map[string]string `yaml:"echars"`
	RuleDir     string            `yaml:"ruledir"`
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

type CustomRules struct {
	Rules []RuleConfig `mapstructure:"custom-rules"`
}

type RuleConfig struct {
	Name            string        `mapstructure:"name"`
	Group           string        `mapstructure:"group"`
	IncludeFile     string        `mapstructure:"include-file"`
	ExcludeFile     string        `mapstructure:"exclude-file"`
	RangeStart      string        `mapstructure:"range-start"`
	RangeEnd        string        `mapstructure:"range-end"`
	Pattern         string        `mapstructure:"pattern"`
	ExportTemplate  RegexTemplate `mapstructure:"export-template"`
	ReplaceTemplate RegexTemplate `mapstructure:"replace-template"`

	FullPatterns  CheckPatternConfig `mapstructure:"full-patterns"`
	RangePatterns CheckPatternConfig `mapstructure:"range-patterns"`
	MatchPatterns CheckPatternConfig `mapstructure:"match-patterns"`
	MatchFormulas []string           `mapstructure:"match-formulas"`
	MatchEvals    CheckPatternConfig `mapstructure:"match-evals"`
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
