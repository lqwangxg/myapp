package cmd

import (
	"regexp"
)

type Regex struct {
	R           *regexp.Regexp
	Result      RegexResult
	Cache       bool
	CacheKey    string
	FromFile    string
	ToFile      string
	ExportFlag  bool
	ReplaceFlag bool
	Rule        *RegexRule
}

type RedisOption struct {
	Enable   bool   `yaml:"enable"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type RegexRules struct {
	Rules []RegexRule `mapstructure:"regex-rules"`
}

type RegexRule struct {
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
type CheckRules struct {
	Rules []CheckRule `mapstructure:"check-rules"`
}
type CheckRule struct {
	Name     string   `mapstructure:"name"`
	IsOr     bool     `mapstructure:"isor"`
	Skip     bool     `mapstructure:"skip"`
	Patterns []string `mapstructure:"patterns"`
	Formulas []string `mapstructure:"formulas"`
}
