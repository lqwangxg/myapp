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
	PATTERN_KEY          = `(?P<key>[\w\.\-]+)`
	PATTERN_ANYKEY       = `\$\{` + PATTERN_KEY + `\}`
	PATTERN_REFKEY       = `\$\{(?P<key>%s)\}`
	PATTERN_KIND_KEY     = `kind:\s*` + PATTERN_KEY
	PATTERN_FORMULA_DO   = `[\$\{]?` + PATTERN_KEY + `[\}]?` + `\s*=\s*(?P<formula>.+)$`
	PATTERN_FORMULA_BOOL = `is\w+\(|\|\||\&\&|\!`
	PATTERN_FORMULA_INT  = `\s*[\+\-\*\/\%]\s*|len\(|len\w+\(`
	PATTERN_FORMULA_WORD = `^\s*` + PATTERN_KEY + `\s*$`
	PATTERN_RESERVED_KEY = `^(match|group|param).\w+$`
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
	return fmt.Sprintf(PATTERN_REFKEY, key)
}
func FromParamKey(key string) string {
	if !IsMatchString(PATTERN_ANYKEY, key) {
		return key
	}
	return NewTemplate(key).GetMatchString(PATTERN_KEY)
}

// return true if contains ${key}.
func (t *StringTemplate) HasKey(key string) bool {
	return IsMatchString(ToParamKey(key), t.Template)
}

// return true if contains any key.
// key chars accept: \w+\.\-
func (t *StringTemplate) HasAnyKey() bool {
	return IsMatchString(PATTERN_ANYKEY, t.Template)
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
		} else if IsMatchString(`^true|false$`, val) {
			if v, err := strconv.ParseBool(val); err == nil {
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
	// skip if empty
	if t.Template == "" {
		return
	}

	key, assign := t.ToKeyValue()
	evTemp, _ := NewTemplate(assign).ReplaceByMap(origin)
	if IsMatchString(PATTERN_FORMULA_BOOL, assign) {
		i, err := NewTemplate(evTemp).EvalTrue(origin)
		if err != nil {
			return
		}
		origin[key] = strconv.FormatBool(i)
	} else if IsMatchString(PATTERN_FORMULA_INT, assign) {
		i, err := NewTemplate(evTemp).EvalInt(origin)
		if err != nil {
			return
		}
		origin[key] = strconv.Itoa(i)
	} else if IsMatchString(PATTERN_FORMULA_WORD, assign) {
		if _, ok := origin[assign]; ok {
			origin[key] = origin[assign]
		}
	} else {
		fmt.Printf("skip formula:%s, match.value:%v\n", t.Template, origin["match.value"])
	}
}
func (t *StringTemplate) append(paramMap map[string]string) string {
	if t.Template == "" {
		return ""
	}
	header, _ := t.ReplaceByMap(paramMap)
	return header
}
func (t *StringTemplate) ToKeyValue() (key, value string) {
	r := regexp.MustCompile(`\s*=\s*`)
	strArray := r.Split(t.Template, 2)
	key = strArray[0]
	value = strArray[1]
	return key, value
}
