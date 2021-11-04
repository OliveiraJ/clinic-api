package model

type Interval struct {
	Start CustomHour `json:"start"`
	End   CustomHour `json:"end"`
}
