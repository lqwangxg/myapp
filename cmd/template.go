package cmd

import (
	"fmt"
	"regexp"
)

const (
	anykey = `\$\{[\w\.\-]+\}`
	fmtkey = `\$\{(?P<key>%s)\}`
)

type StringTemplate struct {
	Template string
}
type ITemplate interface {
	HasKey(key string) bool
}

func NewTemplate(s string) *StringTemplate {
	return &StringTemplate{s}
}
func (t *StringTemplate) Pattern(key string) string {
	return fmt.Sprintf(fmtkey, key)
}

// return true if contains ${key}.
func (t *StringTemplate) HasKey(key string) bool {
	return IsMatchString(t.Pattern(key), t.Template)
}

// return true if contains any key.
// key chars accept: \w+\.\-
func (t *StringTemplate) HasAnyKey() bool {
	return IsMatchString(anykey, t.Template)
}

func (t *StringTemplate) ReplaceByKeyValue(key string, value string) (string, bool) {
	if !t.HasKey(key) {
		return t.Template, false
	}

	regex := regexp.MustCompile(t.Pattern(key))
	if regex.MatchString(t.Template) {
		return regex.ReplaceAllString(t.Template, value), true
	} else {
		return t.Template, false
	}
}

func (t *StringTemplate) ReplaceByMap(kvs map[string]string) (string, bool) {
	if !t.HasAnyKey() {
		return t.Template, false
	}
	hasChanged := false
	for key, val := range kvs {
		_, changed := t.ReplaceByKeyValue(key, val)
		if changed {
			hasChanged = true
		}
	}
	return t.Template, hasChanged
}

// func (t *StringTemplate) ReplaceByRegexResult(result RegexResult) string {
// 	t.Template = t.ReplaceByMap(result.Params)
// 	for _, m := range result.Captures {
// 		if m.IsMatch && m.Params != nil {
// 			t.Template = t.ReplaceByMap(m.Params)
// 		}
// 	}
// 	return t.Template
// }

func (t *StringTemplate) ReplaceBy(m Capture) (string, bool) {
	if m.IsMatch && m.Params != nil {
		return t.ReplaceByMap(m.Params)
	}
	return t.Template, false
}
