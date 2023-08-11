package cmd

import (
	"regexp"
)

// Transform content from keyChar to valueChar defined in config.yaml
// keychar:valuechar= {'\n': '乚' , '\r': '刂', '\t': '亠'}
func (cfg *AppConfig) Transform(content *string) string {
	for key, val := range config.EChars {
		encode(content, key, val)
	}
	return *content
}

// Transform content from valueChar to keyChar defined in config.yaml
// keychar:valuechar= {'\n': '乚' , '\r': '刂', '\t': '亠'}
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

// replace all string by fchar to tchar
func encode(input *string, fchar, tchar string) {
	r, _ := regexp.Compile(fchar)
	*input = r.ReplaceAllString(*input, tchar)
}
