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
	config.EncodePattern(&pattern)
	r := regexp.MustCompile(pattern)
	return &Regex{
		R: r,
		Result: RegexResult{
			Pattern:    pattern,
			GroupNames: r.SubexpNames(),
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
	//subMatches := rs.R.FindAllStringSubmatch(input, -1)
	positions := rs.R.FindAllStringSubmatchIndex(input, -1)
	rs.Result.Captures = *SplitBy(&positions, input, false, rs.Result.Captures)

	rs.Result.Params["matches.count"] = strconv.Itoa(len(positions) / 2)
	rs.Result.Params["groups.count"] = strconv.Itoa(len(rs.Result.GroupNames))
	rs.Result.Params["groups.keys"] = strings.Join(rs.Result.GroupNames, ",")

	// match.Index.
	x := 0
	for i, c := range rs.Result.Captures {
		// skip if it's not match
		if !c.IsMatch {
			continue
		}

		match := &rs.Result.Captures[i]
		match.Groups = make([]Capture, 0)
		match.Params = make(map[string]string)
		position := positions[x]

		for x := 0; x < len(rs.Result.GroupNames); x++ {
			gname := rs.Result.GroupNames[x]
			if x == 0 {
				gname = "match.value"
			}
			group := &Capture{Start: position[x*2+0], End: position[x*2+1]}
			group.SetValue(input)
			match.Groups = append(match.Groups, *group)
			match.Params[gname] = group.Value
		}
		x++
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
		log.Printf("Matching file: %s", rs.FromFile)
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
	if rs.ReplaceFlag && rs.HasMatches() {
		if rs.FullCheck(content) {
			replaced, newContent := rs.replaceText(content)
			if replaced {
				config.Decode(&newContent)
				if rs.FromFile != "" {
					rs.writeText(newContent)
				} else {
					log.Println(newContent)
				}
			} else if rs.FromFile != "" {
				log.Printf("Nothing Replaced. File: %s", rs.FromFile)
			} else {
				log.Println("Nothing Replaced.")
				return
			}

		}
	}
	//==========================
	//save to cache
	rs.ToCache()
}

// matched if RegexResult.captures > 1
func (rs *Regex) HasMatches() bool {
	return len(rs.Result.Captures) > 1
}
func (rs *Regex) replaceText(content string) (bool, string) {
	//=replace =============
	var sb strings.Builder
	replaced := false
	for _, c := range rs.Result.Captures {
		if c.IsMatch {
			if rs.IsDestMatch(c, content) {
				mTemplate := NewTemplate(rs.Rule.ReplaceTemplate.Match)
				sb.WriteString(mTemplate.ReplaceByMap(c.Params))
				replaced = true
			} else {
				sb.WriteString(content[c.Start:c.End])
			}
		} else {
			sb.WriteString(content[c.Start:c.End])
		}
	}
	newContent := sb.String()
	//=replace =============
	return replaced, newContent
}

// write content to file
func (rs *Regex) writeText(content string) {
	if rs.ToFile == "" && rs.FromFile != "" {
		rs.ToFile = rs.FromFile
	}
	if rs.ToFile != "" {
		WriteAll(rs.ToFile, content)
		log.Printf("Written To:%s", rs.ToFile)
	}
}

// export Matches using export template defined in rule.yaml
// custom template can include header/match/group/footer
// match/group can export by loop
func (rs *Regex) exportMatches() string {
	var sb strings.Builder

	//---------------------------------------
	// replace Footer: rs.Result.Params
	tHeader := NewTemplate(rs.Rule.ExportTemplate.Header)
	sb.WriteString(tHeader.ReplaceByMap(rs.Result.Params))
	//---------------------------------------
	// replace Loop-Matches: rs.Result.Matches
	for _, item := range rs.Result.Captures {
		if item.IsMatch {
			if rs.Rule.ExportTemplate.Match != "" {
				tContent := NewTemplate(rs.Rule.ExportTemplate.Match)
				tmp := tContent.ReplaceByRegexResult(rs.Result)
				sb.WriteString(tmp)
			} else {
				//when template is empty, export match.value
				sb.WriteString(item.Value)
			}
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
