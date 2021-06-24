package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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
	getCredsMongo() //Get mongo creds
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
	_, err := ioutil.ReadFile("logging/textapilog.txt")
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
