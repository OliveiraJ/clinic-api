package utils

import (
	"log"
	"time"

	"github.com/OliveiraJ/clinic-api/model"
)

const PATH string = "./database/rules.json"

// Checks if there are any existing rules, preventing duplicate rules from being added as those are removed from the slice rule.Intervals.
func CheckInvalidSchedule(comparator model.Rule, newRule model.Rule) bool {
	for _, interval := range comparator.Intervals {
		for _, ruleInterval := range newRule.Intervals {
			if interval == ruleInterval {
				log.Println("Schedule already reserved")
				return true
			}
		}
	}

	return false
}

// Checks if the interval btween end and start are a valid interval.
// Ex.:
// { start: 14:00, end: 14:30 }
// Is a valid interval while,
// { start: 14:00, end: 14:00 } and { start: 14:00, end: 13:00}
// Are exemples of invalid intervals.
// The function returns a bool, "true" if the interval is a Invalid, otherwise "false" is returned
func CheckInvalidInterval(newRule model.Rule) bool {
	for _, interval := range newRule.Intervals {
		if time.Time(interval.End).Before(time.Time(interval.Start)) || time.Time(interval.End).Equal(time.Time(interval.Start)) {
			return true
		}
	}

	return false
}

// Checks if the date interval is valid, in case of 'limit' and 'day' receive equal value or 'limit' preceds 'day' true is returned,
// otherwise false is returned.
func CheckInvalidDate(newRule model.Rule) bool {
	if time.Time(newRule.Limit).Equal(time.Time(newRule.Day)) {
		return true
	} else if time.Time(newRule.Limit).Before(time.Time(newRule.Day)) {
		return true
	} else {
		return false
	}
}

// Checks for overlaping intervals, preventing the insertion of rules that overlap each other.
// Returns a bool, true if there are overlaping schedule times in the rule or false if there isnt.
func CheckOverlappingIntervals(comparator model.Rule, newRule model.Rule) bool {
	check := false
	for _, new := range newRule.Intervals {
		for _, interval := range comparator.Intervals {
			if (time.Time(new.Start).Before(time.Time(interval.Start)) && (time.Time(new.End).Before(time.Time(interval.Start)) || time.Time(new.End).Equal(time.Time(interval.Start)))) || ((time.Time(new.Start).After(time.Time(interval.End)) || time.Time(new.Start).Equal(time.Time(interval.End))) && (time.Time(new.End).After(time.Time(interval.End)))) {
				check = false
			} else {
				check = true
			}
		}
	}

	return check
}

// func Check(rules []model.Rule, newRule model.Rule) int {
// }
