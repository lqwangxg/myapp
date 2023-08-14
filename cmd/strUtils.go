package cmd

import (
	"fmt"
	"regexp"
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

func ReplaceByKeyValue(template *string, key string, value string) {
	pattern := fmt.Sprintf(`\$\{(?P<key>%s)\}`, key)
	regex := NewNoCacheRegex(pattern)
	regex.ScanMatches(*template)
	if ok, _ := regex.ParamValue(key); !ok {
		return
	}
	*template = regex.R.ReplaceAllString(*template, value)
}

func ReplaceByMap(template *string, kvs map[string]string) *string {
	for key, val := range kvs {
		ReplaceByKeyValue(template, key, val)
	}
	return template
}
