package cmd

import (
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/maja42/goval"
)

type EvalFunctions map[string]goval.ExpressionFunction

var evalFuncs EvalFunctions

func EvalFuncs() EvalFunctions {
	if evalFuncs == nil {
		evalFuncs = make(map[string]goval.ExpressionFunction)
		evalFuncs["len"] = func(args ...interface{}) (interface{}, error) {
			p1 := args[0]
			switch reflect.TypeOf(p1).Kind() {
			case reflect.String:
				str := p1.(string)
				return len(str), nil
			}
			return reflect.ValueOf(&p1).Len(), nil
		}
		evalFuncs["lenUTF8"] = func(args ...interface{}) (interface{}, error) {
			p1 := args[0]
			switch reflect.TypeOf(p1).Kind() {
			case reflect.String:
				str := p1.(string)
				return utf8.RuneCountInString(str), nil
			}
			return reflect.ValueOf(&p1).Len(), nil
		}
		evalFuncs["int"] = func(args ...interface{}) (interface{}, error) {
			p1 := args[0]
			str := p1.(string)
			return strconv.Atoi(str)
		}
	}
	return evalFuncs
}
