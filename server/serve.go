package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/OliveiraJ/clinic-api/model"
	"github.com/OliveiraJ/clinic-api/utils"
	"github.com/gorilla/mux"
)

// Creates a rule and add new valid schedules to those that already exists
func CreateRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var exist bool
	var rule model.Rule
	json.NewDecoder(r.Body).Decode(&rule)

	rules := utils.ReadJson(utils.PATH)
	key := time.Time(rule.Day).Format(model.DAY)

	if v, found := rules[key]; found {

		for _, p := range v.Intervals {
			for _, pr := range rule.Intervals {
				if p == pr {
					exist = true
				}
			}
		}
		if exist {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Error: Rule already exists"))

		} else if utils.CheckInvalidInterval(rule) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: Invalid intervals format, no possible to schedule"))

		} else if utils.CheckInvalidDate(rule) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: Limite date and start date do not match"))

		} else if utils.CheckOverlappingIntervals(v, rule) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: Overlaping schedules"))
		} else {
			v.Intervals = append(v.Intervals, rule.Intervals...)
			rules[key] = v
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(rule)
		}

	} else {
		rules[key] = rule
		log.Println("New rule created")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(rule)
	}

	utils.WriteJson(rules, utils.PATH)

}

// Creates a dayly rule on the set of days between day and limit dates.
func CreateDaylyRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var invalidSchedule bool
	var rule model.Rule

	json.NewDecoder(r.Body).Decode(&rule)
	rules := utils.ReadJson(utils.PATH)

	rules, invalidSchedule = utils.Dayly(rules, rule)
	if invalidSchedule {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Error: Schedule Conflict"))
	} else if utils.CheckInvalidInterval(rule) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Error: Interval Conflict"))
	} else if utils.CheckInvalidDate(rule) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Error: Invalid Date Interval (Equal values to day and limit or 'limit' preceds 'day')"))
	} else {
		w.WriteHeader(http.StatusCreated)
		utils.WriteJson(rules, utils.PATH)
		json.NewEncoder(w).Encode(rule)
		log.Println("New dayly rule created successfully!")
	}
}

// Creates a Weekly Rule based on the values of 'limit' and 'day', starting from the date represented by the 'day' variable
// a new rule is created in each 7 days gap until it reach the closest date to 'limit'. As an example:
// If day: 01-11-2021 and limit: 01-12-2021, then rules would be created on 01-11-2021, 08-11-2021, 15-11-2021, 22-11-2021 and
// 29-11-2021, ending on the closest day to the limit date.
func CreateWeeklyRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var invalidSchedule bool
	var rule model.Rule

	json.NewDecoder(r.Body).Decode(&rule)
	rules := utils.ReadJson(utils.PATH)

	rules, invalidSchedule = utils.Weekly(rules, rule)

	if invalidSchedule {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Error: Schedule Conflict"))
	} else if utils.CheckInvalidInterval(rule) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Error: Interval Conflict"))
	} else if utils.CheckInvalidDate(rule) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Error: Invalid Date Interval (Equal values to day and limit or 'limit' preceds 'day')"))
	} else {
		w.WriteHeader(http.StatusCreated)
		utils.WriteJson(rules, utils.PATH)
		json.NewEncoder(w).Encode(rule)
		log.Println("New weekly rule created successfully!")
	}

}

// Return every rule present on the json file (database)
func GetRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rules := utils.ReadJson(utils.PATH)
	var extRules []model.ExtRule

	for _, v := range rules {
		extRules = append(extRules, model.ExtRule(v))

	}

	json.NewEncoder(w).Encode(extRules)
	log.Println("Returning all rules!")
}

// Returns a specific rule
func GetRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	rules := utils.ReadJson(utils.PATH)

	if _, found := rules[params["key"]]; found {
		rule := rules[params["key"]]

		//Conversion to type ExtRule (external rule), so only the needed information is showed to the client
		extRule := model.ExtRule(rule)
		json.NewEncoder(w).Encode(extRule)
		log.Println("Returning rule " + params["key"])
	} else {
		log.Println("Couldn't find " + params["key"] + " rule!")
		w.WriteHeader(http.StatusNotFound)
	}
}

// Deletes a specific rule.
func DeleteRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	rules := utils.ReadJson(utils.PATH)

	if _, found := rules[params["key"]]; found {
		delete(rules, params["key"])
		log.Println("Rule " + params["key"] + " deleted successefully!")
	} else {
		log.Println("Couldn't delete " + params["key"] + " rule!")
		w.WriteHeader(http.StatusNotFound)
	}

	utils.WriteJson(rules, utils.PATH)
	json.NewEncoder(w).Encode(rules)
}

// Deletes a interval or a schedule from an specific rule.
func DeleteInterval(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var period model.Interval
	rules := utils.ReadJson(utils.PATH)

	json.NewDecoder(r.Body).Decode(&period)

	if v, found := rules[params["key"]]; found {
		for i, p := range v.Intervals {
			if p == period {
				v.Intervals = append(v.Intervals[:i], v.Intervals[i+1:]...)
			}
		}

		rules[params["key"]] = v

		log.Println("Schedule " + period.Start.String() + " to " + period.End.String() + " deleted successefully!")
	} else {
		log.Println("Couldn't delete " + period.Start.String() + " to " + period.End.String() + " schedule")
		w.WriteHeader(http.StatusNotFound)
	}

	utils.WriteJson(rules, utils.PATH)
	json.NewEncoder(w).Encode(rules)
}

// Returns a set of rules that fit a interval of days
func AvailableDays(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//queryParams := r.URL.Query()
	rules := utils.ReadJson(utils.PATH)
	var extRules []model.ExtRule

	start := params["start"]
	end := params["end"]

	if startRule, foundStart := rules[start]; foundStart {
		if endRule, foundEnd := rules[end]; foundEnd {
			day := time.Time(startRule.Day)
			for (day.Equal(time.Time(startRule.Day)) || day.After(time.Time(startRule.Day))) && (day.Equal(time.Time(endRule.Day)) || day.Before(time.Time(endRule.Day))) {
				if time.Time(rules[day.Format(model.DAY)].Day).Format(model.DAY) != "01-01-0001" {

					extRules = append(extRules, model.ExtRule(rules[day.Format(model.DAY)]))

				}
				day = day.AddDate(0, 0, 1)
			}

			json.NewEncoder(w).Encode(extRules)
			log.Println("Returning rules!")

		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

// Updates an specific schedule
func UpdateRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	rules := utils.ReadJson(utils.PATH)
	var rule model.Rule

	json.NewDecoder(r.Body).Decode(&rule)

	if x, found := rules[params["key"]]; found {
		if utils.CheckInvalidInterval(rule) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Schedule already exists"))
		} else if utils.CheckOverlappingIntervals(x, rule) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Overlaping schedule, try a diferent one"))
		} else {
			x.Intervals = append(x.Intervals, rule.Intervals...)
			rules[params["key"]] = x

			utils.WriteJson(rules, utils.PATH)
			json.NewEncoder(w).Encode(rule)
		}
	} else {
		log.Println("Rule not found!")
		w.WriteHeader(http.StatusNotFound)
	}

}
