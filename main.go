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
	myRouter.HandleFunc("/v1/{id}", deleteRecord).Methods("DELETE")
	myRouter.HandleFunc("/v1/{id}", getSpecificRecord)

	log.Fatal(http.ListenAndServe(":3000", myRouter))
}

type Record struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Age         int         `json:"age"`
	Description interface{} `json:"description"`
}

var Records []Record

func getSpecificRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, record := range Records {
		if record.Id == key {
			json.NewEncoder(w).Encode(record)
			fmt.Println("HTTP Response Status:", 200, http.StatusText(200))
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("Endpoint Hit: deleteRecord")
	for index, record := range Records {
		if record.Id == id {
			Records = append(Records[:index], Records[index+1:]...)
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}

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
