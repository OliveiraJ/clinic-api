package model

import (
	"fmt"
	"log"
	"strings"
	"time"
)

const DAY string = "02-01-2006"
const HOUR string = "15:04"

type CustomDay time.Time
type CustomHour time.Time

type Interval struct {
	Start CustomHour `json:"start"`
	End   CustomHour `json:"end"`
}

type Rule struct {
	Day       CustomDay  `json:"day"`
	Limit     CustomDay  `json:"limit"`
	Intervals []Interval `json:"intervals"`
}

func (cd CustomDay) MarshalJSON() ([]byte, error) {
	return []byte(cd.String()), nil
}

func (cd *CustomDay) String() string {
	t := time.Time(*cd)
	return fmt.Sprintf("%q", t.Format(DAY))
}

func (cd *CustomDay) UnmarshalJSON(dat []byte) error {
	s := strings.Trim(string(dat), `"`)
	day, err := time.Parse(DAY, s)
	if err != nil {
		log.Fatal(err)
	}
	*cd = CustomDay(day)
	return nil
}

func (ch CustomHour) MarshalJSON() ([]byte, error) {
	return []byte(ch.String()), nil
}

func (ch *CustomHour) String() string {
	t := time.Time(*ch)
	return fmt.Sprintf("%q", t.Format(HOUR))
}

func (ch *CustomHour) UnmarshalJSON(dat []byte) error {
	h := strings.Trim(string(dat), `"`)
	hour, err := time.Parse(HOUR, h)
	if err != nil {
		log.Fatal(err)
	}

	*ch = CustomHour(hour)
	return nil
}
