package utils

import (
	"time"

	"github.com/OliveiraJ/clinic-api/model"
)

// Makes the received rule a dayly rule, propagating it to all the days within the received interval, returns a map of type [string]model.Rule
func Dayly(rules map[string]model.Rule, rule model.Rule) (map[string]model.Rule, bool) {

	var check bool
	limit := time.Time(rule.Limit)
	start := time.Time(rule.Day)
	days := limit.Sub(start).Hours() / (24)

	day := time.Time(rule.Day)

	for i := 0; i < int(days)+1; i++ {
		if foundRule, found := rules[day.Format(model.DAY)]; found {
			check = CheckInvalidSchedule(foundRule, rule)
			if check {
				return rules, check
			}
			check = CheckOverlapingIntervals(foundRule, rule)
			if check {
				return rules, check
			}
			foundRule.Intervals = append(foundRule.Intervals, rule.Intervals...)
			rules[day.Format(model.DAY)] = foundRule
		} else {
			rule.Day = model.CustomDay(day)

			rules[time.Time(rule.Day).Format(model.DAY)] = rule
		}
		day = day.AddDate(0, 0, 1)
	}

	return rules, false
}
