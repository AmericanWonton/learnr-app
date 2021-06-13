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
	"testing"
	"time"
)

/* DEFINE CRUD POSTS URLS FOR LOCALHOST */
const READUSERURL string = "http://localhost:4000/getUser"
const DELETEURL string = "http://localhost:4000/deleteUser"

//UserCrud Create
type UserCrudCreate struct {
	TheUser             User
	ExpectedNum         int
	ExpectedStringArray []string
}

var userCrudCreateResults []UserCrudCreate

//UserCrud Read
type UserCrudRead struct {
	TheUserID           int
	ExpectedNum         int
	ExpectedStringArray []string
}

var userCrudReadResults []UserCrudRead

//UserCrud Update
type UserCrudUpdate struct {
	TheUser             User
	ExpectedNum         int
	ExpectedStringArray []string
}

var userCrudUpdateResults []UserCrudUpdate

//UserCRUD Delete
type UserCrudDelete struct {
	TheUserID           int
	ExpectedNum         int
	ExpectedStringArray []string
}

var userCrudDeleteResults []UserCrudDelete

//User Login
type UserLoginTest struct {
	TheUsername         string
	ThePassword         string
	ExpectedNum         int
	ExpectedStringArray []string
}

var userLoginResults []UserLoginTest

//This creates our Crud Testing cases for Creating Users
func createCreateUserCrud() {
	theTimeNow := time.Now() //Used for creating time later
	//Good User Crud Create
	userCrudCreateResults = append(userCrudCreateResults, UserCrudCreate{User{
		UserName:    "TestUsername",
		Password:    hex.EncodeToString([]byte("testpword")),
		Firstname:   "Test",
		Lastname:    "User",
		PhoneNums:   []string{"13143228594"},
		UserID:      1111,
		Email:       []string{"jbkeller0303@gmail.com"},
		Whoare:      "I am a test User and I like to write tests",
		AdminOrgs:   []int{1111},
		OrgMember:   []int{1111},
		Banned:      false,
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}, 0, []string{"User successfully added in addUser"}})
	//Empty User Crud
	userCrudCreateResults = append(userCrudCreateResults, UserCrudCreate{User{}, 1,
		[]string{"Error adding User in addUser", "Error reading the request"}})
	//User with Zero value
	userCrudCreateResults = append(userCrudCreateResults, UserCrudCreate{User{UserID: 0}, 1,
		[]string{"Error adding User in addUser", "Error reading the request"}})
	//User with negative UserID value
	userCrudCreateResults = append(userCrudCreateResults, UserCrudCreate{User{UserID: -1}, 1,
		[]string{"Error adding User in addUser", "Error reading the request"}})
}

//This creates our CRUD Testing cases for Reading Users
func createUserReadCrud() {
	//Good User Crud Read
	userCrudReadResults = append(userCrudReadResults, UserCrudRead{1111, 0, []string{"User successfully read in getUser"}})
	//Bad User CRUD Read
	userCrudReadResults = append(userCrudReadResults, UserCrudRead{0, 1,
		[]string{"Error adding User in addUser", "Error reading the request"}})
	//Not seen UserID
	userCrudReadResults = append(userCrudReadResults, UserCrudRead{4000000, 1,
		[]string{"Error adding User in addUser", "Error reading the request"}})
	//Another not seen UserID
	userCrudReadResults = append(userCrudReadResults, UserCrudRead{-1, 1,
		[]string{"Error adding User in addUser", "Error reading the request"}})
}

//This creates our CRUD Update cases for Updating Users
func createUserUpdateCrud() {
	theTimeNow := time.Now() //Used for creating time later
	//Good User Crud Create
	userCrudUpdateResults = append(userCrudUpdateResults, UserCrudUpdate{User{
		UserName:    "StickyMicky",
		Password:    hex.EncodeToString([]byte("testpword")),
		Firstname:   "Test",
		Lastname:    "User",
		PhoneNums:   []string{"13143228594"},
		UserID:      1111,
		Email:       []string{"jbkeller0303@gmail.com"},
		Whoare:      "I am a test User and I like to write tests",
		AdminOrgs:   []int{1111, 2222},
		OrgMember:   []int{1111, 4444, 2222},
		Banned:      false,
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}, 0, []string{"User successfully added in addUser"}})
	//Bad Non-Existent UserID
	userCrudUpdateResults = append(userCrudUpdateResults, UserCrudUpdate{User{
		UserName:    "TheWrongValue",
		Password:    hex.EncodeToString([]byte("testpword")),
		Firstname:   "Test",
		Lastname:    "User",
		PhoneNums:   []string{"13143228594"},
		UserID:      4000000,
		Email:       []string{"jbkeller0303@gmail.com"},
		Whoare:      "I am a test User and I like to write tests",
		AdminOrgs:   []int{1111, 2222},
		OrgMember:   []int{1111, 4444, 2222},
		Banned:      false,
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}, 1, []string{"Error adding User in addUser", "Error reading the request"}})
	//Bad Empty User Crud
	userCrudUpdateResults = append(userCrudUpdateResults, UserCrudUpdate{User{}, 1,
		[]string{"Error adding User in addUser", "Error reading the request"}})
}

//This creates our CRUD Delete Cases for deleting Users
func createUserDeleteCrud() {
	//Good User Crud Read
	userCrudDeleteResults = append(userCrudDeleteResults, UserCrudDelete{1111, 0, []string{"User successfully deleted in deleteUser"}})
	//Bad User CRUD Read
	userCrudDeleteResults = append(userCrudDeleteResults, UserCrudDelete{0, 1,
		[]string{"Error deleting User in deleteUser", "Error reading the request"}})
	//Not seen UserID
	userCrudDeleteResults = append(userCrudDeleteResults, UserCrudDelete{4000000, 1,
		[]string{"Error deleting User in deleteUser", "Error reading the request"}})
	//Another not seen UserID
	userCrudDeleteResults = append(userCrudDeleteResults, UserCrudDelete{-1, 1,
		[]string{"Error deleting User in deleteUser", "Error reading the request"}})
}

