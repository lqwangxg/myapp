package cmd

import (
	"regexp"
)

func beforeMatch(content string) string {
	ret := content
	for key, val := range config.echars {
		ret = replaceText(ret, key, val)
	}
	return ret
}

func afterMatch(content string) string {
	ret := content
	for key, val := range config.echars {
		ret = replaceText(ret, val, key)
	}
	return ret
}

func replaceText(input, pattern, replacement string) string {
	r, _ := regexp.Compile(pattern)
	return r.ReplaceAllString(input, replacement)
}
