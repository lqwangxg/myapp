package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func NewRegexFromCmd() *Regex {
	var rs *Regex
	if flags.RuleName != "" {
		cr := customRules.GetRule(flags.RuleName)
		if cr != nil {
			rs = cr.NewRegex()
		} else {
			var rule RuleConfig
			if rule.findByName(flags.RuleName) {
				rs = rule.NewRegex()
			}
		}
	} else if flags.Pattern != "" {
		rs = NewRegex(flags.Pattern)
	}
	if rs == nil {
		panic("pattern is empty, ruleName is empty, or ruleName not avaiable.")
	}
	return rs
}

// var regex regexp.Regexp
func NewRegex(pattern string) *Regex {
	return NewCacheRegex(pattern, config.RedisOption.Enable, nil)
}

// var regex regexp.Regexp
func (rule RuleConfig) NewRegex() *Regex {
	return NewCacheRegex(rule.Pattern, config.RedisOption.Enable, &rule)
}
func NewRegexByPattern(pattern string) *Regex {
	return NewCacheRegex(pattern, false, nil)
}

// var regex regexp.Regexp
func NewCacheRegex(pattern string, cache bool, rule *RuleConfig) *Regex {
	r := regexp.MustCompile(pattern)
	return &Regex{
		R: r,
		Result: RegexResult{
			Pattern:    pattern,
			GroupNames: r.SubexpNames(),
			Matches:    make([]RegexMatch, 0),
			Params:     make(map[string]string),
		},
		Cache: cache,
		Rule:  rule,
	}
}

// get hashsum of (pattern + ": "+ input) to string.
func (rs *Regex) Hashsum(input string) string {
	h := sha1.New()
	h.Write([]byte(rs.Result.Pattern + ": " + input))
	return hex.EncodeToString(h.Sum(nil))
}
func (rs *Regex) IsMatch(input string) bool {
	return rs.R.MatchString(input)
}

// var regex regexp.Regexp
func (rs *Regex) ScanMatches(input string) {
	if !rs.IsMatch(input) {
		return
	}
	subMatches := rs.R.FindAllStringSubmatch(input, -1)
	positions := rs.R.FindAllStringSubmatchIndex(input, -1)
	rs.Result.Params["matches.count"] = strconv.Itoa(len(subMatches))
	rs.Result.Params["groups.count"] = strconv.Itoa(len(rs.Result.GroupNames))
	rs.Result.Params["groups.keys"] = strings.Join(rs.Result.GroupNames, ",")

	for i, smatch := range subMatches {
		position := positions[i]
		match := &RegexMatch{
			Capture: &Capture{
				Start: position[0],
				End:   position[1],
				Value: smatch[0],
				RType: MatchType,
			},
			Groups: make([]RegexGroup, 0),
			Params: make(map[string]string),
		}
		//match.Params["match.value"] = match.Value
		for x := 0; x < len(rs.Result.GroupNames); x++ {
			gname := rs.Result.GroupNames[x]
			if x == 0 {
				gname = "match.value"
			}
			group := &RegexGroup{
				Name: gname,
				Capture: &Capture{
					Start: position[x*2+0],
					End:   position[x*2+1],
					Value: smatch[x],
					RType: GroupType,
				},
			}
			match.Groups = append(match.Groups, *group)
			match.Params[group.Name] = group.Value
		}
		rs.Result.Matches = append(rs.Result.Matches, *match)
	}
}

// split input by matches
func (rs *Regex) SplitMatches(input string) {
	cpos := 0
	epos := len(input)
	rs.Result.Ranges = rs.Result.Ranges[:cap(rs.Result.Ranges)]
	for i, match := range rs.Result.Matches {
		if cpos < epos && cpos < match.Start {
			//append string before match
			h := &RegexRange{
				Capture: &Capture{
					Start: cpos,
					End:   match.Start,
					Value: input[cpos:match.Start],
					RType: UnMatchType,
				},
			}
			rs.Result.Ranges = append(rs.Result.Ranges, *h)
		}
		// append match.value
		m := &RegexRange{
			Capture: &Capture{
				Start: match.Start,
				End:   match.End,
				Value: input[match.Start:match.End],
				RType: MatchType,
			},
			MatchIndex: i,
		}
		rs.Result.Ranges = append(rs.Result.Ranges, *m)
		cpos = match.End
	}
	if cpos < epos {
		//append last string
		f := &RegexRange{
			Capture: &Capture{
				Start: cpos,
				End:   epos,
				Value: input[cpos:epos],
				RType: UnMatchType,
			},
		}
		rs.Result.Ranges = append(rs.Result.Ranges, *f)
	}
}