//This creates our login tests for logging Users in
func createUserLogin() {
	//Good User Login
	userLoginResults = append(userLoginResults, UserLoginTest{TheUsername: "TestUsername", ThePassword: hex.EncodeToString([]byte("testpword")),
		ExpectedNum: 0, ExpectedStringArray: []string{"User should be successfully decoded."}})
	//Bad User Login Username
	userLoginResults = append(userLoginResults, UserLoginTest{TheUsername: "BadUsername", ThePassword: hex.EncodeToString([]byte("testpword")),
		ExpectedNum: 1, ExpectedStringArray: []string{"No User was returned."}})
	//Bad User Password Username
	userLoginResults = append(userLoginResults, UserLoginTest{TheUsername: "BadUsername", ThePassword: hex.EncodeToString([]byte("badPWord")),
		ExpectedNum: 1, ExpectedStringArray: []string{"No User was returned."}})
	//Bad nil username
	userLoginResults = append(userLoginResults, UserLoginTest{TheUsername: "", ThePassword: hex.EncodeToString([]byte("badPWord")),
		ExpectedNum: 1, ExpectedStringArray: []string{"No User was returned."}})
	//Bad nil password
	userLoginResults = append(userLoginResults, UserLoginTest{TheUsername: "BadUsername", ThePassword: hex.EncodeToString([]byte("")),
		ExpectedNum: 1, ExpectedStringArray: []string{"No User was returned."}})
}

//Add a User
func TestUserAdd(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range userCrudCreateResults {
		/* start listener */
		/* 1. Create Context */
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		/* 2. Marshal test case to JSON expect */
		theJSONMessage, err := json.Marshal(test.TheUser)
		if err != nil {
			fmt.Println(err)
			logWriter(err.Error())
			log.Fatal(err)
		}
		/* 3. Create Post to JSON */
		payload := strings.NewReader(string(theJSONMessage))
		req, err := http.NewRequest("POST", ADDUSERURL, payload)
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

//Test for User Login with Username and Password
func TestUserLogin(t *testing.T) {
	/* start listener */
	type LoginData struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}
	loginData := LoginData{Username: "TestUsername", Password: hex.EncodeToString([]byte("testpword"))}
	/* 1. Create Context */
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(loginData)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
		t.Fatal("Could not Marshal JSON: " + err.Error())
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", GETUSERLOGIN, payload)
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
		log.Fatal(theMessage)
	}
}

//Test for updating Users
func TestUserUpdate(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range userCrudUpdateResults {
		/* start listener */
		/* 1. Create Context */
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		/* 2. Marshal test case to JSON expect */
		theJSONMessage, err := json.Marshal(test.TheUser)
		if err != nil {
			fmt.Println(err)
			logWriter(err.Error())
			log.Fatal(err)
			t.Fatal("Could not Marshal JSON: " + err.Error())
		}
		/* 3. Create Post to JSON */
		payload := strings.NewReader(string(theJSONMessage))
		req, err := http.NewRequest("POST", UPDATEURL, payload)
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

//Test for Reading Users
func TestUserRead(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range userCrudReadResults {
		/* 1. Create Context */
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		/* 2. Marshal test case to JSON expect */
		type UserIDUser struct {
			TheUserID int `json:"TheUserID"`
		}
		aUserID := UserIDUser{TheUserID: test.TheUserID}
		theJSONMessage, err := json.Marshal(aUserID)
		if err != nil {
			fmt.Println(err)
			logWriter(err.Error())
			log.Fatal(err)
			t.Fatal("Could not marshal JSON")
		}
		/* 3. Create Post to JSON */
		payload := strings.NewReader(string(theJSONMessage))
		req, err := http.NewRequest("POST", READUSERURL, payload)
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

//Test for getting all Usernames
func TestGetAllUsernames(t *testing.T) {
	mapOusernameToReturn := make(map[string]bool) //Username to load our values into
	//Call our crudOperations Microservice in order to get our Usernames
	//Create a context for timing out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", GETALLUSERNAMESURL, nil)
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
	} else {
		mapOusernameToReturn = returnedMessage.ReturnedUserMap
		fmt.Printf("Here is our map: %v\n", mapOusernameToReturn)
	}
}

//Test for Deleting Users
func TestUserDelete(t *testing.T) {
	time.Sleep(2 * time.Second) //Might needed for CRUD updating
	testNum := 0                //Used for incrementing
	for _, test := range userCrudDeleteResults {
		/* 1. Create Context */
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		/* 2. Marshal test case to JSON expect */
		type UserDelete struct {
			UserID int `json:"UserID"`
		}
		aUserID := UserDelete{UserID: test.TheUserID}
		theJSONMessage, err := json.Marshal(aUserID)
		if err != nil {
			fmt.Println(err)
			logWriter(err.Error())
			log.Fatal(err)
			t.Fatal("Could not marshal JSON")
		}
		/* 3. Create Post to JSON */
		payload := strings.NewReader(string(theJSONMessage))
		req, err := http.NewRequest("POST", DELETEURL, payload)
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
			fmt.Printf("We here test un-expected\n")
			t.Fatal("Wrong num recieved on testcase " + strconv.Itoa(testNum) +
				" :" + strconv.Itoa(returnedMessage.SuccOrFail) + " Expected: " + strconv.Itoa(test.ExpectedNum))
		}
		/* Maybe we can test the strings at some point... */
		testNum = testNum + 1 //Increment this number for testing
	}
}
