package cmd

import (
	"log"
	"regexp"

	"golang.org/x/exp/slices"
)

/*
 * is skip check
 */
func (rule *RuleConfig) IsFullSkip(src string) bool {
	return rule.FullCheckPatterns.IsSkip(src)
}
func (rule *RuleConfig) IsFullDest(src string) bool {
	return rule.FullCheckPatterns.IsDest(src)
}
func (rule *RuleConfig) IsRangeSkip(src string) bool {
	return rule.RangeCheckPatterns.IsSkip(src)
}
func (rule *RuleConfig) IsRangeDest(src string) bool {
	return rule.RangeCheckPatterns.IsDest(src)
}

func (rule *RuleConfig) IsMatchSkip(src string) bool {
	return rule.MatchCheckPatterns.IsSkip(src)
}
func (rule *RuleConfig) IsMatchDest(src string) bool {
	return rule.MatchCheckPatterns.IsDest(src)
}

/*
 * is skip check
 */
func (checker *CheckPatternConfig) IsSkip(src string) bool {

	if IfAnyMatch(src, checker.SkipIfany) {
		return true
	}
	if WhenMatch(src, checker.SkipWhen) {
		return true
	}
	return false
}

/*
 * is dest check
 */
func (checker *CheckPatternConfig) IsDest(src string) bool {

	if IfAnyMatch(src, checker.DoIfany) {
		return true
	}
	if WhenMatch(src, checker.DoWhen) {
		return true
	}
	return false
}

/*
 * if empty, return false
 * if any true,  return true.
 */
func IfAnyMatch(src string, patterns []string) bool {
	var boolSlice []bool
	for _, pattern := range patterns {
		if pattern == "" {
			continue
		}

		if ok, err := regexp.MatchString(pattern, src); err == nil {
			boolSlice = append(boolSlice, ok)
		} else {
			log.Fatal(err)
		}
	}
	// if contains true, return true
	return slices.Contains(boolSlice, true)
}

/*
 * if empty, return false
 * if any false,  return false.
 */
func WhenMatch(src string, patterns []string) bool {
	var boolSlice []bool
	for _, pattern := range patterns {
		if pattern == "" {
			continue
		}

		if ok, err := regexp.MatchString(pattern, src); err == nil {
			boolSlice = append(boolSlice, ok)
		} else {
			log.Fatal(err)
		}
	}
	//if empty, return false
	if len(boolSlice) == 0 {
		return false
	}
	// if contains false, return false
	return !slices.Contains(boolSlice, false)
}
