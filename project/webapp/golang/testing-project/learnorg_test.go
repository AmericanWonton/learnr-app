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
	"testing"
	"time"
)

/* DEFINE URL FOR LOCALHOST */
const READLEARNRORGURL string = "http://localhost:4000/getLearnOrg"
const UPDATELEARNRORGURL string = "http://localhost:4000/updateLearnOrg"
const DELETELEARNRORGURL string = "http://localhost:4000/deleteLearnOrg"

//LearnOrg CRUD Create
type LearnOrgCrudCreate struct {
	LearnOrg            LearnrOrg
	ExpectedNum         int
	ExpectedStringArray []string
}

var learnOrgCrudCreateResults []LearnOrgCrudCreate

//LearnOrg Crud Read
type LearnOrgCrudRead struct {
	TheLearnOrgID       int
	ExpectedNum         int
	ExpectedStringArray []string
}

var learnOrgCrudReadResults []LearnOrgCrudRead

//LearnOrg Crud Update
type LearnOrgCrudUpdate struct {
	TheLearnOrg         LearnrOrg
	ExpectedNum         int
	ExpectedStringArray []string
}

var learnOrgCrudUpdateResults []LearnOrgCrudUpdate

//LearnOrg CRUD Delete
type LearnOrgCrudDelete struct {
	TheLearnOrgID       int
	ExpectedNum         int
	ExpectedStringArray []string
}

var learnOrgCrudDeleteResults []LearnOrgCrudDelete

//This creates our Crud Testing cases for Creating LearnOrgs
func createCreateLearnOrgCrud() {
	theTimeNow := time.Now() //Used for creating time later
	//Good LearnOrg Crud Create
	learnOrgCrudCreateResults = append(learnOrgCrudCreateResults, LearnOrgCrudCreate{LearnrOrg{
		OrgID:       1111,
		Name:        "TestOrg",
		OrgGoals:    []string{"Being Super Cool", "Being super awesome"},
		UserList:    []int{1111, 2222},
		AdminList:   []int{1111},
		LearnrList:  []int{4567},
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}, 0, []string{"LearnOrg successfully added in addlearnOrg"}})
	//Empty LearnOrg Crud
	learnOrgCrudCreateResults = append(learnOrgCrudCreateResults, LearnOrgCrudCreate{LearnrOrg{}, 1,
		[]string{"Error adding LearnOrg in addLearnOrg", "Error adding LarnOrg in addLarnOrg in crudoperations API"}})
	//LearnOrg with Zero value
	learnOrgCrudCreateResults = append(learnOrgCrudCreateResults, LearnOrgCrudCreate{LearnrOrg{OrgID: 0}, 1,
		[]string{"Error adding LearnOrg in addLearnOrg", "Error reading the request"}})
	//LearnOrg with negative OrgID value
	learnOrgCrudCreateResults = append(learnOrgCrudCreateResults, LearnOrgCrudCreate{LearnrOrg{OrgID: -1}, 1,
		[]string{"Error adding LearnOrg in addLearnOrg", "Error reading the request"}})
}

//This creates our CRUD Testing cases for Reading LearnOrgs
func createLearnOrgReadCrud() {
	//Good LearnOrg Crud Read
	learnOrgCrudReadResults = append(learnOrgCrudReadResults, LearnOrgCrudRead{1111, 0, []string{"LearnOrg successfully read in getLearnOrg"}})
	//Bad LearnOrg CRUD Read
	learnOrgCrudReadResults = append(learnOrgCrudReadResults, LearnOrgCrudRead{0, 1,
		[]string{"Error adding LearnOrg in addLearnOrg", "Error reading the request"}})
	//Not seen OrgID
	learnOrgCrudReadResults = append(learnOrgCrudReadResults, LearnOrgCrudRead{4000000, 1,
		[]string{"Error adding LearnOrg in addLearnOrg", "Error reading the request"}})
	//Another not seen OrgID
	learnOrgCrudReadResults = append(learnOrgCrudReadResults, LearnOrgCrudRead{-1, 1,
		[]string{"Error adding LearnOrg in addLearnOrg", "Error reading the request"}})
}

