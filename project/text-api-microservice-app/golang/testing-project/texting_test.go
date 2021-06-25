package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

/* Declarative Test structs/values */

type OurJSON struct {
	TheUser        User       `json:"TheUser"`
	TheLearnR      Learnr     `json:"TheLearnR"`
	TheLearnRInfo  LearnrInfo `json:"TheLearnRInfo"`
	PersonName     string     `json:"PersonName"`
	PersonPhoneNum string     `json:"PersonPhoneNum"`
	Introduction   string     `json:"Introduction"`
}
type LearnRTestSends struct {
	JSONSend            OurJSON
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var learnrTestSendResults []LearnRTestSends

func createLearnRTextSession() {
	//Get info for sending
	theUser, _, _ := callGetUser(228778447811)
	_, _, theLearnR := callReadLearnR(102471876033)
	_, _, theLearnRInfo := callReadLearnrInfo(718658150182)
	//Create first LearnR, success
	learnrTestSendResults = append(learnrTestSendResults, LearnRTestSends{
		JSONSend: OurJSON{
			TheUser:        theUser,
			TheLearnR:      theLearnR,
			TheLearnRInfo:  theLearnRInfo,
			PersonName:     "Greg Gregory Test",
			PersonPhoneNum: "13143228594",
			Introduction:   "Hey idiot, I want to educate you about gay frogs...in test",
		},
		ExpectedNum:         0,
		ExpectedTruth:       true,
		ExpectedStringArray: []string{"Good send"},
	})
}

func TestSendLearnR(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range learnrTestSendResults {
		/* start listener */
		/* 1. Create Context */
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		/* 2. Marshal test case to JSON expect */
		fmt.Printf("DEBUG: Made it  to marshalling\n")
		theJSONMessage, err := json.Marshal(test.JSONSend)
		if err != nil {
			t.Fatal("Error marshalling JSON: " + err.Error())
			fmt.Println(err)
			logWriter(err.Error())
		}
		/* 3. Create Post to JSON */
		payload := strings.NewReader(string(theJSONMessage))
		req, err := http.NewRequest("POST", INITIALLEARNRSEND, payload)
		if err != nil {
			t.Fatal(err.Error())
		}
		fmt.Printf("DEBUG: Made it to req headerer\n")
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
			theErr := "There was an error reading response from initialLearnRSend " + err.Error()
			t.Fatal(theErr)
		}
		type SuccessMSG struct {
			Message    string `json:"Message"`
			SuccessNum int    `json:"SuccessNum"`
		}
		var returnedMessage SuccessMSG
		json.Unmarshal(body, &returnedMessage)
		fmt.Printf("DEBUG: made it to reutnredmessage\n")
		/* 5. Evaluate response in returnedMessage for testing */
		if test.ExpectedNum != returnedMessage.SuccessNum {
			t.Fatal("Wrong num recieved on testcase " + strconv.Itoa(testNum) +
				" :" + strconv.Itoa(returnedMessage.SuccessNum) + " Expected: " + strconv.Itoa(test.ExpectedNum))
		}
		/* Maybe we can test the strings at some point... */
		testNum = testNum + 1 //Increment this number for testing
	}
}
