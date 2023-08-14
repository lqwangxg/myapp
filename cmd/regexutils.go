package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func NewRegexFromCmd() *Regex {
	var rs *Regex
	if flags.RuleName != "" {
		found, rule := LoadRule(flags.RuleName)
		if found {
			rs = rule.NewRegex()
			return rs
		}
	}
	if flags.Pattern != "" {
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
func NewNoCacheRegex(pattern string) *Regex {
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
			Ranges:     make([]RegexRange, 0),
			Params:     make(map[string]any),
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
	// get from redis cache
	if rs.FromCache(input) {
		return
	}

	//before match input, transfer special chars
	config.Encode(&input)
	if !rs.IsMatch(input) {
		return
	}
	subMatches := rs.R.FindAllStringSubmatch(input, -1)
	positions := rs.R.FindAllStringSubmatchIndex(input, -1)
	rs.Result.Params["matches.count"] = len(subMatches)
	rs.Result.Params["groups.count"] = len(rs.Result.GroupNames)
	rs.Result.Params["groups.keys"] = rs.Result.GroupNames

	for i, smatch := range subMatches {
		position := positions[i]
		match := RegexMatch{
			Index: i,
			Bound: Bound{
				Start: position[0],
				End:   position[1],
			},
			Groups: make([]RegexGroup, 0),
			Value:  smatch[0],
			Params: make(map[string]string),
		}

		for x, groupName := range rs.Result.GroupNames {
			//from match 0 ~ count - 1
			groupIndex := rs.R.SubexpIndex(groupName)
			if groupName == "" {
				groupIndex = x
			}
			gname := groupName
			if x == 0 {
				gname = "match.value"
			}
			group := RegexGroup{
				Name:  gname,
				Value: smatch[groupIndex],
				Bound: Bound{
					Start: position[x*2+0],
					End:   position[x*2+1],
				},
			}
			match.Params[gname] = config.Decode(&group.Value)
			match.Groups = append(match.Groups, group)
		}
		rs.Result.Matches = append(rs.Result.Matches, match)
	}
	// split input by matches
	rs.SplitMatches(input)
}

// split input by matches
func (rs *Regex) SplitMatches(input string) {
	cpos := 0
	epos := len(input)
	rs.Result.Ranges = rs.Result.Ranges[:cap(rs.Result.Ranges)]
	for _, match := range rs.Result.Matches {
		if cpos < epos && cpos < match.Bound.Start {
			//append string before match
			h := &RegexRange{
				Value:   input[cpos:match.Bound.Start],
				IsMatch: false,
			}
			rs.Result.Ranges = append(rs.Result.Ranges, *h)
		}
		// append match.value
		m := &RegexRange{
			Value:      input[match.Bound.Start:match.Bound.End],
			IsMatch:    true,
			MatchIndex: match.Index,
		}
		rs.Result.Ranges = append(rs.Result.Ranges, *m)
		cpos = match.Bound.End
	}
	if cpos < epos {
		//append last string
		f := &RegexRange{
			Value:   input[cpos:epos],
			IsMatch: false,
		}
		rs.Result.Ranges = append(rs.Result.Ranges, *f)
	}
}

func (rs *Regex) ProcFile(filePath string) {
	// exit if file not exists
	if !IsExists(filePath) {
		return
	}
	if flags.IncludeSuffix != "" || rs.Rule.IncludeSuffix != "" {
		re := NewRegex(flags.IncludeSuffix)
		if !re.IsMatch(filePath) {
			return
		}
		re = NewRegex(rs.Rule.IncludeSuffix)
		if !re.IsMatch(filePath) {
			return
		}
	}
	if flags.ExcludeSuffix != "" || rs.Rule.ExcludeSuffix != "" {
		re := NewRegex(flags.ExcludeSuffix)
		if re.IsMatch(filePath) {
			return
		}
		re = NewRegex(rs.Rule.ExcludeSuffix)
		if re.IsMatch(filePath) {
			return
		}
	}
	if buffer, err := ReadAll(filePath); err == nil {
		rs.fromFile = filePath
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

func (rs *Regex) MatchText(content string) string {
	rs.content = content
	rs.ScanMatches(rs.content)

	//=export log =============
	export := rs.exportMatches()
	if export != "" {
		log.Printf("%s", export)
	} else {
		log.Printf("no matches")
	}
	//=replace log =============
	if rs.IsReplace && rs.FullCheck(rs.content) {
		newContent := rs.replaceText()
		if rs.fromFile != "" {
			rs.writeText(newContent)
		} else {
			log.Println(newContent)
		}
	}
	//==========================
	rs.close()
	return export
}
func (rs *Regex) getTemplate() string {
	// private template first.
	if rs.template != "" {
		return rs.template
	}
	// then flags.ReplaceTemplate
	if flags.ReplaceTemplate != "" {
		return flags.ReplaceTemplate
	}
	// then rule.ReplaceTemplate
	if rs.Rule.ReplaceTemplate != "" {
		return rs.Rule.ReplaceTemplate
	}
	return ""
}
func (rs *Regex) replaceText() string {
	//=replace log =============
	var sb strings.Builder
	template := rs.getTemplate()

	for _, m := range rs.Result.Ranges {
		if m.IsMatch && template != "" {
			match := rs.Result.Matches[m.MatchIndex]
			//
			if rs.IsDestMatch(rs.content, match) {
				sb.WriteString(rs.replaceMatch(m.MatchIndex, template))
			} else {
				sb.WriteString(m.Value)
			}
		} else {
			sb.WriteString(m.Value)
		}
	}
	newContent := sb.String()
	config.Decode(&newContent)
	//=replace log =============
	return newContent
}
func (rs *Regex) writeText(content string) {
	if rs.fromFile != "" {
		WriteAll(rs.fromFile, content)
	}
}
func (rs *Regex) exportMatches() string {
	template := rs.Rule.ExportTemplate
	if flags.ExportTemplate != "" {
		template = flags.ExportTemplate
	}
	var sb strings.Builder
	for i := 0; i < len(rs.Result.Matches); i++ {
		if template != "" {
			tmp := rs.replaceMatch(i, template)
			rs.ReplaceLoop(&tmp, ReplaceByMap)
			sb.WriteString(tmp)
		} else {
			//when template is empty, export match.value
			sb.WriteString(rs.Result.Matches[i].Value)
		}
	}
	exports := sb.String()
	config.Decode(&exports)
	return exports
}

func (rs *Regex) replaceMatch(index int, template string) string {
	var sb strings.Builder
	mval := template
	ReplaceByMap(&mval, rs.Result.Params)
	//ReplaceByMap(&mval, rs.Result.Matches[index].Params)
	sb.WriteString(mval)

	buffer := sb.String()
	config.Decode(&buffer)
	return buffer
}

func (rs *Regex) close() {
	//match restore
	for i := 0; i < len(rs.Result.Matches); i++ {
		config.Decode(&rs.Result.Matches[i].Value)
		for x := 0; x < len(rs.Result.Matches[i].Groups); x++ {
			config.Decode(&rs.Result.Matches[i].Groups[x].Value)
		}
	}
	//range restore
	for x := 0; x < len(rs.Result.Ranges); x++ {
		config.Decode(&rs.Result.Ranges[x].Value)
	}
	//save to cache
	rs.ToCache()
}
