package cmd

import (
	"strings"
)

type RegexResult struct {
	Pattern    string
	GroupNames []string
	Captures   []Capture
	Params     map[string]string
	MatchCount int
}

func (rs *RegexResult) Export(template *RegexTemplate, matchOnly bool) string {
	var sb strings.Builder
	//---------------------------------------
	// replace Footer: rs.Result.Params
	if template != nil {
		tHeader := NewTemplate(template.Header)
		sb.WriteString(tHeader.ReplaceByMap(rs.Params))
	}

	//---------------------------------------
	for _, item := range rs.Captures {
		if item.IsMatch {
			if template != nil && template.Match != "" {
				tMatch := NewTemplate(template.Match)
				tmp := tMatch.ReplaceBy(item)
				sb.WriteString(tmp)
				if template.Group != "" {
					//replace by template.group
					for _, g := range item.Groups {
						tGroup := NewTemplate(template.Group)
						tmp := tGroup.ReplaceBy(g)
						sb.WriteString(tmp)
					}
				}
			} else {
				//when template is empty, export match.value
				sb.WriteString(item.Value)
			}
		} else if !matchOnly {
			// export not only match
			sb.WriteString(item.Value)
		}
	}
	//---------------------------------------
	if template != nil {
		tFooter := NewTemplate(template.Footer)
		sb.WriteString(tFooter.ReplaceByMap(rs.Params))
	}
	//---------------------------------------
	return sb.String()
}
