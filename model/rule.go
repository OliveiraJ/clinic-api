package model

type Rule struct {
	Day       CustomDay  `json:"day"`
	Limit     CustomDay  `json:"limit"`
	Intervals []Interval `json:"intervals"`
}
