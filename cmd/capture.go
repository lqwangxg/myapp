package cmd

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
	for key, fromValue := range fromMatch.Params {
		// If the key not exists, add key/value
		// if exists, do nothing
		if _, ok := rs.Params[key]; !ok {
			rs.Params[key] = fromValue
		}
	}
}