//This creates our CRUD Update cases for Updating LearnOrgs
func createLearnOrgUpdateCrud() {
	theTimeNow := time.Now() //Used for creating time later
	//Good LearnOrg Crud Create
	learnOrgCrudUpdateResults = append(learnOrgCrudUpdateResults, LearnOrgCrudUpdate{LearnrOrg{
		OrgID:       1111,
		Name:        "TestOrg Revised",
		OrgGoals:    []string{"Being Super NOT Cool", "Being super NOT awesome"},
		UserList:    []int{1111, 2222, 3333},
		AdminList:   []int{1111},
		LearnrList:  []int{4567, 7859},
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}, 0, []string{"LearnOrg successfully updated in addLearnOrg"}})
	//Bad Non-Existent OrgID
	learnOrgCrudUpdateResults = append(learnOrgCrudUpdateResults, LearnOrgCrudUpdate{LearnrOrg{
		OrgID:       400000,
		Name:        "TestOrg Revised",
		OrgGoals:    []string{"Being Super NOT Cool", "Being super NOT awesome"},
		UserList:    []int{1111, 2222, 3333},
		AdminList:   []int{1111},
		LearnrList:  []int{4567, 7859},
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}, 1, []string{"Error adding LearnOrg in addUser", "Error reading the request"}})
	//Bad Empty LearnOrg Crud
	learnOrgCrudUpdateResults = append(learnOrgCrudUpdateResults, LearnOrgCrudUpdate{LearnrOrg{}, 1,
		[]string{"Error adding LearnOrg in addLearnOrg", "Error reading the request"}})
}

//This creates our CRUD Delete Cases for deleting LearnOrgs
func createLearnOrgDeleteCrud() {
	//Good LearnOrg Crud Read
	learnOrgCrudDeleteResults = append(learnOrgCrudDeleteResults, LearnOrgCrudDelete{1111, 0,
		[]string{"LearnOrg successfully deleted in deleteLearnOrg"}})
	//Bad LearnOrg CRUD Read
	learnOrgCrudDeleteResults = append(learnOrgCrudDeleteResults, LearnOrgCrudDelete{0, 1,
		[]string{"Error deleting LearnOrg in deleteLearnOrg", "Error reading the request"}})
	//Not seen OrgID
	learnOrgCrudDeleteResults = append(learnOrgCrudDeleteResults, LearnOrgCrudDelete{4000000, 1,
		[]string{"Error deleting LearnOrg in deleteLearnOrg", "Error reading the request"}})
	//Another not seen OrgID
	learnOrgCrudDeleteResults = append(learnOrgCrudDeleteResults, LearnOrgCrudDelete{-1, 1,
		[]string{"Error deleting LearnOrg in deleteLearnOrg", "Error reading the request"}})
}

/* Here are our Test calls to the CRUD APIs */
//Add a LearnROrg
func TestLearnROrgAdd(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range learnOrgCrudCreateResults {
		/* start listener */
		/* 1. Create Context */
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		/* 2. Marshal test case to JSON expect */
		theJSONMessage, err := json.Marshal(test.LearnOrg)
		if err != nil {
			fmt.Println(err)
			logWriter(err.Error())
			log.Fatal(err)
		}
		/* 3. Create Post to JSON */
		payload := strings.NewReader(string(theJSONMessage))
		req, err := http.NewRequest("POST", ADDLEARNRORGURL, payload)
		if err != nil {
			log.Fatal(err)
		}
		//req.Header.Add("Content-Type", "text/plain")
		req.Header.Add("Content-Type", "application/json")
		/* 4. Get response from Post */
		resp, err := http.DefaultClient.Do(req.WithContext(ctx))
		if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
			theRespCode := strconv.Itoa(resp.StatusCode)
			t.Fatal("We have the wrong response code: " + theRespCode)
		} else if err != nil {
			t.Fatal("Had an error creating response: " + err.Error())
		}
		//Declare message we expect to see returned
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			theErr := "There was an error reading response from UserCreate " + err.Error()
			t.Fatal(theErr)
		}
		type ReturnMessage struct {
			TheErr     []string `json:"TheErr"`
			ResultMsg  []string `json:"ResultMsg"`
			SuccOrFail int      `json:"SuccOrFail"`
		}
		var returnedMessage ReturnMessage
		json.Unmarshal(body, &returnedMessage)
		/* 5. Evaluate response in returnedMessage for testing */
		if test.ExpectedNum != returnedMessage.SuccOrFail {
			t.Fatal("Wrong num recieved on testcase " + strconv.Itoa(testNum) +
				" :" + strconv.Itoa(returnedMessage.SuccOrFail) + " Expected: " + strconv.Itoa(test.ExpectedNum))
		}
		/* Maybe we can test the strings at some point... */
		testNum = testNum + 1 //Increment this number for testing
	}
}

