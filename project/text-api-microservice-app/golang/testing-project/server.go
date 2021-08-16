package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const TESTPINGURL string = "http://13.59.100.23/testLocalPing"

//Handles all requests coming in
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	fmt.Printf("DEBUG: Running on port 3000...\n")
	//Web request/Text Request handling
	myRouter.HandleFunc("/initialLearnRStart", initialLearnRStart).Methods("POST")         //Handle incoming learnr initiations
	myRouter.HandleFunc("/initialBulkLearnRStart", initialBulkLearnRStart).Methods("POST") //Handle incoming Bulk learnr initiations
	myRouter.HandleFunc("/textWebhook", textWebhook).Methods("POST")                       //Handle incoming webhook texts from Users
	//Test Ping to our Server
	myRouter.HandleFunc("/testLocalPing", testLocalPing).Methods("POST")
	//Serve response for services checking if we're up
	myRouter.HandleFunc("/available", available).Methods("GET") //See if this service is available
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

//A service that returns if this Microservice is up and running
func available(w http.ResponseWriter, r *http.Request) {
	//Declare data to return
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	theReturnMessage := ReturnMessage{
		TheErr:     []string{""},
		ResultMsg:  []string{"Good return from available for this text Microservice"},
		SuccOrFail: 0,
	}

	//Format the JSON map for returning our results
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in available: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}
