package cmd

import (
	"log"
	"regexp"

	"golang.org/x/exp/slices"
)

// check whether the content is should be skiped by skip-pattern,
// or is required by dest-pattern
// func (rs *Regex) IsFullOK(content string, rule *CheckRule) bool {
// 	return rs.Rule.FullCheckRule.CanDo(content)
// }
// func (rs *Regex) IsRangeOK(match Capture, content string) bool {
// 	return true
// }
// func (rs *Regex) IsMatchOK(match Capture, content string) bool {
// 	return true
// }

// func (rs *RegexText) MergeRangeStartEnd(content string) *[]Capture {
// 	Rule := rs.RegexRule
// 	ranges := Rule.MergeRangeStartEnd(content)
// }

// check the given match[start,end] whether between the range of matches
// by range-start and range-end
// func (rs *RegexText) isDestMatch(match Capture) bool {
// 	ranges := rs.RegexRule.MergeRangeStartEnd(rs.Content)
// 	inRange := false
// 	rangeValue := ""
// 	for _, r := range *ranges {
// 		if r.Start <= match.Start && match.End <= r.End {
// 			inRange = true
// 			rangeValue = rs.Content[r.Start:r.End]
// 			break
// 		}
// 	}
// 	log.Printf("rangeValue:%s", rangeValue)
// 	return inRange
// }

// func (rule *RegexRule) IsRangeSkip(src string) bool {
// 	return rule.RangePatterns.IsSkip(src)
// }
// func (rule *RegexRule) IsRangeDest(src string) bool {
// 	return rule.RangePatterns.IsDest(src)
// }

// func (rule *RegexRule) IsMatchSkip(src string) bool {
// 	return rule.MatchPatterns.IsSkip(src)
// }
// func (rule *RegexRule) IsMatchDest(src string) bool {
// 	return rule.MatchPatterns.IsDest(src)
// // }

// type ICheck interface {
// 	IsSkip() bool
// 	IsDest() bool
// }

// /*
//  * is skip check
//  */
// func (checker *CheckRule) IsSkip(src string) bool {
// 	if len(checker.Patterns) == 0 {
// 		return checker.Skip
// 	}
// 	if checker.isOK(src) {
// 		return checker.Skip
// 	} else {
// 		return !checker.Skip
// 	}
// }

// check if can do, default is true.
func (checker *CheckRule) CanDo(src string) bool {
	if len(checker.Patterns) == 0 {
		return !checker.Skip
	}
	matched := checker.isOK(src)
	if !matched {
		//not matched do skip
		return checker.Skip
	}

	//matched do not skip
	return !checker.Skip
}

/*
 * is destination check
 */
func (checker *CheckRule) isOK(src string) bool {
	patterns := checker.Patterns
	if len(patterns) == 0 {
		// if condition is empty return true
		return true
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
	if checker.IsOr {
		// if contains any true, return true
		return slices.Contains(boolSlice, true)
	} else {
		// if contains any false, return false
		return slices.Contains(boolSlice, false)
	}
}

// /*
//  * if empty, return emptyFlag
//  * if any true,  return true.
//  */
// func IfAnyMatch(src string, patterns []string, emptyFlag bool) bool {
// 	if len(patterns) == 0 {
// 		return emptyFlag
// 	}
// 	var boolSlice []bool
// 	for _, pattern := range patterns {
// 		if pattern == "" {
// 			continue
// 		}

// 		if ok, err := regexp.MatchString(pattern, src); err == nil {
// 			boolSlice = append(boolSlice, ok)
// 		} else {
// 			log.Fatal(err)
// 		}
// 	}
// 	// if contains true, return true
// 	return slices.Contains(boolSlice, true)
// }

// /*
//  * if empty, return emptyFlag
//  * if any false,  return false.
//  */
// func WhenMatch(src string, patterns []string, emptyFlag bool) bool {
// 	if len(patterns) == 0 {
// 		return emptyFlag
// 	}
// 	var boolSlice []bool
// 	for _, pattern := range patterns {
// 		if pattern == "" {
// 			continue
// 		}

// 		if ok, err := regexp.MatchString(pattern, src); err == nil {
// 			boolSlice = append(boolSlice, ok)
// 		} else {
// 			log.Fatal(err)
// 		}
// 	}
// 	//if empty, return false
// 	if len(boolSlice) == 0 {
// 		return false
// 	}
// 	// if contains false, return false
// 	return !slices.Contains(boolSlice, false)
// }
