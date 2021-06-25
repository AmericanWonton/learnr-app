package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const TESTPINGURL string = "http://13.59.100.23/testLocalPing"

//Handles all requests coming in
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	defer close(learnSessChannel)       //Close the channel when needed
	defer close(learnSessResultChannel) //Close this channel when needed

	go learnRSession(learnSessChannel, learnSessResultChannel)

	fmt.Printf("DEBUG: Running on port 3000...\n")
	//Web request/Text Request handling
	myRouter.HandleFunc("/initialLearnRStart", initialLearnRStart).Methods("POST") //Handle incoming learnr initiations
	//Test Ping to our Server
	myRouter.HandleFunc("/testLocalPing", testLocalPing).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", myRouter))
	//Log results of our channel jobs being completed
	for a := 0; a <= len(learnSessResultChannel); a++ {
		aUserSess := <-learnSessResultChannel
		aMessage := "We are done with this learnRSess: " + strconv.Itoa(aUserSess.TheSession.ID) + " for this LearnR: \n" +
			aUserSess.TheSession.LearnRName + "\n"
		logWriter(aMessage)
		fmt.Println(aMessage)
	}
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