//Test for updating LearnR Orgs
func TestLearnROrgUpdate(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range learnOrgCrudUpdateResults {
		/* start listener */
		/* 1. Create Context */
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		/* 2. Marshal test case to JSON expect */
		theJSONMessage, err := json.Marshal(test.TheLearnOrg)
		if err != nil {
			fmt.Println(err)
			logWriter(err.Error())
			log.Fatal(err)
			t.Fatal("Could not Marshal JSON: " + err.Error())
		}
		/* 3. Create Post to JSON */
		payload := strings.NewReader(string(theJSONMessage))
		req, err := http.NewRequest("POST", UPDATELEARNRORGURL, payload)
		if err != nil {
			log.Fatal(err)
			t.Fatal("Had an issue creating a request: " + err.Error())
		}
		//req.Header.Add("Content-Type", "text/plain")
		req.Header.Add("Content-Type", "application/json")
		/* 4. Get response from Post */
		resp, err := http.DefaultClient.Do(req.WithContext(ctx))
		if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
			theRespCode := strconv.Itoa(resp.StatusCode)
			t.Fatal("We have the wrong response code: " + theRespCode)
		} else if err != nil {
			t.Fatal("Had an error creating response: " + err.Error())
		}
		//Declare message we expect to see returned
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			theErr := "There was an error reading response from UserCreate " + err.Error()
			t.Fatal(theErr)
		}
		type ReturnMessage struct {
			TheErr     []string `json:"TheErr"`
			ResultMsg  []string `json:"ResultMsg"`
			SuccOrFail int      `json:"SuccOrFail"`
		}
		var returnedMessage ReturnMessage
		json.Unmarshal(body, &returnedMessage)
		/* 5. Evaluate response in returnedMessage for testing */
		if test.ExpectedNum != returnedMessage.SuccOrFail {
			t.Fatal("Wrong num recieved on testcase " + strconv.Itoa(testNum) +
				" :" + strconv.Itoa(returnedMessage.SuccOrFail) + " Expected: " + strconv.Itoa(test.ExpectedNum))
		}
		/* Maybe we can test the strings at some point... */
		testNum = testNum + 1 //Increment this number for testing
	}
}

//Test for Reading LearnR Orgs
func TestLearnROrgRead(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range learnOrgCrudReadResults {
		/* 1. Create Context */
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		/* 2. Marshal test case to JSON expect */
		type LearnOrgID struct {
			TheLearnOrgID int `json:"TheLearnOrgID"`
		}
		aLearnOrgID := LearnOrgID{TheLearnOrgID: test.TheLearnOrgID}
		theJSONMessage, err := json.Marshal(aLearnOrgID)
		if err != nil {
			fmt.Println(err)
			logWriter(err.Error())
			log.Fatal(err)
			t.Fatal("Could not marshal JSON")
		}
		/* 3. Create Post to JSON */
		payload := strings.NewReader(string(theJSONMessage))
		req, err := http.NewRequest("POST", READLEARNRORGURL, payload)
		if err != nil {
			log.Fatal(err)
			t.Fatal("Had an error making request: " + err.Error())
		}
		req.Header.Add("Content-Type", "application/json")
		/* 4. Get response from Post */
		resp, err := http.DefaultClient.Do(req.WithContext(ctx))
		//defer resp.Body.Close()
		if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
			resp.Body.Close()
			theRespCode := strconv.Itoa(resp.StatusCode)
			t.Fatal("We have the wrong response code: " + theRespCode)
			return
		} else if err != nil {
			resp.Body.Close()
			t.Fatal("Had an error creating response: " + err.Error())
			return
		}
		//Declare message we expect to see returned
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			theErr := "There was an error reading response from UserCreate " + err.Error()
			t.Fatal(theErr)
		}
		type ReturnMessage struct {
			TheErr     []string `json:"TheErr"`
			ResultMsg  []string `json:"ResultMsg"`
			SuccOrFail int      `json:"SuccOrFail"`
		}
		var returnedMessage ReturnMessage
		json.Unmarshal(body, &returnedMessage)
		/* 5. Evaluate response in returnedMessage for testing */
		if test.ExpectedNum != returnedMessage.SuccOrFail {
			fmt.Printf("We here test un-expected\n")
			t.Fatal("Wrong num recieved on testcase " + strconv.Itoa(testNum) +
				" :" + strconv.Itoa(returnedMessage.SuccOrFail) + " Expected: " + strconv.Itoa(test.ExpectedNum))
		}
		/* Maybe we can test the strings at some point... */
		testNum = testNum + 1 //Increment this number for testing
	}
}

