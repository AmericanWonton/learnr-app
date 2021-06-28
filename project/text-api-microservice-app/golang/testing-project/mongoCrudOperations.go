package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Get User
//User
const READUSERURL string = "http://localhost:4000/getUser"

//LearnRGet
const READLEARNRURL string = "http://localhost:4000/getLearnR"

//LearnRInfoGet
const READLEARNRINFOURL string = "http://localhost:4000/getLearnrInfo"
const UPDATELEARNRINFOURL string = "http://localhost:4000/updateLearnrInfo"

//ID Random
const GETRANDOMID string = "http://localhost:4000/randomIDCreationAPI"

//LearnRSession
const ADDLEARNRSESSIONSURL string = "http://localhost:4000/addLearnRSession"
const READLEARNRSESSIONSURL string = "http://localhost:4000/getLearnRSession"
const UPDATELEARNRSESSIONSURL string = "http://localhost:4000/updateLearnRSession"
const DELETELEARNRSESSIONSURL string = "http://localhost:4000/deleteLearnRSession"

//Mongo DB Declarations
var mongoClient *mongo.Client

var theContext context.Context
var mongoURI string //Connection string loaded

/* Variable definitions for User/Learnr */
type User struct {
	UserName    string   `json:"UserName"`
	Password    string   `json:"Password"`
	Firstname   string   `json:"Firstname"`
	Lastname    string   `json:"Lastname"`
	PhoneNums   []string `json:"PhoneNums"`
	UserID      int      `json:"UserID"`
	Email       []string `json:"Email"`
	Whoare      string   `json:"Whoare"`
	AdminOrgs   []int    `json:"AdminOrgs"`
	OrgMember   []int    `json:"OrgMember"`
	Banned      bool     `json:"Banned"`
	DateCreated string   `json:"DateCreated"`
	DateUpdated string   `json:"DateUpdated"`
}

//LearnR Org
type LearnrOrg struct {
	OrgID       int      `json:"OrgID"` //Unique ID of this organization
	Name        string   `json:"Name"`  //Name of this organization
	OrgGoals    []string //A list of goals for this organization
	UserList    []int    //All the Users in this organization
	AdminList   []int    //A list of all the Admins in this organization,(UserIDs)
	LearnrList  []int    //A list of all learnr ints in this organization
	DateCreated string   `json:"DateCreated"`
	DateUpdated string   `json:"DateUpdated"`
}

//LearnR
type Learnr struct {
	ID            int             `json:"ID"`            //ID of this LearnR
	InfoID        int             `json:"InfoID"`        //Links to the LearnRInfo object which holds data
	OrgID         int             `json:"OrgID"`         //Which organization does this belong to
	Name          string          `json:"Name"`          //Name of this LearnR
	Tags          []string        `json:"Tags"`          //Tags that describe this LearnR
	Description   []string        `json:"Description"`   //Description of this LearnR
	PhoneNums     []string        `json:"PhoneNums"`     //Phone Nums attatched to this LearnR
	LearnRInforms []LearnRInforms `json:"LearnRInforms"` //What we'll text to our Users
	Active        bool            `json:"Active"`        //Whether this LearnR is still active
	DateCreated   string          `json:"DateCreated"`
	DateUpdated   string          `json:"DateUpdated"`
}

//LearnRInfo
type LearnrInfo struct {
	ID               int             `json:"ID"`               //ID of this LearnR Info
	LearnRID         int             `json:"LearnRID"`         //The LearnR ID related to this info
	AllSessions      []LearnRSession `json:"AllSessions"`      //An array of all the sessions
	FinishedSessions []LearnRSession `json:"FinishedSessions"` //An array of complete sessions only
	DateCreated      string          `json:"DateCreated"`
	DateUpdated      string          `json:"DateUpdated"`
}

