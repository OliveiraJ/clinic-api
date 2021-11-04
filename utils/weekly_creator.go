package utils

import (
	"log"
	"time"

	"github.com/OliveiraJ/clinic-api/model"
)

// Weekly transforms the received rule in a weekly rule, propagating it to every week within the received interval, returns a map of type [string]model.Rule
func Weekly(rules map[string]model.Rule, rule model.Rule) (map[string]model.Rule, bool) {
	var exist bool
	limit := time.Time(rule.Limit)
	start := time.Time(rule.Day)
	days := limit.Sub(start).Hours() / (24 * 7)

	day := time.Time(rule.Day)

	for i := 0; i < int(days)+1; i++ {
		if v, found := rules[day.Format(model.DAY)]; found {
			exist = CheckInvalidSchedule(v, rule)
			if exist {
				return rules, exist
			}
			v.Intervals = append(v.Intervals, rule.Intervals...)
			rules[day.Format(model.DAY)] = v
		} else {
			log.Println(day)
			rule.Day = model.CustomDay(day)

			rules[day.Format(model.DAY)] = rule
		}
		day = day.AddDate(0, 0, 7)
	}

	return rules, false
}