func (rs *Regex) ProcFile(filePath string) {
	// exit if file not exists
	if !IsExists(filePath) {
		return
	}
	if flags.IncludeSuffix != "" {
		if !IsMatchString(flags.IncludeSuffix, filePath) {
			return
		}
	}
	if rs.Rule.IncludeFile != "" {
		if !IsMatchString(rs.Rule.IncludeFile, filePath) {
			return
		}
	}
	if flags.ExcludeSuffix != "" {
		if IsMatchString(flags.ExcludeSuffix, filePath) {
			return
		}
	}
	if rs.Rule.ExcludeFile != "" {
		if IsMatchString(rs.Rule.ExcludeFile, filePath) {
			return
		}
	}

	if buffer, err := ReadAll(filePath); err == nil {
		rs.FromFile = filePath
		rs.MatchText(buffer)
	}
}

func (rs *Regex) ProcDir(dirPath string) {
	if !IsExists(dirPath) {
		log.Printf("dirPath is not found. dirPath=%s", dirPath)
		return
	}
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, file := range files {
		fullPath := filepath.Join(dirPath, file.Name())
		ok, err := IsDir(fullPath)
		if err == nil {
			if ok {
				rs.ProcDir(fullPath)
			} else {
				rs.ProcFile(fullPath)
			}
		}
	}
}

func (rs *Regex) MatchText(content string) {
	//before match input, transfer special chars
	config.Encode(&content)
	// get from redis cache
	if !rs.FromCache(content) {
		rs.ScanMatches(content)
	}

	//=export log =============
	if rs.ExportFlag {
		export := rs.exportMatches()
		config.Decode(&export)
		if export != "" {
			log.Printf("%s", export)
		} else {
			log.Printf("no matches")
		}
	}
	//=replace log =============
	if rs.ReplaceFlag {
		if rs.FullCheck(content) {
			newContent := rs.replaceText(content)
			config.Decode(&newContent)
			if rs.FromFile != "" {
				rs.writeText(newContent)
			} else {
				log.Println(newContent)
			}
		}
	}
	//==========================
	//save to cache
	rs.ToCache()
}

func (rs *Regex) replaceText(content string) string {
	//=replace log =============
	var sb strings.Builder
	// template := rs.getTemplate()
	// split input by matches
	rs.SplitMatches(content)
	for _, m := range rs.Result.Ranges {
		if m.RType == MatchType {
			match := rs.Result.Matches[m.MatchIndex]
			if rs.IsDestMatch(match, content) {
				mTemplate := NewTemplate(rs.Rule.ReplaceTemplate.Match)
				sb.WriteString(mTemplate.ReplaceByMap(match.Params))
			} else {
				sb.WriteString(m.Value)
			}
		} else {
			sb.WriteString(m.Value)
		}
	}
	newContent := sb.String()
	//=replace log =============
	return newContent
}
func (rs *Regex) writeText(content string) {
	if rs.FromFile != "" {
		WriteAll(rs.FromFile, content)
	}
}
func (rs *Regex) exportMatches() string {
	var sb strings.Builder

	//---------------------------------------
	// replace Footer: rs.Result.Params
	tHeader := NewTemplate(rs.Rule.ExportTemplate.Header)
	sb.WriteString(tHeader.ReplaceByMap(rs.Result.Params))
	//---------------------------------------
	// replace Loop-Matches: rs.Result.Matches
	for i := 0; i < len(rs.Result.Matches); i++ {
		if rs.Rule.ExportTemplate.Match != "" {
			tContent := NewTemplate(rs.Rule.ExportTemplate.Match)
			tmp := tContent.ReplaceByRegexResult(rs.Result)
			sb.WriteString(tmp)
		} else {
			//when template is empty, export match.value
			sb.WriteString(rs.Result.Matches[i].Value)
		}
	}
	//---------------------------------------
	// replace Footer: rs.Result.Params
	tFooter := NewTemplate(rs.Rule.ExportTemplate.Footer)
	sb.WriteString(tFooter.ReplaceByMap(rs.Result.Params))
	//---------------------------------------

	exports := sb.String()
	return exports
}

// func (rs *Regex) replaceMatch(index int, template string) string {
// 	var sb strings.Builder

// 	mTemplate := NewTemplate(template)
// 	mTemplate.ReplaceByMap(rs.Result.Params)
// 	mTemplate.ReplaceByMap(rs.Result.Matches[index].Params)
// 	sb.WriteString(mTemplate.Template)

// 	buffer := sb.String()
// 	config.Decode(&buffer)
// 	return buffer
// }

// func (rs *Regex) close() {
// 	// //match restore
// 	// for i := 0; i < len(rs.Result.Matches); i++ {
// 	// 	config.Decode(&rs.Result.Matches[i].Value)
// 	// 	for x := 0; x < len(rs.Result.Matches[i].Groups); x++ {
// 	// 		config.Decode(&rs.Result.Matches[i].Groups[x].Value)
// 	// 	}
// 	// }
// 	// //range restore
// 	// for x := 0; x < len(rs.Result.Ranges); x++ {
// 	// 	config.Decode(&rs.Result.Ranges[x].Value)
// 	// }
// 	//save to cache
// 	rs.ToCache()
// }
