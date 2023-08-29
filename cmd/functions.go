package cmd

import (
	"reflect"
	"strconv"
	"strings"
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
		evalFuncs["isEmpty"] = func(args ...interface{}) (interface{}, error) {
			p1 := args[0]
			str := p1.(string)
			if str == "" {
				return true, nil
			}
			return false, nil
		}
		evalFuncs["match"] = func(args ...interface{}) (interface{}, error) {
			pattern := args[0].(string)
			input := args[1].(string)
			istrue := IsMatchString(pattern, input)
			return istrue, nil
		}
		evalFuncs["strjoin"] = func(args ...interface{}) (interface{}, error) {

			var buffer strings.Builder
			for _, v := range args {
				s := v.(string)
				if s != "" {
					buffer.WriteString(s)
				}
			}
			return buffer.String(), nil
		}
	}
	return evalFuncs
}
