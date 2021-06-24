package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/* Credentials for Twilio...THESE NEED TO BE READ IN AT SOME POINT */
var accountSID string = "ACceed35250b686bf863fe15fc117d7ac0" //THIS NEEDS TO GET READ IN AT SOME POINT
var authToken string = "44fb9877523c38df14407000419be0d3"    //THIS ALSO NEEDS TO GET READ IN
var urlStr string = "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json"
var coolquote string = "Hello it's joe from you on Twilio. You can check out my resume at http://josephkeller.me/. Reply with STOP, CANCEL, or QUIT to quit this."

type UserSession struct {
	TheUser        User          `json:"TheUser"`
	TheLearnR      Learnr        `json:"TheLearnR"`
	PersonName     string        `json:"PersonName"`
	PersonPhoneNum string        `json:"PersonPhoneNum"`
	TheSession     LearnRSession `json:"TheSession"`
}

//Channel for Go-Routines
var learnSessChannel chan UserSession

/* Called from our webpage to initiate a learnr request to another person */
func initialLearnRStart(w http.ResponseWriter, r *http.Request) {
	//Declare Ajax return statements to be sent back
	type SuccessMSG struct {
		Message    string `json:"Message"`
		SuccessNum int    `json:"SuccessNum"`
	}
	theSuccMessage := SuccessMSG{
		Message:    "User created successfully",
		SuccessNum: 0,
	}

	//Get the byte slice from the request
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	type OurJSON struct {
		TheUser        User   `json:"TheUser"`
		TheLearnR      Learnr `json:"TheLearnR"`
		PersonName     string `json:"PersonName"`
		PersonPhoneNum string `json:"PersonPhoneNum"`
	}
	//Marshal it into our type
	var theJSON OurJSON
	json.Unmarshal(bs, &theJSON)

	/* Create Session for this LearnR to this person*/
	//Get Random ID
	theTimeNow := time.Now()
	goodGet, message, randomID := randomAPICall()
	if !goodGet {
		theErr := "Failure to get random API in session: " + message
		theSuccMessage.Message = theErr
		logWriter(theErr)
	} else {
		newLearnRSession := LearnRSession{
			ID:               randomID,
			LearnRID:         theJSON.TheLearnR.ID,
			LearnRName:       theJSON.TheLearnR.Name,
			TheLearnR:        theJSON.TheLearnR,
			TheUser:          theJSON.TheUser,
			TargetUserNumber: theJSON.PersonPhoneNum,
			Ongoing:          true,
			TextsSent:        []LearnRInforms{},
			UserResponses:    []string{},
			DateCreated:      theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated:      theTimeNow.Format("2006-01-02 15:04:05"),
		}
		goodAdd, message := callAddLearnRSession(newLearnRSession)
		if !goodAdd {
			theErr := "Failure to get random API in session: " + message
			theSuccMessage.Message = theErr
			logWriter(theErr)
		} else {
			/* Send the response back to Ajax */
			theJSONMessage, err := json.Marshal(theSuccMessage)
			//Send the response back
			if err != nil {
				errIs := "Error formatting JSON for return in createUser: " + err.Error()
				logWriter(errIs)
			}
			fmt.Fprint(w, string(theJSONMessage))
			/* Session Added. Begin Go routine to start texting them */

		}
	}
}

func learnRSession(userSessChannel chan UserSession) {
	sessDone := false //Does not end session until done
}
