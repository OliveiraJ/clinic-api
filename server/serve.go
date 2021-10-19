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

func InitializeRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/rule", CreateRule).Methods("POST")
	router.HandleFunc("/rule/dayly", CreateDaylyRule).Methods("POST")
	router.HandleFunc("/rule/weekly", CreateWeeklyRule).Methods("POST")
	router.HandleFunc("/rule/{key}", DeleteRule).Methods("DELETE")
	router.HandleFunc("/rule/{key}", UpdateRule).Methods("PUT")
	router.HandleFunc("/rule/{key}", GetRule).Methods("GET")
	router.HandleFunc("/rules", GetRules).Methods("GET")
	router.HandleFunc("/availabledays", AvailableDays).Methods("GET")

	log.Println("Listening on port :9090")
	log.Fatal(http.ListenAndServe(":9090", router))
}

func CreateRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rule model.Rule
	json.NewDecoder(r.Body).Decode(&rule)

	rules := utils.ReadJson(utils.PATH)
	key := time.Time(rule.Day).Format(model.DAY)

	if _, found := rules[key]; found {
		w.WriteHeader(405)
		w.Write([]byte("405 - Rule already exists"))
	} else {
		rules[key] = rule
		log.Println("New rule created with success!")
		json.NewEncoder(w).Encode(rule)
	}

	utils.WriteJson(rules, utils.PATH)

}

func CreateDaylyRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rule model.Rule

	json.NewDecoder(r.Body).Decode(&rule)
	rules := utils.ReadJson(utils.PATH)

	rules = utils.Dayly(rules, rule)

	utils.WriteJson(rules, utils.PATH)
	json.NewEncoder(w).Encode(rule)
	log.Println("New dayly rule created successfully!")

}

func CreateWeeklyRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rule model.Rule

	json.NewDecoder(r.Body).Decode(&rule)
	rules := utils.ReadJson(utils.PATH)

	rules = utils.Weekly(rules, rule)

	utils.WriteJson(rules, utils.PATH)
	json.NewEncoder(w).Encode(rule)
	log.Println("New weekly rule created successfully!")

}

// Return everyu rule present on the json file
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
		w.WriteHeader(404)
	}
}

func DeleteRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	rules := utils.ReadJson(utils.PATH)

	if _, found := rules[params["key"]]; found {
		delete(rules, params["key"])
		log.Println("Rule " + params["key"] + " deleted successefully!")
	} else {
		log.Println("Couldn't delete " + params["key"] + " rule!")
		w.WriteHeader(404)
	}

	utils.WriteJson(rules, utils.PATH)
	json.NewEncoder(w).Encode(rules)
}

func DeleteInterval(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	rules := utils.ReadJson(utils.PATH)

	if _, found := rules[params["key"]]; found {
		delete(rules, params["key"])
		log.Println("Rule " + params["key"] + " deleted successefully!")
	} else {
		log.Println("Couldn't delete " + params["key"] + " rule!")
		w.WriteHeader(404)
	}

	utils.WriteJson(rules, utils.PATH)
	json.NewEncoder(w).Encode(rules)
}

func AvailableDays(w http.ResponseWriter, r *http.Request) {

}

// Pending system to avoid repeating schedules
func UpdateRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	rules := utils.ReadJson(utils.PATH)
	var rule model.Rule

	json.NewDecoder(r.Body).Decode(&rule)

	if x, found := rules[params["key"]]; found {

		for _, v := range rules {
			for _, interval := range v.Intervals {
				for i, ruleInterval := range rule.Intervals {
					if interval == ruleInterval {
						log.Println("horário já existe")

						x.Intervals = append(x.Intervals[:i], rule.Intervals[i+1:]...)
					}
				}
			}
		}
		x.Intervals = append(x.Intervals, rule.Intervals...)
		rules[params["key"]] = x
	} else {
		log.Println("Rule not found!")
		w.WriteHeader(404)
	}
	utils.WriteJson(rules, utils.PATH)
	json.NewEncoder(w).Encode(rule)

}
