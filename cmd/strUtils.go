package cmd

import (
	"regexp"
)

func (cfg *AppConfig) Parse(content *string) string {
	for key, val := range config.EChars {
		encode(content, key, val)
	}
	return *content
}

func (cfg *AppConfig) Restore(content *string) string {
	for key, val := range config.EChars {
		ckey := key
		switch key {
		case "\\n":
			ckey = "\n"
		case "\\r":
			ckey = "\r"
		case "\\t":
			ckey = "\t"
		}
		encode(content, val, ckey)
	}
	return *content
}

func encode(input *string, fchar, tchar string) {
	r, _ := regexp.Compile(fchar)
	*input = r.ReplaceAllString(*input, tchar)
}
