package cmd

import (
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/maja42/goval"
)

/*
 * equal, not equal, not : ==, !=, !
 * and/or: &&, ||
 * ternary ?: example: true? 1:2
 * Array in: "txt"/true/nil/2/[2,3] in ["txt", true, 2, [2,3]]
 * substrings[a:b]: [0,len(str)]
 * str function: len(str), contains(input, substr), match(input,pattern)
 */

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
			if reflect.ValueOf(p1).Kind() == reflect.Int {
				return false, nil
			}
			str := p1.(string)
			if str == "" {
				return true, nil
			}
			return false, nil
		}
		evalFuncs["match"] = func(args ...interface{}) (interface{}, error) {
			input := args[0].(string)
			pattern := args[1].(string)
			istrue := IsMatchString(pattern, input)
			return istrue, nil
		}
		evalFuncs["strjoin"] = func(args ...interface{}) (interface{}, error) {

			var buffer strings.Builder
			for _, v := range args {
				t := reflect.ValueOf(v)
				if t.Kind() == reflect.Int {
					buffer.WriteString(strconv.Itoa(v.(int)))
				} else if t.Kind() == reflect.Bool {
					buffer.WriteString(strconv.FormatBool(v.(bool)))
				} else {
					s := v.(string)
					if s != "" {
						buffer.WriteString(s)
					}
				}
			}
			return buffer.String(), nil
		}
	}
	return evalFuncs
}
