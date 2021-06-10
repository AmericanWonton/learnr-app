package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

/* Both are used for usernames below */
var allUsernames []string
var usernameMap map[string]bool

//Handles the Index requests; Ask User if they're legal here
func index(w http.ResponseWriter, r *http.Request) {
	/* REdirect, Index not needed */
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

//Handles login/ page
func login(w http.ResponseWriter, r *http.Request) {
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "login.gohtml", nil)
	HandleError(w, err1)
}

//Handles the signup page
func signup(w http.ResponseWriter, r *http.Request) {
	usernameMap = loadUsernames() //Load all usernames
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "signup.gohtml", nil)
	HandleError(w, err1)
}

// Handle Errors passing templates
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		fmt.Printf("We had an error loading this template: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

//Calls 'giveAllUsernames' to run a mongo query to get all Usernames, then puts it in a map to return
func loadUsernames() map[string]bool {
	mapOusernameToReturn := make(map[string]bool) //Username to load our values into
	//Call our crudOperations Microservice in order to get our Usernames
	//Create a context for timing out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", GETALLUSERNAMESURL, nil)
	if err != nil {
		theErr := "There was an error getting Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "There was an error reaching out to loadUsername API: " + strconv.Itoa(resp.StatusCode)
		fmt.Println(theErr)
		logWriter(theErr)
	} else if err != nil {
		theErr := "Error from response to loadUsernames: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	//Marshal the response into a type we can read
	type ReturnMessage struct {
		TheErr          []string        `json:"TheErr"`
		ResultMsg       []string        `json:"ResultMsg"`
		SuccOrFail      int             `json:"SuccOrFail"`
		ReturnedUserMap map[string]bool `json:"ReturnedUserMap"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)

	//Assign our map variable to the map varialbe and see if it's okay
	if returnedMessage.SuccOrFail != 0 {
		errString := ""
		for l := 0; l < len(returnedMessage.TheErr); l++ {
			errString = errString + returnedMessage.TheErr[l]
		}
		logWriter(errString)
		fmt.Println(errString)
	} else {
		mapOusernameToReturn = returnedMessage.ReturnedUserMap
	}

	return mapOusernameToReturn
}
