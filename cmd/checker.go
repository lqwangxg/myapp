package cmd

import (
	"log"
	"regexp"

	"golang.org/x/exp/slices"
)

// check whether the content is should be skiped by skip-pattern,
// or is required by dest-pattern
func (rs *Regex) FullCheck(content string) bool {
	rule := rs.Rule
	if rule.IsFullSkip(content) {
		return false
	}
	if !rule.IsFullDest(content) {
		return false
	}
	return true
}

// check the given match[start,end] whether between the range of matches
// by range-start and range-end
func (rs *Regex) IsDestMatch(match Capture, content string) bool {
	ranges := rs.Rule.MergeRangeStartEnd(content)
	inRange := false
	rangeValue := ""
	for _, r := range *ranges {
		if r.Start <= match.Start && match.End <= r.End {
			inRange = true
			rangeValue = content[r.Start:r.End]
			break
		}
	}

	// return false if not included in any range.
	if !inRange {
		return false
	}
	if rs.Rule.IsRangeSkip(rangeValue) {
		return false
	}
	if !rs.Rule.IsRangeDest(rangeValue) {
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
func (rule *RegexRule) IsFullSkip(src string) bool {
	return rule.FullPatterns.IsSkip(src)
}
func (rule *RegexRule) IsFullDest(src string) bool {
	return rule.FullPatterns.IsDest(src)
}
func (rule *RegexRule) IsRangeSkip(src string) bool {
	return rule.RangePatterns.IsSkip(src)
}
func (rule *RegexRule) IsRangeDest(src string) bool {
	return rule.RangePatterns.IsDest(src)
}

func (rule *RegexRule) IsMatchSkip(src string) bool {
	return rule.MatchPatterns.IsSkip(src)
}
func (rule *RegexRule) IsMatchDest(src string) bool {
	return rule.MatchPatterns.IsDest(src)
}

/*
 * is skip check
 */
func (checker *CheckPatternConfig) IsSkip(src string) bool {

	if IfAnyMatch(src, checker.SkipIfany, false) {
		return true
	}
	if WhenMatch(src, checker.SkipWhen, false) {
		return true
	}
	return false
}

/*
 * is dest check
 */
func (checker *CheckPatternConfig) IsDest(src string) bool {

	if IfAnyMatch(src, checker.DoIfany, true) {
		return true
	}
	if WhenMatch(src, checker.DoWhen, true) {
		return true
	}
	return false
}

/*
 * if empty, return emptyFlag
 * if any true,  return true.
 */
func IfAnyMatch(src string, patterns []string, emptyFlag bool) bool {
	if len(patterns) == 0 {
		return emptyFlag
	}
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
 * if empty, return emptyFlag
 * if any false,  return false.
 */
func WhenMatch(src string, patterns []string, emptyFlag bool) bool {
	if len(patterns) == 0 {
		return emptyFlag
	}
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
