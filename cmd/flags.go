package cmd

type FlagConfig struct {
	Action          string
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

var flags FlagConfig
