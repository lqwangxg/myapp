#### rule-name
```yaml
name: button-size
```

#### rule-group
```yaml
group: html-input
```
#### before replace process: check full content skip or do.
#### priority descend: skip_ifany, skip_when, do_ifany, do_when
#### if any true return skip/do;
#### optional: if groupname is included in pattern, 
####  { match.group.name: match.group.nalue } will be saved in ${full_params_map}
```yaml
full_patterns :
  skip_ifany: [""]
  skip_when: [""]
  do_ifany: [""]
  do_when: [""]
```

#### match destination range(spos=startIndex, epos=endIndex + endLength) 
#### if range_start is empty then spos=0
#### if range_end is empty then epos=length(full)
#### optional: if groupname is included in pattern, 
####  { match.group.name: match.group.nalue } will be saved in ${match_params_map}
```yaml
range_pattern: \<html:button.*?\/?\>
```

#### before replace process: get params from range_match
####  { match.group.name: match.group.nalue } will be saved in ${match_params_map}
```yaml
range_params_pattern:
  - 'text=[\s"]*(?<text>.+?)[\s"]*'
  - 'width:[\s"]*(?<width>.+?)px'
  - 'maxlength=[\s"]*(?<maxlength>\d+?)[\s"]*'
  - 'size=[\s"]*(?<size>\d+?)[\s"]*'
```

#### before replace process: check full content skip or do.
#### priority descend: skip_ifany, skip_when, do_ifany, do_when
#### if any true return skip/do;
#### optional: if groupname is included in pattern, 
####  { match.group.name: match.group.nalue } will be saved in ${match_params_map}
```yaml
range_patterns:
  skip_ifany: [""]
  skip_when: [""]
  do_ifany: [""]
  do_when: [""]
```
#### before replace process: get replace destination match
####  { match.group.name: match.group.nalue } will be saved in ${match_params_map}
```yaml
match_pattern: width:[\s"]*(?<width>\d+)
```
#### before replace process: check skip or do replace conditions
```yaml
match_patterns:
  skip_ifany: [""]
  skip_when: [""]
  do_ifany: [""]
  do_when: [""]
```

#### process for refreshing params
#### support functions: 
####  string: len(string)
####  number: +,-,*,/
####  compare: <, <=, >, >=, ==
####  reference: [maja42/goval](https://github.com/maja42/goval)
```yaml
match_formulas:
  - text_width = len(text) * 10 + 40

match_evals:
  skip_ifany: [""]
  skip_when: [""]
  do_ifany: [""]
  do_when: [""]
```
#### replaceByEvalFunc(match_pattern, match_replace)
```yaml
match_replace: 'width: ${text_width}'
```