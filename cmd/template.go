package cmd

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/maja42/goval"
	"github.com/pkg/errors"
)

const (
	anykey  = `\$\{[\w\.\-]+\}`
	fmtkey  = `\$\{(?P<key>%s)\}`
	kindkey = `kind:\s*(?P<key>[\w\.\-]+)`
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
	changed := false
	if regex.MatchString(t.Template) {
		t.Template = regex.ReplaceAllString(t.Template, value)
		changed = true
	}
	return t.Template, changed

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

func (t *StringTemplate) ReplaceBy(m Capture) (string, bool) {
	if m.IsMatch && m.Params != nil {
		return t.ReplaceByMap(m.Params)
	}
	return t.Template, false
}

func (t *StringTemplate) EvalTrue() (bool, error) {
	eval := goval.NewEvaluator()
	result, err := eval.Evaluate(t.Template, nil, nil)
	if err != nil {
		//fmt.Printf("Error: %s\n", err)
		return false, err
	}

	fmt.Printf("%v (%s)\n\n", result, reflect.TypeOf(result))
	if reflect.TypeOf(result).Kind() != reflect.Bool {
		return false, errors.Errorf("result %v is not bool.", result)
	}
	v := reflect.ValueOf(&result).Elem()
	ret := v.Interface().(bool)
	//fmt.Printf("eval(%s)=%v", t.Template, ret)
	return ret, nil
}
