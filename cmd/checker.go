package cmd

import (
	"log"
	"regexp"

	"golang.org/x/exp/slices"
)

func (rs *Regex) FullCheck(src string) bool {
	rule := rs.Rule
	if rule.IsFullSkip(src) {
		return false
	}
	if !rule.IsFullDest(src) {
		return false
	}
	return true
}
func (rs *Regex) IsDestMatch(src string, match RegexMatch) bool {
	// if rs.Rule.RangePattern == "" {
	// 	return true
	// }

	rsRange := NewNoCacheRegex(rs.Rule.RangeStart)
	rsRange.ScanMatches(src)
	rsRange.SplitMatches(src)
	innerCheck := false
	var inMatch RegexMatch
	for _, rm := range rsRange.Result.Matches {
		// match.position(start, end) is inner rm.position(start, end)
		if rm.Bound.Start <= match.Bound.Start && match.Bound.End <= rm.Bound.End {
			innerCheck = true
			inMatch = rm
			break
		}
	}
	// return false if not included in any range.
	if !innerCheck {
		return false
	}
	if rs.Rule.IsRangeSkip(inMatch.Value) {
		return false
	}
	if !rs.Rule.IsRangeDest(inMatch.Value) {
		return false
	}
	if rs.Rule.IsMatchSkip(match.Value) {
		return false
	}
	if !rs.Rule.IsMatchDest(match.Value) {
		return false
	}
	return true
}

/*
 * is skip check
 */
func (rule *RuleConfig) IsFullSkip(src string) bool {
	return rule.FullPatterns.IsSkip(src)
}
func (rule *RuleConfig) IsFullDest(src string) bool {
	return rule.FullPatterns.IsDest(src)
}
func (rule *RuleConfig) IsRangeSkip(src string) bool {
	return rule.RangePatterns.IsSkip(src)
}
func (rule *RuleConfig) IsRangeDest(src string) bool {
	return rule.RangePatterns.IsDest(src)
}

func (rule *RuleConfig) IsMatchSkip(src string) bool {
	return rule.MatchPatterns.IsSkip(src)
}
func (rule *RuleConfig) IsMatchDest(src string) bool {
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
