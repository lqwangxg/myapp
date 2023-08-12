package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// var regex regexp.Regexp
func NewRegex(pattern string) *Regex {
	return NewCacheRegex(pattern, config.RedisOption.Enable)
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
	config.Transform(&input)
	r := rs.R
	if !r.MatchString(input) {
		return
	}
	subMatches := r.FindAllStringSubmatch(input, -1)
	positions := r.FindAllStringSubmatchIndex(input, -1)

	for i, smatch := range subMatches {
		position := positions[i]
		groups := make([]RegexGroup, 0)
		match := RegexMatch{
			Index: i,
			Position: RegexMatchIndex{
				Start: position[0],
				End:   position[1],
			},
			Groups: groups,
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
				Index: x,
				Name:  gname,
				Value: smatch[groupIndex],
				Position: RegexMatchIndex{
					Start: position[x*2+0],
					End:   position[x*2+1],
				},
			}
			match.Params[gname] = group.Value
			groups = append(groups, group)
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
	if flags.IncludeSuffix != "" {
		re := NewRegex(flags.IncludeSuffix)
		if !re.IsMatch(filePath) {
			return
		}
	}
	if flags.ExcludeSuffix != "" {
		re := NewRegex(flags.ExcludeSuffix)
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
	rs.ScanMatches(content)
	export := rs.ExportMatches(flags.TempleteOfExport)
	//=export log =============
	if value, ok := rs.Result.Params["filePath"]; ok {
		log.Printf("file: %s", value)
	}
	if export != "" {
		log.Printf("%s", export)
	} else {
		log.Printf("no matches")
	}
	//=replace log =============
	if rs.Action == ReplaceAction {
		newContent := rs.replaceText()
		fmt.Println(newContent)
	}
	//==========================
	rs.Close()
	return export
}

func (rs *Regex) replaceText() string {
	//=replace log =============
	var sb strings.Builder
	for _, m := range rs.Result.Ranges {
		if m.IsMatch && flags.TemplateOfReplace != "" {
			mval := flags.TemplateOfReplace
			ReplaceTemplate(&mval, rs.Result.Params)
			ReplaceTemplate(&mval, rs.Result.Matches[m.MatchIndex].Params)
			sb.WriteString(mval)
		} else {
			sb.WriteString(m.Value)
		}
	}
	newContent := sb.String()
	config.Restore(&newContent)
	//=replace log =============
	return newContent
}
