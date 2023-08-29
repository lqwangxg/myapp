package cmd

type AppConfig struct {
	EChars      map[string]string `yaml:"echars"`
	RuleDir     string            `yaml:"ruledir"`
	Params      map[string]string `yaml:"params"`
	Indent      string            `yaml:"indent"`
	Prefix      string            `yaml:"prefix"`
	RedisOption RedisOption       `mapstructure:"redis"`
	WebServer   WebServerOption   `yaml:"webserver"`
}

// type Regex struct {
// 	R           *regexp.Regexp
// 	Result      RegexResult
// 	Cache       bool
// 	CacheKey    string
// 	FromFile    string
// 	ToFile      string
// 	ExportFlag  bool
// 	ReplaceFlag bool
// 	Rule        *RegexRule
// }

type WebServerOption struct {
	Port string `yaml:"port"`
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
	Name        string `mapstructure:"name"`
	Group       string `mapstructure:"group"`
	IncludeFile string `mapstructure:"include-file"`
	ExcludeFile string `mapstructure:"exclude-file"`
	RangeStart  string `mapstructure:"range-start"`
	RangeEnd    string `mapstructure:"range-end"`
	Pattern     string `mapstructure:"pattern"`

	ParamPatterns `mapstructure:"param-patterns"`
	Formulas      *[]Formula `mapstructure:"formulas"`

	ExportTemplate      *RegexTemplate `mapstructure:"export-template"`
	ExportTemplateName  string         `mapstructure:"export-template-name"`
	ReplaceTemplate     *RegexTemplate `mapstructure:"replace-template"`
	ReplaceTemplateName string         `mapstructure:"replace-template-name"`

	FullCheckName  string     `mapstructure:"full-check-name"`
	FullCheckRule  *CheckRule `mapstructure:"full-check-rule"`
	RangeCheckName string     `mapstructure:"range-check-name"`
	RangeCheckRule *CheckRule `mapstructure:"range-check-rule"`
	MatchCheckName string     `mapstructure:"match-check-name"`
	MatchCheckRule *CheckRule `mapstructure:"match-check-rule"`
}
type RegexTemplate struct {
	Name        string `mapstructure:"name"`
	Header      string `mapstructure:"header"`
	Match       string `mapstructure:"match"`
	GroupHeader string `mapstructure:"group-header"`
	Group       string `mapstructure:"group"`
	GroupFooter string `mapstructure:"group-footer"`
	ParamHeader string `mapstructure:"param-header"`
	Param       string `mapstructure:"param"`
	ParamFooter string `mapstructure:"param-footer"`
	Footer      string `mapstructure:"footer"`
}

type Formula struct {
	If string `mapstructure:"if"`
	Do string `mapstructure:"do"`
}

type ParamPatterns struct {
	Inits   []string `mapstructure:"init"`
	Fulls   []string `mapstructure:"full"`
	Ranges  []string `mapstructure:"range"`
	Matches []string `mapstructure:"match"`
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
