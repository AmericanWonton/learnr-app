package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculate(t *testing.T) {
	if Calculate(2) != 4 {
		t.Error("Expected 2 + 2 to equal 4")
	}
}

func TestTableCalculate(t *testing.T) {
	var tests = []struct {
		input    int
		expected int
	}{
		{2, 4},
		{-1, 1},
		{0, 2},
		{999, 1001},
	}

	for _, test := range tests {
		if output := Calculate(test.input); output != test.expected {
			t.Error("Test failed: {} inputted, {} expected, recieved: {}",
				test.input, test.expected, output)
		}
	}
}

/* ADVANCED TESTING TECHNIQUE BEGINNING */

/* TEST TABLE EXAMPLE */
type AddResult struct {
	x        int
	y        int
	expected int
}

/* This is an array of test cases we are going to put into TestAdd */
var addResults = []AddResult{
	{1, 1, 2},
}

func TestAdd(t *testing.T) {
	/* For the range of all of our test cases in 'addResults', run them
	through the 'Add' function, get our result, then see if it matches
	the expected result within this element of AddResult */
	for _, test := range addResults {
		result := Add(test.x, test.y)
		if result != test.expected {
			t.Fatal("Expected Result not given")
		}
	}
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
	/* Template definition section
	/* TEMPLATE DEFINITION */
	template1 = template.Must(template.ParseGlob("./static/templates/*"))
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
/*
func TestHandle(t *testing.T) {
	//Here is our waitgroup
	theTimer := 5
	doneCounting := false
	var wg sync.WaitGroup
	go func() {
		for l := 0; l != theTimer; l++ {
			if !doneCounting {
				time.Sleep(1 * time.Second)
				fmt.Printf("The timer is...%v\n", l)
			} else {
				doneCounting = true
				break
			}
		}
		wg.Done()
	}()

	//Declare a func that will deter closing the handleRequest
	for !doneCounting {
		handleRequests()
	}

	wg.Wait()
}
*/

/* Test Index */
