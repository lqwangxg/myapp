package cmd

type Capture struct {
	Start   int
	End     int
	Value   string
	IsMatch bool
	Params  map[string]string
	Groups  []Capture
}

// func NewCapture(Start, End int) *Capture {
// 	capture := &Capture{
// 		Start:   Start,
// 		End:     End,
// 		IsMatch: false,
// 	}
// 	return capture
// }

// func NewMatchCapture(Start, End int) *Capture {
// 	capture := &Capture{
// 		Start:   Start,
// 		End:     End,
// 		IsMatch: true,
// 		Params:  make(map[string]string),
// 		Groups:  make([]Capture, 0),
// 	}
// 	return capture
// }

func (cap *Capture) SetValue(input string) {
	cap.Value = input[cap.Start:cap.End]
}
