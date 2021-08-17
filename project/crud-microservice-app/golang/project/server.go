package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

var mongoCrudURL string
var textAPIURL string

//Handles all requests coming in
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	fmt.Printf("Handling requests on 4000...\n")

	logWriter("We are now handling requests")
	//Serve our User Crud API
	myRouter.HandleFunc("/addUser", addUser).Methods("POST")            //Add a User
	myRouter.HandleFunc("/deleteUser", deleteUser).Methods("POST")      //Delete a User
	myRouter.HandleFunc("/updateUser", updateUser).Methods("POST")      //Update a User
	myRouter.HandleFunc("/getUser", getUser).Methods("POST")            //Get User
	myRouter.HandleFunc("/giveAllEmails", giveAllEmails).Methods("GET") //Get All user Emails
	//Serve our LearnROrg Crud API
	myRouter.HandleFunc("/addLearnOrg", addLearnOrg).Methods("POST")               //Add a LearnROrg
	myRouter.HandleFunc("/deleteLearnOrg", deleteLearnOrg).Methods("POST")         //Delete a LearnROrg
	myRouter.HandleFunc("/updateLearnOrg", updateLearnOrg).Methods("POST")         //Update a LearnROrg
	myRouter.HandleFunc("/getLearnOrg", getLearnOrg).Methods("POST")               //Get LearnROrg
	myRouter.HandleFunc("/getLearnOrgAdminOf", getLearnOrgAdminOf).Methods("POST") //Get LearnROrg this user is admin of
	//Serve our LearnR Crud API
	myRouter.HandleFunc("/addLearnR", addLearnR).Methods("POST")                 //Add a LearnR
	myRouter.HandleFunc("/deleteLearnR", deleteLearnR).Methods("POST")           //Delete a LearnR
	myRouter.HandleFunc("/updateLearnR", updateLearnR).Methods("POST")           //Update a LearnR
	myRouter.HandleFunc("/getLearnR", getLearnR).Methods("POST")                 //Get LearnR
	myRouter.HandleFunc("/specialLearnRGive", specialLearnRGive).Methods("POST") //Gets an array of special learnrs
	myRouter.HandleFunc("/getLearnRArray", getLearnRArray).Methods("POST")       //Gets an array of learnrs
	//Serve our LearnRInfo Crud API
	myRouter.HandleFunc("/addLearnrInfo", addLearnrInfo).Methods("POST")       //Add a LearnRInfo
	myRouter.HandleFunc("/deleteLearnrInfo", deleteLearnrInfo).Methods("POST") //Delete a LearnRInfo
	myRouter.HandleFunc("/updateLearnrInfo", updateLearnrInfo).Methods("POST") //Update a LearnRInfo
	myRouter.HandleFunc("/getLearnrInfo", getLearnrInfo).Methods("POST")       //Get LearnRInfo
	//Serve our LearnRSessions CRUD API
	myRouter.HandleFunc("/addLearnRSession", addLearnRSession).Methods("POST")       //Add a LearnRSession
	myRouter.HandleFunc("/deleteLearnRSession", deleteLearnRSession).Methods("POST") //Delete a LearnRSession
	myRouter.HandleFunc("/updateLearnRSession", updateLearnRSession).Methods("POST") //Update a LearnRSession
	myRouter.HandleFunc("/getLearnRSession", getLearnRSession).Methods("POST")       //Get LearnRSession
	//Serve our LearnRInforms CRUD API
	myRouter.HandleFunc("/addLearnRInforms", addLearnRInforms).Methods("POST")       //Add a LearnRInfo
	myRouter.HandleFunc("/deleteLearnRInforms", deleteLearnRInforms).Methods("POST") //Delete a LearnRInfo
	myRouter.HandleFunc("/updateLearnRInforms", updateLearnRInforms).Methods("POST") //Update a LearnRInfo
	myRouter.HandleFunc("/getLearnRInforms", getLearnRInforms).Methods("POST")       //Get LearnRInfo
	//Serve our validation APIs
	myRouter.HandleFunc("/giveAllUsernames", giveAllUsernames).Methods("GET") //Get all our Usernames
	myRouter.HandleFunc("/giveAllLearnROrg", giveAllLearnROrg).Methods("GET") //Get all our LearnROrg Names
	myRouter.HandleFunc("/giveAllLearnr", giveAllLearnr).Methods("GET")
	myRouter.HandleFunc("/randomIDCreationAPI", randomIDCreationAPI).Methods("GET") //Get a random ID
	myRouter.HandleFunc("/userLogin", userLogin).Methods("POST")                    //Check if User can login
	myRouter.HandleFunc("/addEmailVerif", addEmailVerif).Methods("POST")            //Adds email verif
	myRouter.HandleFunc("/getEmailVerif", getEmailVerif).Methods("POST")            //Get email verif
	myRouter.HandleFunc("/deleteEmailVerify", deleteEmailVerify).Methods("POST")    //Delete email verif
	myRouter.HandleFunc("/updateEmailVerify", updateEmailVerify).Methods("POST")    //Update Email verif
	//Serve our test ping
	myRouter.HandleFunc("/testPingPost", testPingPost).Methods("POST") //Get a random ID
	myRouter.HandleFunc("/testPingGet", testPingGet).Methods("GET")    //Get a random ID
	//Serve response for services checking if we're up
	myRouter.HandleFunc("/available", available).Methods("GET") //See if this service is available
	//Serve our static files
	log.Fatal(http.ListenAndServe(":4000", myRouter))
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

func testPingPost(w http.ResponseWriter, req *http.Request) {
	//Declare data to return
	type ReturnMessage2 struct {
		TheErr          []string        `json:"TheErr"`
		ResultMsg       []string        `json:"ResultMsg"`
		SuccOrFail      int             `json:"SuccOrFail"`
		ReturnedUserMap map[string]bool `json:"ReturnedUserMap"`
	}
	theReturnMessage := ReturnMessage2{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success
	theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, "You've got a successful response")

	fmt.Printf("DEBUG: Successful ping to testPingPost\n")

	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	type LoginData struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}

	//Marshal the user data into our type
	var dataForLogin LoginData
	json.Unmarshal(bs, &dataForLogin)

	/* Test get User */
	/* User collection */
	userCollection := mongoClient.Database("learnR").Collection("users") //Here's our collection
	var testAUser User
	theErr := userCollection.FindOne(theContext, bson.M{"userid": 228778447811}).Decode(&testAUser)
	if theErr != nil {
		if strings.Contains(theErr.Error(), "no documents in result") {
			fmt.Printf("DEBUG: We didn't find this User in the search...\n")
		} else {
			theErr := "There is another error getting random ID: " + err.Error()
			fmt.Println(theErr)
			logWriter(theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			log.Fatal(theErr)
		}
	}

	/* Return the marshaled response */
	//Send the response back
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in testPingPost: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

func testPingGet(w http.ResponseWriter, r *http.Request) {
	type ReturnMessage struct {
		TheErr          []string        `json:"TheErr"`
		ResultMsg       []string        `json:"ResultMsg"`
		SuccOrFail      int             `json:"SuccOrFail"`
		ReturnedUserMap map[string]bool `json:"ReturnedUserMap"`
	}
	theReturnMessage := ReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success
	theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, "You've got a successful response")

	fmt.Printf("DEBUG: Successful ping to testPingGET\n")

	/* Return the marshaled response */
	//Send the response back
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in testPingGet: " + err.Error()
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
		ResultMsg:  []string{"Good return from available for this  CRUD Microservice"},
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
