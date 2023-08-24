package cmd

import "strconv"

type Capture struct {
	Start   int
	End     int
	Value   string
	IsMatch bool
	Params  map[string]string
	Groups  []Capture
}

func (rs *Capture) MergeParams(fromMatch *Capture) {
	if !fromMatch.IsMatch {
		return
	}
	//not overwrite
	MergeMap(fromMatch.Params, rs.Params, false)
}

func MergeMap(fromParams, toParams map[string]string, overwrite bool) {
	for key, fromValue := range fromParams {
		if _, ok := toParams[key]; ok {
			// if exists, can overwrite.
			if overwrite {
				toParams[key] = fromValue
			}
		} else {
			// if not exists.
			toParams[key] = fromValue
		}
	}
}
func (rs *Capture) Skip() bool {
	if !rs.IsMatch {
		return true
	}
	if val, ok := rs.Params["match.skip"]; ok {
		if skip, err := strconv.ParseBool(val); err == nil {
			return skip
		}
	}
	return false
}
