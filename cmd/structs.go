package cmd

import (
	"regexp"
)

type AppConfig struct {
	EChars      map[string]string `yaml:"echars"`
	RuleDir     string            `yaml:"ruledir"`
	Params      map[string]string `yaml:"params"`
	Indent      string            `yaml:"indent"`
	Prefix      string            `yaml:"prefix"`
	RedisOption RedisOption       `mapstructure:"redis"`
}

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
	Name                string         `mapstructure:"name"`
	Group               string         `mapstructure:"group"`
	IncludeFile         string         `mapstructure:"include-file"`
	ExcludeFile         string         `mapstructure:"exclude-file"`
	RangeStart          string         `mapstructure:"range-start"`
	RangeEnd            string         `mapstructure:"range-end"`
	Pattern             string         `mapstructure:"pattern"`
	ExportTemplate      *RegexTemplate `mapstructure:"export-template"`
	ExportTemplateName  string         `mapstructure:"export-template-name"`
	ReplaceTemplate     *RegexTemplate `mapstructure:"replace-template"`
	ReplaceTemplateName string         `mapstructure:"replace-template-name"`

	FullCheckName  string `mapstructure:"full-check-name"`
	RangeCheckName string `mapstructure:"range-check-name"`
	MatchCheckName string `mapstructure:"match-check-name"`
}
type RegexTemplate struct {
	Name   string `mapstructure:"name"`
	Header string `mapstructure:"header"`
	Match  string `mapstructure:"match"`
	Group  string `mapstructure:"group"`
	Footer string `mapstructure:"footer"`
}

//	type CheckPatternConfig struct {
//		SkipIfany []string `mapstructure:"skip_ifany"`
//		SkipWhen  []string `mapstructure:"skip_when"`
//		DoIfany   []string `mapstructure:"do_ifany"`
//		DoWhen    []string `mapstructure:"do_when"`
//	}
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