//LearnRSession
type LearnRSession struct {
	ID               int             `json:"ID"`               //ID of this session
	LearnRID         int             `json:"LearnRID"`         //ID of this LearnR
	LearnRName       string          `json:"LearnRName"`       //Name of this LearnR
	TheLearnR        Learnr          `json:"TheLearnR"`        //The actual LearnR
	TheUser          User            `json:"TheUser"`          //Who is the User that sent this LearnR to someone?
	TargetUserNumber string          `json:"TargetUserNumber"` //User this session started to
	Ongoing          bool            `json:"Ongoing"`          //Is this session ongoing? Determined by time
	TextsSent        []LearnRInforms `json:"TextsSent"`        //All the Informs our program sent to User
	UserResponses    []string        `json:"UserResponses"`    //All the text responses sent by the User
	DateCreated      string          `json:"DateCreated"`
	DateUpdated      string          `json:"DateUpdated"`
}

//LearnRInforms
type LearnRInforms struct {
	ID          int    `json:"ID"`         //ID of this Inform
	Name        string `json:"Name"`       //Name of this Inform
	LearnRID    int    `json:"LearnRID"`   //ID of the LearnR this belongs to
	LearnRName  string `json:"LearnRName"` //Name this LearnR belongs to
	Order       int    `json:"Order"`      //The Order in the LearnR this will be
	TheInfo     string `json:"TheInfo"`    //What you want to say to someone
	ShouldWait  bool   `json:"ShouldWait"` //Should this info wait for User Response?
	WaitTime    int    `json:"WaitTime"`   //How much time should User be given to read this text?
	DateCreated string `json:"DateCreated"`
	DateUpdated string `json:"DateUpdated"`
}

//This gets the client to connect to our DB
func connectDB() *mongo.Client {
	//Setup Mongo connection to Atlas Cluster
	theClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("Errored getting mongo client: %v\n", err)
		log.Fatal(err)
	}
	theContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = theClient.Connect(theContext)
	if err != nil {
		fmt.Printf("Errored getting mongo client context: %v\n", err)
		log.Fatal(err)
	}
	//Double check to see if we've connected to the database
	err = theClient.Ping(theContext, readpref.Primary())
	if err != nil {
		fmt.Printf("Errored pinging MongoDB: %v\n", err)
		log.Fatal(err)
	}

	return theClient
}

/* Calls our MongoDB Microservice to add a learnrsession to our DB */
func callAddLearnRSession(newLearnRSession LearnRSession) (bool, string) {
	goodAdd, message := true, ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(newLearnRSession)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", ADDLEARNRSESSIONSURL, payload)
	if err != nil {
		theErr := "There was an error posting LearnrSession: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "Failed response from addLearnRSession: " + strconv.Itoa(resp.StatusCode)
		logWriter(theErr)
		goodAdd, message = false, theErr
	} else if err != nil {
		theErr := "Failed response from addLearnRSession: " + strconv.Itoa(resp.StatusCode) + " " + err.Error()
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from learnRSessionCreate " + err.Error()
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)
	/* 5. Evaluate response in returnedMessage */
	if returnedMessage.SuccOrFail != 0 {
		theErr := ""
		for n := 0; n < len(returnedMessage.TheErr); n++ {
			theErr = theErr + returnedMessage.TheErr[n]
		}
		goodAdd, message = false, theErr
	} else {
		goodAdd, message = true, "LearnrSession successfully added and able to log in"
	}

	return goodAdd, message
}

/* An intended go routine service that calls and manages logging for callAddLearnRSession */
func fastAddLearnRSession(newLearnRSession LearnRSession) {
	goodAdd, theMessage := callAddLearnRSession(newLearnRSession)
	if !goodAdd {
		errMsg := "Could not add this learnRSession: " + strconv.Itoa(newLearnRSession.ID) + "\n" + theMessage
		logWriter(errMsg)
	} else {
		message := "Added this LearnRSession to the DB: " + strconv.Itoa(newLearnRSession.ID)
		logWriter(message)
	}

	wg.Done()
}

