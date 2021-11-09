package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitializeRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/clinicapi/v1/rule/{start}/{end}", AvailableDays).Methods("GET")
	router.HandleFunc("/clinicapi/v1/rule", CreateRule).Methods("POST")
	router.HandleFunc("/clinicapi/v1/rule/dayly", CreateDaylyRule).Methods("POST")
	router.HandleFunc("/clinicapi/v1/rule/weekly", CreateWeeklyRule).Methods("POST")
	router.HandleFunc("/clinicapi/v1/rule/{key}", DeleteRule).Methods("DELETE")
	router.HandleFunc("/clinicapi/v1/rule/interval/{key}", DeleteInterval).Methods("DELETE")
	router.HandleFunc("/clinicapi/v1/rule/{key}", UpdateRule).Methods("PATCH")
	router.HandleFunc("/clinicapi/v1/rules", GetRules).Methods("GET")
	router.HandleFunc("/clinicapi/v1/rule/{key}", GetRule).Methods("GET")

	log.Println("Listening on port :9090")
	log.Fatal(http.ListenAndServe(":9090", router))
}
