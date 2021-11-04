package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/OliveiraJ/clinic-api/model"
)

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
