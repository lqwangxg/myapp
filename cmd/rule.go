/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

func (rs *RegexRules) GetDefaultRule() *RegexRule {
	for _, r := range rs.Rules {
		if r.Name == "default" {
			r.ResetTemplate()
			return &r
		}
	}
	return nil
}
func (rs *RegexRules) GetRule(name string) *RegexRule {
	for _, r := range rs.Rules {
		if r.Name == name {
			r.ResetTemplate()
			return &r
		}
	}
	return nil
}
func (rule *RegexRule) ResetTemplate() {
	if rule.ExportTemplate == nil {
		if rule.ExportTemplateName == "" {
			rule.ExportTemplateName = "default"
		}
		for _, t := range appContext.RegexTemplates.Templates {
			if t.Name == rule.ExportTemplateName {
				rule.ExportTemplate = &t
				break
			}
		}
	}
	if rule.ReplaceTemplate == nil {
		if rule.ReplaceTemplateName == "" {
			rule.ReplaceTemplateName = "default"
		}
		for _, t := range appContext.RegexTemplates.Templates {
			if t.Name == rule.ReplaceTemplateName {
				rule.ReplaceTemplate = &t
				break
			}
		}
	}
}

// func (rule *RegexRule) findByName(ruleName string) bool {
// 	return rule.findRule(config.RuleDir, ruleName)
// }
// func (rule *RegexRule) findRule(dirPath string, ruleName string) bool {
// 	files, err := os.ReadDir(dirPath)
// 	if err != nil {
// 		return false
// 	}

// 	r := regexp.MustCompile(ruleName + `\.(yaml|yml)$`)
// 	for _, file := range files {
// 		fullPath := filepath.Join(dirPath, file.Name())
// 		isdir, err := IsDir(fullPath)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if isdir {
// 			if rule.findRule(fullPath, ruleName) {
// 				return true
// 			}
// 			continue
// 		}
// 		//skip ifnot destfile of yaml or yml
// 		if !r.MatchString(file.Name()) {
// 			continue
// 		}
// 		//=============================
// 		//TODO
// 		//return LoadConfig(fullPath, &rule)
// 	}
// 	return false
// }
