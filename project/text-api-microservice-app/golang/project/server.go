package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var mongoCrudURL string
var textAPIURL string

var TESTPINGURL string = mongoCrudURL + "/testLocalPing"

//Handles all requests coming in
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	fmt.Printf("DEBUG: Running on port 3000...\n")
	//Web request/Text Request handling
	myRouter.HandleFunc("/initialLearnRStart", initialLearnRStart).Methods("POST") //Handle incoming learnr initiations
	myRouter.HandleFunc("/textWebhook", textWebhook).Methods("POST")               //Handle incoming webhook texts from Users
	//Test Ping to our Server
	myRouter.HandleFunc("/testLocalPing", testLocalPing).Methods("POST")
	myRouter.HandleFunc("/httpTakerFunc", httpTakerFunc).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", myRouter))
}

func testLocalPing(w http.ResponseWriter, r *http.Request) {
	//Declare Ajax return statements to be sent back
	type SuccessMSG struct {
		Message    string `json:"Message"`
		SuccessNum int    `json:"SuccessNum"`
	}
	theSuccMessage := SuccessMSG{
		Message:    "LearnR Organization created successfully",
		SuccessNum: 0,
	}

	//Declare struct we are expecting
	type SendJSON struct {
		TestNum int `json:"TestNum"`
	}
	//Get the byte slice from the request
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	//Marshal it into our type
	var ourJSON SendJSON
	json.Unmarshal(bs, &ourJSON)

	/* Send the response back to Ajax */
	theJSONMessage, err := json.Marshal(theSuccMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in createUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//Loads in the initial text API and MongoCrud URLS
func loadInMicroServiceURL() {
	//Check to see if ENV Creds are available first
	_, ok := os.LookupEnv("CRUD_URL")
	if !ok {
		message := "This ENV Variable is not present: " + "CRUD_URL"
		panic(message)
	}
	_, ok2 := os.LookupEnv("TEXT_API")
	if !ok2 {
		message := "This ENV Variable is not present: " + "TEXT_API"
		panic(message)
	}

	mongoCrudURL = os.Getenv("CRUD_URL")
	textAPIURL = os.Getenv("TEXT_API")

	fmt.Printf("DEBUG: Here is mongo: %v\n and here is text: %v\n", mongoCrudURL, textAPIURL)
}
