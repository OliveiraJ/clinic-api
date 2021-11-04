package utils

import (
	"log"
	"time"

	"github.com/OliveiraJ/clinic-api/model"
)

const PATH string = "./database/rules.json"

// Checks if there are any existing rules, preventing duplicate rules from being added as those are removed from the slice rule.Intervals.
func CheckInvalidSchedule(x model.Rule, rule model.Rule) bool {
	for _, interval := range x.Intervals {
		for _, ruleInterval := range rule.Intervals {
			if interval == ruleInterval {
				log.Println("Schedule already reserved")
				return true
			}
		}
	}

	return false
}

func CheckInvalidInterval(rule model.Rule) bool {
	for _, interval := range rule.Intervals {
		if time.Time(interval.End).Before(time.Time(interval.Start)) || time.Time(interval.End).Equal(time.Time(interval.Start)) {
			return true
		}
	}

	return false
}

func CheckInvalidDate(rule model.Rule) bool {
	if time.Time(rule.Limit).Equal(time.Time(rule.Day)) {
		return false
	} else {
		return true
	}
}

func CheckOverlapingIntervals(x model.Rule, rule model.Rule) bool {
	check := false
	for _, new := range rule.Intervals {
		for _, interval := range x.Intervals {
			if (time.Time(new.Start).Before(time.Time(interval.Start)) && (time.Time(new.End).Before(time.Time(interval.Start)) || time.Time(new.End).Equal(time.Time(interval.Start)))) || ((time.Time(new.Start).After(time.Time(interval.End)) || time.Time(new.Start).Equal(time.Time(interval.End))) && (time.Time(new.End).After(time.Time(interval.End)))) {
				check = false
			} else {
				check = true
			}
		}
	}

	return check
}
