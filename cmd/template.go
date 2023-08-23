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
	asignformula = `[\$\{]?` + keypattern + `[\}]?` + `\s*=\s*(?P<formula>.+)$`
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
		t.Template = regex.ReplaceAllLiteralString(t.Template, value)
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

func (t *StringTemplate) EvalTrue(p map[string]string) (bool, error) {
	v, err := t.Eval2(reflect.Bool, p)
	if err == nil {
		return v.Interface().(bool), nil
	} else {
		return false, err
	}
}
func (t *StringTemplate) EvalInt(p map[string]string) (int, error) {
	v, err := t.Eval2(reflect.Int, p)
	if err == nil {
		return v.Interface().(int), nil
	} else {
		return 0, err
	}
}

func (t *StringTemplate) Eval2(kind reflect.Kind, params map[string]string) (reflect.Value, error) {
	variables := make(map[string]any)
	for key, val := range params {
		if IsMatchString(`^\d+$`, val) {
			if v, err := strconv.Atoi(val); err == nil {
				variables[key] = v
				continue
			}
		}
		variables[key] = val
	}
	return t.Evaluate(kind, variables)
}

func (t *StringTemplate) Evaluate(kind reflect.Kind, variables map[string]any) (reflect.Value, error) {
	eval := goval.NewEvaluator()
	result, err := eval.Evaluate(t.Template, variables, EvalFuncs())
	if err != nil {
		fmt.Printf("%v", err)
		return reflect.ValueOf(result), err
	}

	if reflect.TypeOf(result).Kind() != kind {
		return reflect.ValueOf(result), errors.Errorf("result %v is not %s.", result, kind.String())
	}

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

	evTemp, _ := NewTemplate(match.Params["formula"]).ReplaceByMap(origin)
	i, err := NewTemplate(evTemp).EvalInt(origin)
	if err != nil {
		return
	}
	key := FromParamKey(match.Params["key"])
	origin[key] = strconv.Itoa(i)
}
func (t *StringTemplate) append(paramMap map[string]string) string {
	if t.Template == "" {
		return ""
	}
	header, _ := t.ReplaceByMap(paramMap)
	return header
}
