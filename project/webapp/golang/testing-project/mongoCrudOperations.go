package main

import (
	"context"
	"encoding/hex"
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

//Mongo DB Declarations
var mongoClient *mongo.Client

var theContext context.Context
var mongoURI string //Connection string loaded

/* DEFINE URL FOR LOCALHOST */
//LearnROrg
const READLEARNRORGURL string = "http://localhost:4000/getLearnOrg"
const UPDATELEARNRORGURL string = "http://localhost:4000/updateLearnOrg"
const DELETELEARNRORGURL string = "http://localhost:4000/deleteLearnOrg"

//LearnR
const ADDLEARNRURL string = "http://localhost:4000/addLearnR"
const READLEARNRURL string = "http://localhost:4000/getLearnR"
const UPDATELEARNRURL string = "http://localhost:4000/updateLearnR"
const DELETELEARNRURL string = "http://localhost:4000/deleteLearnR"

//LearnRInfo
const ADDLEARNRINFOURL string = "http://localhost:4000/addLearnrInfo"
const READLEARNRINFOURL string = "http://localhost:4000/getLearnrInfo"
const UPDATELEARNRINFOURL string = "http://localhost:4000/updateLearnrInfo"
const DELETELEARNRINFOURL string = "http://localhost:4000/deleteLearnrInfo"

//LearnRSession
const ADDLEARNRSESSIONSURL string = "http://localhost:4000/addLearnRSession"
const READLEARNRSESSIONSURL string = "http://localhost:4000/getLearnRSession"
const UPDATELEARNRSESSIONSURL string = "http://localhost:4000/updateLearnRSession"
const DELETELEARNRSESSIONSURL string = "http://localhost:4000/deleteLearnRSession"

//LearnRInform
const ADDLEARNRINFORMURL string = "http://localhost:4000/addLearnRInforms"
const READLEARNRINFORMURL string = "http://localhost:4000/getLearnRInforms"
const UPDATELEARNRINFORMURL string = "http://localhost:4000/updateLearnRInforms"
const DELETELEARNRINFORMURL string = "http://localhost:4000/deleteLearnRInforms"

/* App/Data type declarations for our application */
// Desc: This person uses our app
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

/* Calls our CRUD service to add our new User */
func callAddUser(newUser User) (bool, string) {
	goodAdd, message := true, ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(newUser)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", ADDUSERURL, payload)
	if err != nil {
		theErr := "There was an error posting User: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "Failed response from addUser: " + strconv.Itoa(resp.StatusCode)
		logWriter(theErr)
		goodAdd, message = false, theErr
	} else if err != nil {
		theErr := "Failed response from addUser: " + strconv.Itoa(resp.StatusCode) + " " + err.Error()
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from UserCreate " + err.Error()
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
		goodAdd, message = true, "User successfully added and able to log in"
	}

	return goodAdd, message
}

/* Calls our CRUD Service to update our User */
func callUpdateUser(updatedUser User) (bool, string) {
	goodAdd, message := true, ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(updatedUser)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", UPDATEURL, payload)
	if err != nil {
		theErr := "There was an error posting Updated User: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "Failed response from updateUser: " + strconv.Itoa(resp.StatusCode)
		logWriter(theErr)
		goodAdd, message = false, theErr
	} else if err != nil {
		theErr := "Failed response from updatedUser: " + strconv.Itoa(resp.StatusCode) + " " + err.Error()
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from userUpdate " + err.Error()
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
		goodAdd, message = true, "User successfully updated"
	}

	return goodAdd, message
}

/* Calls our CRUD service to see if this User can login with the passed Username/Password */
func callUserLogin(username string, password string) (bool, string, User) {
	goodLogin, message, theUser := true, "", User{}

	/* Call the API */
	type LoginData struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}
	loginData := LoginData{Username: username, Password: hex.EncodeToString([]byte(password))}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(loginData)
	if err != nil {
		theErr := "Error marshaling JSON: " + err.Error()
		goodLogin, message, theUser = false, theErr, User{}
		logWriter(theErr)
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", GETUSERLOGIN, payload)
	if err != nil {
		theErr := "Error with request: " + err.Error()
		goodLogin, message, theUser = false, theErr, User{}
		logWriter(theErr)
	}
	//req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "Error getting response code: " + strconv.Itoa(resp.StatusCode)
		goodLogin, message, theUser = false, theErr, User{}
		logWriter(theErr)
	} else if err != nil {
		theErr := "Error with response: " + err.Error()
		goodLogin, message, theUser = false, theErr, User{}
		logWriter(theErr)
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from UserCreate " + err.Error()
		goodLogin, message, theUser = false, theErr, User{}
		logWriter(theErr)
	}
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
		TheUser    User     `json:"TheUser"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)
	/* 5. Evaluate response in returnedMessage for testing */
	if returnedMessage.SuccOrFail != 0 {
		theMessage := ""
		for n := 0; n < len(returnedMessage.TheErr); n++ {
			theMessage = theMessage + returnedMessage.TheErr[n]
		}
		goodLogin, message, theUser = false, theMessage, User{}
		logWriter(theMessage)
	} else {
		message = "User found"
		theUser = returnedMessage.TheUser
	}

	return goodLogin, message, theUser
}

/* Calls our CRUD service to add a new LearnR Organization */
func calladdLearnOrg(newLearnROrg LearnrOrg) (bool, string) {
	goodAdd, message := true, ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(newLearnROrg)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", ADDLEARNRORGURL, payload)
	if err != nil {
		theErr := "There was an error posting User: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "Failed response from addLearnROrg: " + strconv.Itoa(resp.StatusCode)
		logWriter(theErr)
		goodAdd, message = false, theErr
	} else if err != nil {
		theErr := "Failed response from LearnROrg: " + strconv.Itoa(resp.StatusCode) + " " + err.Error()
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from learnROrg Create " + err.Error()
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
		goodAdd, message = true, "LearnR Org successfully created"
	}

	return goodAdd, message
}

/* CRUD Operations for LearnR */
func callAddLearnR(newLearnR Learnr) (bool, string) {
	goodAdd, message := true, ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(newLearnR)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", ADDLEARNRURL, payload)
	if err != nil {
		theErr := "There was an error posting Learnr: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "Failed response from addLearnr: " + strconv.Itoa(resp.StatusCode)
		logWriter(theErr)
		goodAdd, message = false, theErr
	} else if err != nil {
		theErr := "Failed response from addLearnr: " + strconv.Itoa(resp.StatusCode) + " " + err.Error()
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from learnRCreate " + err.Error()
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
		goodAdd, message = true, "Learnr successfully added and able to log in"
	}

	return goodAdd, message
}

func callUpdateLearnR(updatedLearnr Learnr) (bool, string) {
	goodAdd, message := true, ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(updatedLearnr)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", UPDATELEARNRURL, payload)
	if err != nil {
		theErr := "There was an error posting Updated Learnr: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "Failed response from updateLearnR: " + strconv.Itoa(resp.StatusCode)
		logWriter(theErr)
		goodAdd, message = false, theErr
	} else if err != nil {
		theErr := "Failed response from updatedLearnr: " + strconv.Itoa(resp.StatusCode) + " " + err.Error()
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from learnrUpdate " + err.Error()
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
		goodAdd, message = true, "Learnr successfully updated"
	}

	return goodAdd, message
}

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

func callDeleteLearnR(theid int) (bool, string) {
	goodAdd, message := true, ""

	/* 1. Create Context */
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	type LearnRDelete struct {
		ID int `json:"ID"`
	}
	theID := LearnRDelete{ID: theid}
	theJSONMessage, err := json.Marshal(theID)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", DELETELEARNRURL, payload)
	if err != nil {
		theErr := "We had an error with this request: %v\n" + err.Error()
		fmt.Println(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
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
		theErr := "There was an error reading response from learnRDelete " + err.Error()
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
		goodAdd, message = true, "Learnr successfully deleted"
	}

	return goodAdd, message
}

/* CRUD Operations for LearnrInfo */

func callAddLearnrInfo(newLearnRInfo LearnrInfo) (bool, string) {
	goodAdd, message := true, ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(newLearnRInfo)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", ADDLEARNRINFOURL, payload)
	if err != nil {
		theErr := "There was an error posting Learnr: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "Failed response from addLearnRInfo: " + strconv.Itoa(resp.StatusCode)
		logWriter(theErr)
		goodAdd, message = false, theErr
	} else if err != nil {
		theErr := "Failed response from addLearnRInfo: " + strconv.Itoa(resp.StatusCode) + " " + err.Error()
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from learnRInfoCreate " + err.Error()
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
		goodAdd, message = true, "LearnrInfo successfully added and able to log in"
	}

	return goodAdd, message
}

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

func callDeleteLearnrInfo(theid int) (bool, string) {
	goodAdd, message := true, ""

	/* 1. Create Context */
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	type LearnRInfoDelete struct {
		ID int `json:"ID"`
	}
	theID := LearnRInfoDelete{ID: theid}
	theJSONMessage, err := json.Marshal(theID)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", DELETELEARNRINFOURL, payload)
	if err != nil {
		theErr := "We had an error with this request: %v\n" + err.Error()
		fmt.Println(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
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
		theErr := "There was an error reading response from learnRInfoDelete " + err.Error()
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
		goodAdd, message = true, "LearnrInfo successfully deleted"
	}

	return goodAdd, message
}

/* CRUD operations for LearnRSession */

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

func callUpdateLearnRSession(updatedLearnRSession LearnRSession) (bool, string) {
	goodAdd, message := true, ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(updatedLearnRSession)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", UPDATELEARNRSESSIONSURL, payload)
	if err != nil {
		theErr := "There was an error posting Updated LearnrSession: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "Failed response from updateLearnRSession: " + strconv.Itoa(resp.StatusCode)
		logWriter(theErr)
		goodAdd, message = false, theErr
	} else if err != nil {
		theErr := "Failed response from updatedLearnrSession: " + strconv.Itoa(resp.StatusCode) + " " + err.Error()
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from learnrSessionUpdate " + err.Error()
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
		goodAdd, message = true, "LearnrSession successfully updated"
	}

	return goodAdd, message
}

func callReadLearnRSession(theid int) (bool, string, LearnRSession) {
	goodAdd, message := true, ""

	/* 1. Create Context */
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	type LearnRSessionID struct {
		ID int `json:"ID"`
	}
	theID := LearnRSessionID{ID: theid}
	theJSONMessage, err := json.Marshal(theID)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", READLEARNRSESSIONSURL, payload)
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
		theErr := "There was an error reading response from LearnRSessionRead" + err.Error()
		goodAdd, message = false, theErr
	}
	type ReturnMessage struct {
		TheErr          []string      `json:"TheErr"`
		ResultMsg       []string      `json:"ResultMsg"`
		SuccOrFail      int           `json:"SuccOrFail"`
		ReturnedSession LearnRSession `json:"ReturnedSession"`
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
		goodAdd, message = true, "LearnrSession successfully gotten"
	}

	return goodAdd, message, returnedMessage.ReturnedSession
}

func callDeleteLearnRSession(theid int) (bool, string) {
	goodAdd, message := true, ""

	/* 1. Create Context */
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	type LearnRSessionDelete struct {
		ID int `json:"ID"`
	}
	theID := LearnRSessionDelete{ID: theid}
	theJSONMessage, err := json.Marshal(theID)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", DELETELEARNRSESSIONSURL, payload)
	if err != nil {
		theErr := "We had an error with this request: %v\n" + err.Error()
		fmt.Println(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
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
		theErr := "There was an error reading response from learnRSessionDelete " + err.Error()
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
		goodAdd, message = true, "LearnrSession successfully deleted"
	}

	return goodAdd, message
}

/* CRUD operations for LearnRInform */

func callAddLearnRInform(newLearnRInform LearnRInforms) (bool, string) {
	goodAdd, message := true, ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(newLearnRInform)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", ADDLEARNRINFORMURL, payload)
	if err != nil {
		theErr := "There was an error posting LearnRInform: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "Failed response from addLearnRInform: " + strconv.Itoa(resp.StatusCode)
		logWriter(theErr)
		goodAdd, message = false, theErr
	} else if err != nil {
		theErr := "Failed response from addLearnRInform: " + strconv.Itoa(resp.StatusCode) + " " + err.Error()
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from learnRInformCreate" + err.Error()
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
		goodAdd, message = true, "LearnrInform successfully added and able to log in"
	}

	return goodAdd, message
}

func callUpdateLearnRInform(updatedLearnRInform LearnRInforms) (bool, string) {
	goodAdd, message := true, ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(updatedLearnRInform)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", UPDATELEARNRINFORMURL, payload)
	if err != nil {
		theErr := "There was an error posting Updated LearnRInform: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "Failed response from updateLearnRInform: " + strconv.Itoa(resp.StatusCode)
		logWriter(theErr)
		goodAdd, message = false, theErr
	} else if err != nil {
		theErr := "Failed response from updatedLearnrInform: " + strconv.Itoa(resp.StatusCode) + " " + err.Error()
		logWriter(theErr)
		goodAdd, message = false, theErr
	}
	//Declare message we expect to see returned
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error reading response from learnrInformUpdate " + err.Error()
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
		goodAdd, message = true, "LearnrInform successfully updated"
	}

	return goodAdd, message
}

func callReadLearnRInform(theid int) (bool, string, LearnRInforms) {
	goodAdd, message := true, ""

	/* 1. Create Context */
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	type LearnRInformID struct {
		ID int `json:"ID"`
	}
	theID := LearnRInformID{ID: theid}
	theJSONMessage, err := json.Marshal(theID)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", READLEARNRINFORMURL, payload)
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
		theErr := "There was an error reading response from LearnRInform" + err.Error()
		goodAdd, message = false, theErr
	}
	type ReturnMessage struct {
		TheErr               []string      `json:"TheErr"`
		ResultMsg            []string      `json:"ResultMsg"`
		SuccOrFail           int           `json:"SuccOrFail"`
		ReturnedLearnRInform LearnRInforms `json:"ReturnedLearnRInform"`
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
		goodAdd, message = true, "LearnrInform successfully gotten"
	}

	return goodAdd, message, returnedMessage.ReturnedLearnRInform
}

func callDeleteLearnRInform(theid int) (bool, string) {
	goodAdd, message := true, ""

	/* 1. Create Context */
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	type LearnRInformsDelete struct {
		ID int `json:"ID"`
	}
	theID := LearnRInformsDelete{ID: theid}
	theJSONMessage, err := json.Marshal(theID)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
		goodAdd, message = false, err.Error()
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", DELETELEARNRINFORMURL, payload)
	if err != nil {
		theErr := "We had an error with this request: %v\n" + err.Error()
		fmt.Println(theErr)
		goodAdd, message = false, theErr
	}
	req.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
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
		theErr := "There was an error reading response from learnRInformsDelete " + err.Error()
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
		goodAdd, message = true, "LearnrInform successfully deleted"
	}

	return goodAdd, message
}
