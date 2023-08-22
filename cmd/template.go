package cmd

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/maja42/goval"
	"github.com/pkg/errors"
)

const (
	keypattern   = `(?P<key>[\w\.\-]+)`
	anykey       = `\$\{` + keypattern + `\}`
	fmtkey       = `\$\{(?P<key>%s)\}`
	kindkey      = `kind:\s*` + keypattern
	asignformula = anykey + `\s*=\s*(?P<formula>.+)$`
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
func ToParamKey(key string) string {
	return fmt.Sprintf(fmtkey, key)
}
func FromParamKey(key string) string {
	if !IsMatchString(anykey, key) {
		return key
	}
	return NewTemplate(key).GetMatchString(keypattern)
}

// return true if contains ${key}.
func (t *StringTemplate) HasKey(key string) bool {
	return IsMatchString(ToParamKey(key), t.Template)
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

	regex := regexp.MustCompile(ToParamKey(key))
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
	// eval := goval.NewEvaluator()
	// result, err := eval.Evaluate(t.Template, nil, nil)
	// if err != nil {
	// 	//fmt.Printf("Error: %s\n", err)
	// 	return false, err
	// }

	// //fmt.Printf("%v (%s)\n\n", result, reflect.TypeOf(result))
	// if reflect.TypeOf(result).Kind() != reflect.Bool {
	// 	return false, errors.Errorf("result %v is not bool.", result)
	// }
	// v := reflect.ValueOf(&result).Elem()
	// ret := v.Interface().(bool)
	// //fmt.Printf("eval(%s)=%v", t.Template, ret)
	// return ret, nil
	v, err := t.Eval(reflect.Bool)
	if err == nil {
		return v.Interface().(bool), nil
	} else {
		return false, err
	}
}
func (t *StringTemplate) EvalInt() (int, error) {
	v, err := t.Eval(reflect.Int)
	if err == nil {
		return v.Interface().(int), nil
	} else {
		return 0, err
	}
}
func (t *StringTemplate) Eval(kind reflect.Kind) (reflect.Value, error) {
	eval := goval.NewEvaluator()
	result, err := eval.Evaluate(t.Template, nil, nil)
	if err != nil {
		//fmt.Printf("Error: %s\n", err)
		return reflect.ValueOf(result), err
	}

	//fmt.Printf("%v (%s)\n\n", result, reflect.TypeOf(result))
	if reflect.TypeOf(result).Kind() != kind {
		return reflect.ValueOf(result), errors.Errorf("result %v is not %s.", result, kind.String())
	}
	//v := reflect.ValueOf(&result).Elem()
	//ret := v.Interface().(reflect.Value)
	//fmt.Printf("eval(%s)=%v", t.Template, ret)
	return reflect.ValueOf(&result).Elem(), nil
}
func (t *StringTemplate) GetMatchString(pattern string) string {
	regex := regexp.MustCompile(pattern)
	if regex.MatchString(t.Template) {
		return regex.FindString(t.Template)
	}
	return ""
}
func (t *StringTemplate) GetKey(pattern string) string {
	regex := regexp.MustCompile(pattern)
	if regex.MatchString(t.Template) {
		return regex.FindString(t.Template)
	}
	return t.Template
}
func (t *StringTemplate) ResetParam(origin map[string]string) {
	// ${key} = ${formula}
	isMatched, matches := NewRegexText(asignformula, t.Template).GetMatchResult(true)
	if !isMatched {
		return
	}

	match := matches.FirstMatch()
	if match == nil {
		return
	}

	evTemp, changed := NewTemplate(match.Params["formula"]).ReplaceByMap(origin)
	if !changed {
		return
	}
	i, err := NewTemplate(evTemp).EvalInt()
	if err != nil {
		return
	}
	key := FromParamKey(match.Params["key"])
	origin[key] = strconv.Itoa(i)
}
