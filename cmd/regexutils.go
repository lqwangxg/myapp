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
		rule, found := localRules[flags.RuleName]
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
	return NewCacheRegex(pattern, config.RedisOption.Enable)
}

// var regex regexp.Regexp
func (rule RuleConfig) NewRegex() *Regex {
	rs := NewCacheRegex(rule.Pattern, config.RedisOption.Enable)
	rs.Rule = &rule
	return rs
}
func NewNoCacheRegex(pattern string) *Regex {
	return NewCacheRegex(pattern, false)
}

// var regex regexp.Regexp
func NewCacheRegex(pattern string, cache bool) *Regex {
	r := regexp.MustCompile(pattern)
	return &Regex{
		R: r,
		Result: RegexResult{
			Pattern:    pattern,
			GroupNames: r.SubexpNames(),
			Params:     make(map[string]string),
		},
		Cache: cache,
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
	r := rs.R
	if !r.MatchString(input) {
		return
	}
	subMatches := r.FindAllStringSubmatch(input, -1)
	positions := r.FindAllStringSubmatchIndex(input, -1)
	rs.Result.Params["matches.count"] = strconv.Itoa(len(subMatches))
	rs.Result.Params["groups.count"] = strconv.Itoa(len(rs.Result.GroupNames))
	rs.Result.Params["groups.keys"] = strings.Join(rs.Result.GroupNames, ",")

	for i, smatch := range subMatches {
		position := positions[i]
		match := RegexMatch{
			Index: i,
			Position: RegexMatchIndex{
				Start: position[0],
				End:   position[1],
			},
			Groups: make([]RegexGroup, 0),
			Value:  smatch[0],
			Params: make(map[string]string),
		}

		for x, groupName := range rs.Result.GroupNames {
			//from match 0 ~ count - 1
			groupIndex := r.SubexpIndex(groupName)
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
				Position: RegexMatchIndex{
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
func (rs *Regex) SplitMatches(input string) []RegexRange {
	cpos := 0
	epos := len(input)
	rs.Result.Ranges = rs.Result.Ranges[:cap(rs.Result.Ranges)]
	for _, match := range rs.Result.Matches {
		if cpos < epos && cpos < match.Position.Start {
			//append string before match
			h := &RegexRange{
				Value:   input[cpos:match.Position.Start],
				IsMatch: false,
			}
			rs.Result.Ranges = append(rs.Result.Ranges, *h)
		}
		// append match.value
		m := &RegexRange{
			Value:      input[match.Position.Start:match.Position.End],
			IsMatch:    true,
			MatchIndex: match.Index,
		}
		rs.Result.Ranges = append(rs.Result.Ranges, *m)
		cpos = match.Position.End
	}
	if cpos < epos {
		//append last string
		f := &RegexRange{
			Value:   input[cpos:epos],
			IsMatch: false,
		}
		rs.Result.Ranges = append(rs.Result.Ranges, *f)
	}
	return rs.Result.Ranges
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
		rs.Result.Params["filePath"] = filePath
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
	filePath, hasFilePath := rs.Result.Params["filePath"]
	rs.Content = content
	if hasFilePath {
		log.Printf("file: %s", filePath)
	}

	rs.ScanMatches(rs.Content)
	export := rs.exportMatches()
	//=export log =============
	if export != "" {
		log.Printf("%s", export)
	} else {
		log.Printf("no matches")
	}
	//=replace log =============
	if rs.IsReplace && rs.fullCheck(rs.Content) {
		newContent := rs.replaceText("")
		if hasFilePath {
			WriteAll(filePath, newContent)
		} else {
			log.Println(newContent)
		}
	}
	//==========================
	rs.close()
	return export
}

func (rs *Regex) replaceText(template string) string {
	//=replace log =============
	var sb strings.Builder
	if template == "" {
		template = rs.Rule.ReplaceTemplate
	}
	if template == "" && flags.ReplaceTemplate != "" {
		template = flags.ReplaceTemplate
	}
	for _, m := range rs.Result.Ranges {
		if m.IsMatch {
			match := rs.Result.Matches[m.MatchIndex]
			if rs.matchCheck(rs.Content, match) {
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
	ReplaceByMap(&mval, rs.Result.Matches[index].Params)
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
