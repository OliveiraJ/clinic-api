package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/OliveiraJ/clinic-api/model"
)

const PATH string = "./database/rules.json"

// Reads the rules.json file and returns all the data from it to a map of type [string]model.Rule, the exists function is called so it can be verified
// if the rules.json file exists, if it doesnt, the function will create a new one in the following path ./database/rules.json
func ReadJson(pathFileJson string) map[string]model.Rule {
	if !exists(pathFileJson) {
		jsonFile, err := os.Create(pathFileJson)
		if err != nil {
			log.Fatal(err)
		}
		defer jsonFile.Close()
	}

	jsonFile, err := os.Open(pathFileJson)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValueJSON, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	Rules := make(map[string]model.Rule)
	json.Unmarshal(byteValueJSON, &Rules)

	return Rules
}

// Writes the data from a map of type [string]model.Rule to a json file.
func WriteJson(rules map[string]model.Rule, pathFileJson string) {
	file, err := json.MarshalIndent(rules, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(pathFileJson, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Verifys if a file/folders exists.
func exists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Makes the received rule a dayly rule, propagating it to all the days within the received interval, returns a map of type [string]model.Rule
func Dayly(rules map[string]model.Rule, rule model.Rule) (map[string]model.Rule, bool) {

	var check bool
	limit := time.Time(rule.Limit)
	start := time.Time(rule.Day)
	days := limit.Sub(start).Hours() / (24)

	day := time.Time(rule.Day)

	for i := 0; i < int(days)+1; i++ {
		log.Println(i, time.Time(rule.Day))
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
			log.Println(day)
			rule.Day = model.CustomDay(day)

			rules[time.Time(rule.Day).Format(model.DAY)] = rule
		}
		day = day.AddDate(0, 0, 1)
		log.Println("Imprimindo o day: ", day)
	}

	return rules, false
}

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
