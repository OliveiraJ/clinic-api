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

// A auxiliar struct to help filtering only the needed data to be sent to the client.
type ExtRule struct {
	Day       CustomDay  `json:"day"`
	Limit     CustomDay  `json:"-"`
	Intervals []Interval `json:"intervals"`
}

// Implements the MarshalJSON interface to the CustomDay type
func (cd CustomDay) MarshalJSON() ([]byte, error) {
	return []byte(cd.String()), nil
}

// Implements the String interface to the CustomDay type, allowing it to me encoded and decoded with a custom format defined by the DAY const.
func (cd *CustomDay) String() string {
	t := time.Time(*cd)
	return fmt.Sprintf("%q", t.Format(DAY))
}

// Implements the UnmarshalJSON interface to the CustomDay type, allowing it to be decoded to the tipe time.Time with a custom formta defined by the DAY const
func (cd *CustomDay) UnmarshalJSON(dat []byte) error {
	s := strings.Trim(string(dat), `"`)
	day, err := time.Parse(DAY, s)
	if err != nil {
		log.Fatal(err)
	}
	*cd = CustomDay(day)
	return nil
}

// Implements the MarshalJSON interface to the CustomHour type
func (ch CustomHour) MarshalJSON() ([]byte, error) {
	return []byte(ch.String()), nil
}

// Implements the String interface to the CustomHour type, allowing it to me encoded and decoded with a custom format defined by the HOUR const
func (ch *CustomHour) String() string {
	t := time.Time(*ch)
	return fmt.Sprintf("%q", t.Format(HOUR))
}

// Implements the UnmarshalJSON interface to the CustomHour type, allowing it to be decoded to the tipe time.Time with a custom formta defined by the HOUR const
func (ch *CustomHour) UnmarshalJSON(dat []byte) error {
	h := strings.Trim(string(dat), `"`)
	hour, err := time.Parse(HOUR, h)
	if err != nil {
		log.Fatal(err)
	}

	*ch = CustomHour(hour)
	return nil
}