/* Calls our CRUD API service to update our LearnRInfo */
func callUpdateLearnrInfo(updatedLearnRInfo LearnrInfo) (bool, string) {
	goodAdd, message := true, ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(updatedLearnRInfo)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", UPDATELEARNRINFOURL, payload)
	if err != nil {
		theErr := "There was an error posting Updated LearnrInfo: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "Failed response from updateLearnRInfo: " + strconv.Itoa(resp.StatusCode)
		logWriter(theErr)
		goodAdd, message = false, theErr
	} else if err != nil {
		theErr := "Failed response from updatedLearnrInfo: " + strconv.Itoa(resp.StatusCode) + " " + err.Error()
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from learnrInfoUpdate " + err.Error()
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)
	/* 5. Evaluate response in returnedMessage */
	if returnedMessage.SuccOrFail != 0 {
		theErr := ""
		for n := 0; n < len(returnedMessage.TheErr); n++ {
			theErr = theErr + returnedMessage.TheErr[n]
		}
		goodAdd, message = false, theErr
	} else {
		goodAdd, message = true, "LearnrInfo successfully updated"
	}

	return goodAdd, message
}

/* A fast GoRoutine Service to get our LearnRInfo and update it with the new Session */
func fastUpdateLearnRInform(newLearnRInfo LearnrInfo) {
	goodAdd, theMessage := callUpdateLearnrInfo(newLearnRInfo)

	if !goodAdd {
		errMsg := "Could not add this learnRSession: " + strconv.Itoa(newLearnRInfo.ID) + "\n" + theMessage
		logWriter(errMsg)
	} else {
		message := "Updated this LearnRInfo in the DB: " + strconv.Itoa(newLearnRInfo.ID)
		logWriter(message)
	}

	wg.Done()
}

/* Gets a random API after calling our random API */
func randomAPICall() (bool, string, int) {
	goodGet, message, finalInt := true, "", 0
	//Call our crudOperations Microservice in order to get our Usernames
	//Create a context for timing out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", GETRANDOMID, nil)
	if err != nil {
		theErr := "There was an error getting Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
		goodGet, message = false, theErr
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		goodGet, message = false, "Wrong response code gotten; failed to create random ID: "+strconv.Itoa(resp.StatusCode)
	} else if err != nil {
		theErr := "Had an error getting good random ID: " + err.Error()
		logWriter(theErr)
		goodGet, message = false, theErr
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
		goodGet, message = false, theErr
	}

	//Marshal the response into a type we can read
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
		RandomID   int      `json:"RandomID"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)

	//Assign our map variable to the map varialbe and see if it's okay
	if returnedMessage.SuccOrFail != 0 {
		errString := ""
		for l := 0; l < len(returnedMessage.TheErr); l++ {
			errString = errString + returnedMessage.TheErr[l]
		}
		goodGet, message = false, errString
	} else {
		finalInt = returnedMessage.RandomID
	}

	return goodGet, message, finalInt
}

/* Calls our CRUD Service to Get our User */
func callGetUser(userID int) (User, bool, string) {
	goodAdd, message := true, ""

	/* 1. Create Context */
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	type UserIDUser struct {
		TheUserID int `json:"TheUserID"`
	}
	theID := UserIDUser{TheUserID: userID}
	theJSONMessage, err := json.Marshal(theID)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", READUSERURL, payload)
	if err != nil {
		theErr := "We had an error with this request: %v\n" + err.Error()
		fmt.Println(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	//defer resp.Body.Close()
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "We had an error with this response: " + strconv.Itoa(resp.StatusCode)
		goodAdd, message = false, theErr
		resp.Body.Close()
		logWriter(theErr)
	} else if err != nil {
		theErr := "We had an error with this response: " + strconv.Itoa(resp.StatusCode)
		goodAdd, message = false, theErr
		resp.Body.Close()
		logWriter(theErr)
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from UserGet " + err.Error()
		goodAdd, message = false, theErr
	}
	type ReturnMessage struct {
		TheErr       []string `json:"TheErr"`
		ResultMsg    []string `json:"ResultMsg"`
		SuccOrFail   int      `json:"SuccOrFail"`
		ReturnedUser User     `json:"ReturnedUser"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)
	/* 5. Evaluate response in returnedMessage */
	if returnedMessage.SuccOrFail != 0 {
		theErr := ""
		for n := 0; n < len(returnedMessage.TheErr); n++ {
			theErr = theErr + returnedMessage.TheErr[n]
		}
		goodAdd, message = false, theErr
	} else {
		goodAdd, message = true, "User successfully got"
	}

	return returnedMessage.ReturnedUser, goodAdd, message
}

