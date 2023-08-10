package cmd

import (
	"log"
	"regexp"

	"golang.org/x/exp/slices"
)

/*
 * is skip check
 */
func (rule *ReplaceRuleConfig) IsFullSkip(src string) bool {
	return rule.FullPatterns.IsSkip(src)
}
func (rule *ReplaceRuleConfig) IsFullDest(src string) bool {
	return rule.FullPatterns.IsDest(src)
}
func (rule *ReplaceRuleConfig) IsRangeSkip(src string) bool {
	return rule.RangePatterns.IsSkip(src)
}
func (rule *ReplaceRuleConfig) IsRangeDest(src string) bool {
	return rule.RangePatterns.IsDest(src)
}

func (rule *ReplaceRuleConfig) IsMatchSkip(src string) bool {
	return rule.MatchPatterns.IsSkip(src)
}
func (rule *ReplaceRuleConfig) IsMatchDest(src string) bool {
	return rule.MatchPatterns.IsDest(src)
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
