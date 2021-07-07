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
)

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
	req, err := http.NewRequest("POST", mongoCrudURL+"/addLearnRSession", payload)
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
		fmt.Println(errMsg)
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
	req, err := http.NewRequest("POST", mongoCrudURL+"/updateLearnrInfo", payload)
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
		fmt.Println(errMsg)
	} else {
		message := "Updated this LearnRInfo in the DB: " + strconv.Itoa(newLearnRInfo.ID)
		logWriter(message)
	}

	wg.Done()
}

/* Gets a random API after calling our random API */
func randomAPICall() (bool, string, int) {
	goodGet, message, finalInt := true, "", 0
	/* Keeping this in here until we can figure out how to do GET properly... */
	type LoginData struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}
	theID := LoginData{Username: "tESTUsername", Password: "TestPassword"}
	theJSONMessage, err := json.Marshal(theID)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
	}
	//Create a context for timing out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	payload := strings.NewReader(string(theJSONMessage)) //Debug
	req, err := http.NewRequest("POST", mongoCrudURL+"/randomIDCreationAPI", payload)
	if err != nil {
		theErr := "There was an error getting Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
		goodGet, message = false, theErr
	}
	defer req.Body.Close()

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
	req, err := http.NewRequest("POST", mongoCrudURL+"/getUser", payload)
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
	req, err := http.NewRequest("POST", mongoCrudURL+"/getLearnR", payload)
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
	req, err := http.NewRequest("POST", mongoCrudURL+"/getLearnrInfo", payload)
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

/* Test Calls */
func testPingCrud() {
	fmt.Printf("DEBUG: We are calling the testPing in Crud API first\n")
	//Call our crudOperations Microservice in order to get our Usernames
	//Create a context for timing out
	type LoginData struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}
	theID := LoginData{Username: "tESTUsername", Password: "TestPassword"}
	theJSONMessage, err := json.Marshal(theID)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Printf("DEBUG: Making a request to here...%v\n", mongoCrudURL+"/testPing")
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", mongoCrudURL+"/randomIDCreationAPI", payload)
	if err != nil {
		theErr := "There was an error getting Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
	}
	defer req.Body.Close()

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
	} else if err != nil {
		theErr := "Had an error getting good random ID: " + err.Error()
		logWriter(theErr)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
	}
	type ReturnMessage struct {
		TheErr          []string        `json:"TheErr"`
		ResultMsg       []string        `json:"ResultMsg"`
		SuccOrFail      int             `json:"SuccOrFail"`
		ReturnedUserMap map[string]bool `json:"ReturnedUserMap"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)

	fmt.Printf("DEBUG: Here is our returned Message succ: %v\n", returnedMessage.SuccOrFail)
}
