package cmd

import "regexp"

type FlagConfig struct {
	RuleName          string
	ConfigFile        string
	Pattern           string
	Content           string
	DestFile          string
	DestDir           string
	IncludeSuffix     string
	ExcludeSuffix     string
	TemplateOfReplace string
	TempleteOfExport  string
	TemplateFile      string
}
type RegexAction int

const (
	MatchAction RegexAction = iota
	ReplaceAction
)

type Regex struct {
	R        *regexp.Regexp
	Result   RegexResult
	CacheKey string
	Cache    bool
	Action   RegexAction
	Rule     *RuleConfig
}
type RegexResult struct {
	Pattern    string
	GroupNames []string
	Matches    []RegexMatch
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
	Params   map[string]string
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
	Name            string `mapstructure:"name"`
	Group           string `mapstructure:"group"`
	IncludeSuffix   string `mapstructure:"include-suffix"`
	ExcludeSuffix   string `mapstructure:"exclude-suffix"`
	Pattern         string `mapstructure:"pattern"`
	ExportTemplate  string `mapstructure:"export-template"`
	ReplaceTemplate string `mapstructure:"replace-template"`
	RangePattern    string `mapstructure:"range_pattern"`

	FullPatterns  CheckPatternConfig `mapstructure:"full_patterns"`
	RangePatterns CheckPatternConfig `mapstructure:"range_patterns"`
	MatchPatterns CheckPatternConfig `mapstructure:"match_patterns"`
	MatchFormulas []string           `mapstructure:"match_formulas"`
	MatchEvals    CheckPatternConfig `mapstructure:"match_evals"`
}
type CheckPatternConfig struct {
	SkipIfany []string `mapstructure:"skip_ifany"`
	SkipWhen  []string `mapstructure:"skip_when"`
	DoIfany   []string `mapstructure:"do_ifany"`
	DoWhen    []string `mapstructure:"do_when"`
}
