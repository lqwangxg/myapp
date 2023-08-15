package cmd

import "fmt"

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

func (t *StringTemplate) ReplaceByKeyValue(key string, value string) string {
	if !t.HasKey(key) {
		return t.Template
	}

	regex := NewNoCacheRegex(t.Pattern(key))
	regex.ScanMatches(t.Template)
	return regex.R.ReplaceAllString(t.Template, value)
}

func (t *StringTemplate) ReplaceByMap(kvs map[string]string) string {
	if !t.HasAnyKey() {
		return t.Template
	}

	for key, val := range kvs {
		t.Template = t.ReplaceByKeyValue(key, val)
	}
	return t.Template
}

func (t *StringTemplate) ReplaceByRegexResult(result RegexResult) string {
	t.Template = t.ReplaceByMap(result.Params)
	for _, m := range result.Matches {
		t.Template = t.ReplaceByMap(m.Params)
	}
	return t.Template
}