//Test for getting all LearnR Orgs
func TestGetAllLearnROrgNames(t *testing.T) {
	//Call our crudOperations Microservice in order to get our Usernames
	//Create a context for timing out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", GETALLLEARNRORGURL, nil)
	if err != nil {
		theErr := "There was an error getting Usernames in loadUsernames: " + err.Error()
		t.Fatal(theErr)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theRespCode := strconv.Itoa(resp.StatusCode)
		t.Fatal("We have the wrong response code: " + theRespCode)
		return
	} else if err != nil {
		t.Fatal("Had an error creating response: " + err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for Usernames in loadUsernames: " + err.Error()
		t.Fatal(theErr)
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
		t.Fatal("Had an error getting map: " + errString)
	}
}

//Test for Deleting LearnR Orgs
func TestLearnROrgsDelete(t *testing.T) {
	time.Sleep(4 * time.Second) //Might needed for CRUD updating
	testNum := 0                //Used for incrementing
	for _, test := range learnOrgCrudDeleteResults {
		/* 1. Create Context */
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		/* 2. Marshal test case to JSON expect */
		type OrgDelete struct {
			OrgID int `json:"OrgID"`
		}
		orgID := OrgDelete{OrgID: test.TheLearnOrgID}
		theJSONMessage, err := json.Marshal(orgID)
		if err != nil {
			fmt.Println(err)
			logWriter(err.Error())
			log.Fatal(err)
			t.Fatal("Could not marshal JSON")
		}
		/* 3. Create Post to JSON */
		payload := strings.NewReader(string(theJSONMessage))
		req, err := http.NewRequest("POST", DELETELEARNRORGURL, payload)
		if err != nil {
			log.Fatal(err)
			t.Fatal("Had an error making request: " + err.Error())
		}
		req.Header.Add("Content-Type", "application/json")
		/* 4. Get response from Post */
		resp, err := http.DefaultClient.Do(req.WithContext(ctx))
		if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
			resp.Body.Close()
			theRespCode := strconv.Itoa(resp.StatusCode)
			t.Fatal("We have the wrong response code: " + theRespCode)
			return
		} else if err != nil {
			resp.Body.Close()
			t.Fatal("Had an error creating response: " + err.Error())
			return
		}
		defer resp.Body.Close()
		//Declare message we expect to see returned
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			theErr := "There was an error reading response from UserCreate " + err.Error()
			t.Fatal(theErr)
		}
		type ReturnMessage struct {
			TheErr     []string `json:"TheErr"`
			ResultMsg  []string `json:"ResultMsg"`
			SuccOrFail int      `json:"SuccOrFail"`
		}
		var returnedMessage ReturnMessage
		json.Unmarshal(body, &returnedMessage)
		/* 5. Evaluate response in returnedMessage for testing */
		if test.ExpectedNum != returnedMessage.SuccOrFail {
			fmt.Printf("We here testexpected\n")
			t.Fatal("Wrong num recieved on testcase " + strconv.Itoa(testNum) +
				" :" + strconv.Itoa(returnedMessage.SuccOrFail) + " Expected: " + strconv.Itoa(test.ExpectedNum))
		}
		/* Maybe we can test the strings at some point... */
		testNum = testNum + 1 //Increment this number for testing
	}
}
