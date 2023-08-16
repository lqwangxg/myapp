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

func (rs *RegexResult) Export(template *RegexTemplate, matchOnly bool) (string, bool) {
	var sb strings.Builder
	hasChanged := false
	//---------------------------------------
	// replace Footer: rs.Result.Params
	if template != nil && template.Header != "" {
		tHeader := NewTemplate(template.Header)
		header, changed := tHeader.ReplaceByMap(rs.Params)
		sb.WriteString(header)
		if changed {
			hasChanged = true
		}
	}

	//---------------------------------------
	for _, item := range rs.Captures {
		if item.IsMatch {
			if template != nil && template.Match != "" {
				tMatch := NewTemplate(template.Match)
				tmp, changed := tMatch.ReplaceBy(item)
				sb.WriteString(tmp)
				if changed {
					hasChanged = true
				}
				if template.Group != "" {
					//replace by template.group
					for _, g := range item.Groups {
						tGroup := NewTemplate(template.Group)
						tmp, changed := tGroup.ReplaceBy(g)
						sb.WriteString(tmp)
						if changed {
							hasChanged = true
						}
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
	if template != nil && template.Footer != "" {
		tFooter := NewTemplate(template.Footer)
		footer, changed := tFooter.ReplaceByMap(rs.Params)
		sb.WriteString(footer)
		if changed {
			hasChanged = true
		}
	}
	//---------------------------------------
	return sb.String(), hasChanged
}
