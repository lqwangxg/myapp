package cmd

import (
	"regexp"
	"strings"
)

// Encode content from keyChar to valueChar defined in config.yaml
// keychar:valuechar= {'\n': '乚' , '\r': '刂', '\t': '亠'}
func (cfg *AppConfig) Encode(content *string) string {
	for key, val := range config.EChars {
		encode(content, key, val)
	}
	return *content
}

// Transform content from valueChar to keyChar defined in config.yaml
// keychar:valuechar= {'\n': '乚' , '\r': '刂', '\t': '亠'}
func (cfg *AppConfig) Decode(content *string) string {
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

func (cfg *AppConfig) EncodePattern(input *string) string {
	for key, val := range config.EChars {
		*input = strings.ReplaceAll(*input, key, val)
	}
	return *input
}
func (cfg *AppConfig) DecodePattern(input *string) string {
	for key, val := range config.EChars {
		*input = strings.ReplaceAll(*input, val, key)
	}
	return *input
}

// replace all string by fchar to tchar
func encode(input *string, fchar, tchar string) {
	r, _ := regexp.Compile(fchar)
	*input = r.ReplaceAllString(*input, tchar)
}
func (rs *Regex) ParamValue(key string) (bool, string) {
	for _, m := range rs.Result.Matches {
		for _, val := range m.Params {
			if val == key {
				return true, m.Params[key]
			}
		}
	}
	return false, ""
}
func IsMatchString(pattern, input string) bool {
	//pattern = `\$\{(?P<key>[\w\.\-]+)\}`
	r := regexp.MustCompile(pattern)
	return r.MatchString(input)
}
