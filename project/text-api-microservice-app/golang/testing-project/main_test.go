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
	createLearnRTextSession()
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

//This pings our AWS service to check it it's up
/*
func TestServerPing(t *testing.T) {
	type SendJSON struct {
		TestNum int `json:"TestNum"`
	}
	sendJSON := SendJSON{TestNum: 0}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	theJSONMessage, err := json.Marshal(sendJSON)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
	}

	payload := strings.NewReader(string(theJSONMessage))
	req, err := http.NewRequest("POST", TESTPINGURL, payload)
	if err != nil {
		log.Fatal(err)
	}
	//req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Content-Type", "application/json")

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
		theErr := "There was an error reading response from testLocalPing " + err.Error()
		t.Fatal(theErr)
	}
	type SuccessMSG struct {
		Message    string `json:"Message"`
		SuccessNum int    `json:"SuccessNum"`
	}
	var returnedMessage SuccessMSG
	json.Unmarshal(body, &returnedMessage)

	if returnedMessage.SuccessNum != 0 {
		t.Fatal("Wrong reponse gotten from testLocalPing: " + returnedMessage.Message +
			" SuccessNum = " + strconv.Itoa(returnedMessage.SuccessNum))
	}
}
*/
