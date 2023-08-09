package cmd

import (
	"regexp"
)

func beforeMatch(content string) string {
	ret := content
	for key, val := range config.echars {
		ret = replaceText(ret, key, val)
	}
	return ret
}

func afterMatch(content string) string {
	ret := content
	for key, val := range config.echars {
		ret = replaceText(ret, val, key)
	}
	return ret
}

func replaceText(input, pattern, replacement string) string {
	r, _ := regexp.Compile(pattern)
	return r.ReplaceAllString(input, replacement)
}

// func match(pattern, content string) {
// 	r := regexp.MustCompile(pattern)
// 	groupNames := r.SubexpNames()
// 	matches := r.FindAllStringSubmatch(content, -1)
// 	matchIdxes := r.FindAllStringSubmatchIndex(content, -1)
// 	for i, m := range matchIdxes {
// 		log.Printf("[%d]: %v", i, m)

// 	}

// 	for i, m := range matches {
// 		//slice :=matchIdxes[i]
// 		log.Printf("[%d]: spos:%d %s", i, m)
// 		for _, n := range groupNames {
// 			if n != "" {
// 				lastIndex := r.SubexpIndex(n)
// 				log.Printf("lastIndex=%d, group[%s]: %s", lastIndex, n, m[lastIndex])
// 			}
// 		}
// 	}

// 	// fmt.Println(r.FindAllStringSubmatchIndex(content, -1))
// }
