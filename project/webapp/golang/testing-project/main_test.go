package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

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
	createUserUpdateCrud() // Add our User Crud testing values for updating
	createUserDeleteCrud() //Add our User Crud testing values for deleting
	createUserLogin()      //Create creds for logging Users in
	/* Add values for LearnR Org test cases */
	createCreateLearnOrgCrud()
	createLearnOrgCrud()
	createLearnOrgUpdateCrud()
	createLearnOrgDeleteCrud()
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

//Get a random ID we can use for any struct
func TestRandomID(t *testing.T) {
	//Call our crudOperations Microservice in order to get our Usernames
	//Create a context for timing out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", GETRANDOMID, nil)
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
		t.Fatal("Had an error getting map: " + errString)
	} else {
		fmt.Printf("Here is our random ID: %v\n", returnedMessage.RandomID)
	}
}