//Calls Crud to get our LearnR
func callReadLearnR(theid int) (bool, string, Learnr) {
	goodAdd, message := true, ""

	/* 1. Create Context */
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	type LearnRID struct {
		ID int `json:"ID"`
	}
	theID := LearnRID{ID: theid}
	theJSONMessage, err := json.Marshal(theID)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", READLEARNRURL, payload)
	if err != nil {
		theErr := "We had an error with this request: %v\n" + err.Error()
		fmt.Println(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	//defer resp.Body.Close()
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "We had an error with this response: " + strconv.Itoa(resp.StatusCode)
		goodAdd, message = false, theErr
		resp.Body.Close()
		logWriter(theErr)
	} else if err != nil {
		theErr := "We had an error with this response: " + strconv.Itoa(resp.StatusCode)
		goodAdd, message = false, theErr
		resp.Body.Close()
		logWriter(theErr)
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from UserCreate " + err.Error()
		goodAdd, message = false, theErr
	}
	type ReturnMessage struct {
		TheErr         []string `json:"TheErr"`
		ResultMsg      []string `json:"ResultMsg"`
		SuccOrFail     int      `json:"SuccOrFail"`
		ReturnedLearnR Learnr   `json:"ReturnedLearnR"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)
	/* 5. Evaluate response in returnedMessage */
	if returnedMessage.SuccOrFail != 0 {
		theErr := ""
		for n := 0; n < len(returnedMessage.TheErr); n++ {
			theErr = theErr + returnedMessage.TheErr[n]
		}
		goodAdd, message = false, theErr
	} else {
		goodAdd, message = true, "Learnr successfully updated"
	}

	return goodAdd, message, returnedMessage.ReturnedLearnR
}

//Get our LearnRInfo
func callReadLearnrInfo(theid int) (bool, string, LearnrInfo) {
	goodAdd, message := true, ""

	/* 1. Create Context */
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	type LearnRInfoID struct {
		ID int `json:"ID"`
	}
	theID := LearnRInfoID{ID: theid}
	theJSONMessage, err := json.Marshal(theID)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", READLEARNRINFOURL, payload)
	if err != nil {
		theErr := "We had an error with this request: %v\n" + err.Error()
		fmt.Println(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	//defer resp.Body.Close()
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "We had an error with this response: " + strconv.Itoa(resp.StatusCode)
		goodAdd, message = false, theErr
		resp.Body.Close()
		logWriter(theErr)
	} else if err != nil {
		theErr := "We had an error with this response: " + strconv.Itoa(resp.StatusCode)
		goodAdd, message = false, theErr
		resp.Body.Close()
		logWriter(theErr)
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from UserCreate " + err.Error()
		goodAdd, message = false, theErr
	}
	type ReturnMessage struct {
		TheErr             []string   `json:"TheErr"`
		ResultMsg          []string   `json:"ResultMsg"`
		SuccOrFail         int        `json:"SuccOrFail"`
		ReturnedLearnRInfo LearnrInfo `json:"ReturnedLearnRInfo"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)
	/* 5. Evaluate response in returnedMessage */
	if returnedMessage.SuccOrFail != 0 {
		theErr := ""
		for n := 0; n < len(returnedMessage.TheErr); n++ {
			theErr = theErr + returnedMessage.TheErr[n]
		}
		goodAdd, message = false, theErr
	} else {
		goodAdd, message = true, "LearnrInfo successfully gotten"
	}

	return goodAdd, message, returnedMessage.ReturnedLearnRInfo
}
