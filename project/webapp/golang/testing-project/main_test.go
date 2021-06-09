package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

/* DEFINE CRUD POSTS URLS FOR LOCALHOST */
const ADDUSERURL string = "http://localhost:4000/addUser"
const READUSERURL string = "http://localhost:4000/getUser"

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

func TestMain(m *testing.M) {
	//Build stuff for beginning of tests
	log.Println("Starting stuff in TestMain")
	fmt.Println("Starting stuff in TestMain")
	setup()
	code := m.Run()
	//Do stuff for ending of tests
	log.Println("Ending stuff in Test main")
	fmt.Println("Ending stuff in test main")
	shutdown()

	os.Exit(code)
}

//This is setup values declared for testing
func setup() {
	fmt.Printf("Setting up test values...\n")
	/* Start by connecting to Mongo client */
	getCredsMongo()        //Get mongo creds
	createCreateUserCrud() //Add our User Crud testing values for Create
	createUserReadCrud()   //Add our User Crud testing values for Reading
}

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
	userCrudCreateResults = append(userCrudCreateResults, UserCrudCreate{User{}, 0,
		[]string{"Error adding User in addUser", "Error reading the request"}})
}

//This creates our CRUD Testing cases for Reading Users
func createUserReadCrud() {
	//Good User Crud Read
	userCrudReadResults = append(userCrudReadResults, UserCrudRead{1111, 0, []string{"User successfully added in addUser"}})
	//Bad User CRUD Read
	userCrudReadResults = append(userCrudReadResults, UserCrudRead{0, 0,
		[]string{"Error adding User in addUser", "Error reading the request"}})
}

//This creates our Crud Testing cases for Reading Users
func createReadUserCrud() {
	//Good User Read Crud Create
}

//This is shutdown values/actions for testing
func shutdown() {
	fmt.Printf("Setting up shutdown values/functions...\n")
}

/* Test DIRECTORY EXAMPLE */
func TestReadFile(t *testing.T) {
	data, err := ioutil.ReadFile("test-data/test.data")
	if err != nil {
		t.Fatal("Could not open file:\n" + err.Error())
	}
	if string(data) != "hello world from test.data" {
		t.Fatal("String contents do not match expected")
	}
}

/* Test logwrite example */
func TestLogWriter(t *testing.T) {
	/* Test read */
	_, err := ioutil.ReadFile("logging/weblog.txt")
	if err != nil {
		t.Fatal("Could not open file:\n" + err.Error())
	}
	/* Test logwriter write */
	logWriter("This is a test message")
}

/* Test init example */
func Testinit(t *testing.T) {

}

/* Test HTTP Example */
func TestHTTPRequest(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "{ \"status\": \"good\" }")
	}

	r := httptest.NewRequest("GET", "http://josephkeller.me/", nil)
	w := httptest.NewRecorder()
	handler(w, r)

	resp := w.Result()
	body, theErr := ioutil.ReadAll(resp.Body)
	fmt.Printf("Here is our response code: %v\n", string(body))
	if 200 != resp.StatusCode {
		t.Fatal("Status Code not okay: " + theErr.Error())
	}
}

/* Test handle routes */
//TestUserAdd
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

//Test for Reading Users
func TestUserRead(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range userCrudReadResults {
		/* 1. Create Context */
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		/* 2. Marshal test case to JSON expect */
		theJSONMessage, err := json.Marshal(test.TheUserID)
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
		}
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
