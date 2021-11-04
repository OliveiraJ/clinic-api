package model

type ExtRule struct {
	Day       CustomDay  `json:"day"`
	Limit     CustomDay  `json:"-"`
	Intervals []Interval `json:"intervals"`
}
