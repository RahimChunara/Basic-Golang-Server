package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"text/template"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

// Redis Client
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", returnAllRecords)
	myRouter.HandleFunc("/record", createNewRecord).Methods("POST")
	myRouter.HandleFunc("/v1/{id}", deleteRecord).Methods("DELETE")
	myRouter.Handle("/w1", websocket.Handler(Server))
	myRouter.HandleFunc("/x1", serveStatic)
	myRouter.HandleFunc("/v1", getSpecificRecord)
	myRouter.HandleFunc("/v1/random", randomImage)

	log.Fatal(http.ListenAndServe(":3001", myRouter))

}

var ctx = context.Background()

// Stores All Records
type Record struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Age         int         `json:"age"`
	Description interface{} `json:"description"`
}

var Records []Record

// Render HTML Doc
func serveStatic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: html")
	fp := path.Join("templates", "client.html")
	t, err := template.ParseFiles(fp)
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, nil)
}

// Get Specific Record, if none return all
func getSpecificRecord(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("ojmefrooegrwfnjm")
	fmt.Println("query params:", r.URL.Query())
	id := r.URL.Query().Get("id")

	val, err := rdb.HGet(ctx, "AllRecords", "id:"+id).Result()
	if err != nil {
		fmt.Println("Throw Error")
		w.WriteHeader(404)
	} else {
		json.NewEncoder(w).Encode(val)
	}

}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("Endpoint Hit: deleteRecord")
	val := rdb.HDel(ctx, "AllRecords", "id:"+id).Val()
	if val == 0 {
		w.WriteHeader(404)
	}
}

func returnAllRecords(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllRecords")
	val := rdb.HGetAll(ctx, "AllRecords").Val()
	json.NewEncoder(w).Encode(val)
	fmt.Print(val)
}

// Create New Record, also updates record if same if
func createNewRecord(w http.ResponseWriter, r *http.Request) {
	var newRecord Record
	const Record string = "Record:"

	fmt.Println("Endpoint Hit: Create new record")
	err := json.NewDecoder(r.Body).Decode(&newRecord)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Print(newRecord)

	b, err := json.Marshal(newRecord)
	if err != nil {
		return
	}

	er := rdb.HSet(ctx, "AllRecords", "id:"+newRecord.Id, b).Err()
	if er != nil {
		panic(err)
	}

}

// Echo Server
func Server(ws *websocket.Conn) {
	fmt.Println("Endpoint Hit: Echo Server")
	var err error

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

// Unmarsha; Json to this struct
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

var client = http.Client{}

// Fetches Random image from the api and displays it on client
func randomImage(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://dog.ceo/api/breeds/image/random")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var responseObject Response
	json.Unmarshal(body, &responseObject)

	reqImg, err := client.Get(responseObject.Message)
	if err != nil {
		log.Fatalln(err)
	}

	defer reqImg.Body.Close()

	if _, err = io.Copy(w, reqImg.Body); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	Records = []Record{}
	handleRequests()
}
