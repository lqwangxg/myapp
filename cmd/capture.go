package cmd

type Capture struct {
	Start   int
	End     int
	Value   string
	IsMatch bool
	Params  map[string]string
	Groups  []Capture
}

func (cap *Capture) SetValue(input string) {
	cap.Value = input[cap.Start:cap.End]
}
