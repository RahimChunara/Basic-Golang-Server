package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", returnAllRecords)
	myRouter.HandleFunc("/record", createNewRecord).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", myRouter))
}

type Record struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var Records []Record

func returnAllRecords(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllRecords")
	json.NewEncoder(w).Encode(Records)
}

func createNewRecord(w http.ResponseWriter, r *http.Request) {
	var newRecord Record

	err := json.NewDecoder(r.Body).Decode(&newRecord)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, s := range Records {
		if s.Id == newRecord.Id {
			Records[i] = newRecord
			break
		} else {
			Records = append(Records, newRecord)
		}
	}
	if len(Records) == 0 {
		Records = append(Records, newRecord)
	}
}

func main() {
	Records = []Record{}
	handleRequests()
}
