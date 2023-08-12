package cmd

import (
	"fmt"
	"log"
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
func (rs *Regex) ContainsGroupKey(key string) bool {
	for _, m := range rs.Result.Matches {
		for _, val := range m.Params {
			if val == key {
				return true
			}
		}
	}
	return false
}

func ReplaceTemplateByKeyValue(template *string, key, value string) {
	pattern := fmt.Sprintf(`\$\{(?P<key>%s)\}`, key)
	regex := NewCacheRegex(pattern, false)
	regex.ScanMatches(*template)
	if !regex.ContainsGroupKey(key) {
		return
	}
	*template = regex.R.ReplaceAllString(*template, value)
}

func ReplaceTemplate(template *string, kvs map[string]string) *string {
	for key, val := range kvs {
		ReplaceTemplateByKeyValue(template, key, val)
	}
	return template
}

func TestReplaceMap() {
	input := "${t1}+${t2}=${t3}"
	kvs := map[string]string{
		"t1": "1000",
		"t2": "1500",
		"t3": "2500",
	}
	//key := "t1"
	//val := "1000"
	//result := ReplaceTemplateByKeyValue(input, key, val)
	result := ReplaceTemplate(&input, kvs)
	log.Println(result)
}
