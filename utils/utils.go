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
func Dayly(rules map[string]model.Rule, rule model.Rule) map[string]model.Rule {
	for rule.Day != rule.Limit {
		day := time.Time(rule.Day).AddDate(0, 0, 1)
		rule.Day = model.CustomDay(day)
		key := time.Time(rule.Day).Format(model.DAY)
		rules[key] = rule
	}

	return rules
}

// Weekly transforms the received rule in a weekly rule, propagating it to every week within the received interval, returns a map of type [string]model.Rule
func Weekly(rules map[string]model.Rule, rule model.Rule) map[string]model.Rule {
	limit := time.Time(rule.Limit)
	start := time.Time(rule.Day)
	days := limit.Sub(start).Hours() / (24 * 7)

	if v, found := rules[start.Format(model.DAY)]; found {

		rule = CheckSchedule(v, rule)

		if len(rule.Intervals) != 0 {

			v.Intervals = append(v.Intervals, rule.Intervals...)
			rules[start.Format(model.DAY)] = v
		}

	} else {

		rule.Day = model.CustomDay(start)
		log.Println(time.Time(rule.Day))
		rules[start.Format(model.DAY)] = rule
		log.Println("novo criado")

	}

	for i := 0; i < int(days); i++ {

		day := time.Time(rule.Day).AddDate(0, 0, 7)

		if v, found := rules[day.Format(model.DAY)]; found {

			rule = CheckSchedule(v, rule)
			v.Intervals = append(v.Intervals, rule.Intervals...)
			rules[time.Time(rule.Day).Format(model.DAY)] = v

		} else {

			log.Println(day)
			rule.Day = model.CustomDay(day)

			rules[time.Time(rule.Day).Format(model.DAY)] = rule

		}
	}

	return rules
}

// Checks if there are any existing rules, preventing duplicate rules from being added as those are removed from the slice rule.Intervals.
func CheckSchedule(x model.Rule, rule model.Rule) model.Rule {

	for _, interval := range x.Intervals {
		for i, ruleInterval := range rule.Intervals {
			if interval == ruleInterval {
				log.Println("horário já existe")

				rule.Intervals = append(rule.Intervals[:i], rule.Intervals[i+1:]...)
				log.Println(rule.Intervals)
			}
		}
	}

	return rule
}
