package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/OliveiraJ/clinic-api/model"
)

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
