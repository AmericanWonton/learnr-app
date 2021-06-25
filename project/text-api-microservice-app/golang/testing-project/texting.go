package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/* Credentials for Twilio...THESE NEED TO BE READ IN AT SOME POINT */
var accountSID string = "ACceed35250b686bf863fe15fc117d7ac0" //THIS NEEDS TO GET READ IN AT SOME POINT
var authToken string = "44fb9877523c38df14407000419be0d3"    //THIS ALSO NEEDS TO GET READ IN
var urlStr string = "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json"
var coolquote string = "Hello it's joe from you on Twilio. You can check out my resume at http://josephkeller.me/. Reply with STOP, CANCEL, or QUIT to quit this."

type UserSession struct {
	LocalSessID    int           `json:"LocalSessID"`
	TheUser        User          `json:"TheUser"`
	TheLearnR      Learnr        `json:"TheLearnR"`
	TheLearnRInfo  LearnrInfo    `json:"TheLearnRInfo"`
	PersonName     string        `json:"PersonName"`
	PersonPhoneNum string        `json:"PersonPhoneNum"`
	TheSession     LearnRSession `json:"TheSession"`
	LogInfo        []string      `json:"LogInfo"`
}

//Channel for Go-Routines
var learnSessChannel chan UserSession
var learnSessResultChannel chan UserSession

/* A map of our active sessions */
var UserSessionActiveMap map[int]UserSession

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
		TheUser        User       `json:"TheUser"`
		TheLearnR      Learnr     `json:"TheLearnR"`
		TheLearnRInfo  LearnrInfo `json:"TheLearnRInfo"`
		PersonName     string     `json:"PersonName"`
		PersonPhoneNum string     `json:"PersonPhoneNum"`
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
			theInt, theErr := fmt.Fprint(w, string(theJSONMessage))
			if theErr != nil {
				logWriter("Error writing back to User in UserSession Addition: " + theErr.Error() + " " + strconv.Itoa(theInt))
			}
			/* Session Added. Begin Go routine to start texting them.
			Create User Session to add onto Channel */
			newUserSession := UserSession{
				LocalSessID:    getRandomID(),
				TheUser:        theJSON.TheUser,
				TheLearnR:      theJSON.TheLearnR,
				TheLearnRInfo:  theJSON.TheLearnRInfo,
				PersonName:     theJSON.PersonName,
				PersonPhoneNum: theJSON.PersonPhoneNum,
				TheSession:     newLearnRSession,
				LogInfo:        []string{},
			}
			learnSessChannel <- newUserSession
			go learnRSession(getRandomID(), learnSessChannel, learnSessResultChannel)
		}
	}
}

func learnRSession(workerID int, userSessionChan <-chan UserSession, userSessCloseChan chan<- UserSession) {
	for a := range userSessionChan {
		//Get this UserSession off the channel and into a good variable
		theUserSession := a
		//Start sending texts on this session
		for l := 0; l < len(theUserSession.TheLearnR.LearnRInforms); l++ {
			//Do this only if the 'Ongoing' variable in the LearnRSession is true
			if theUserSession.TheSession.Ongoing {
				//Add this User Session to our map
				UserSessionActiveMap[theUserSession.LocalSessID] = theUserSession
				theTimeNow := time.Now()
				goodSend, resultMessages := sendText(l, theUserSession.PersonPhoneNum, theUserSession.TheLearnR.PhoneNums[0],
					theUserSession.TheLearnR.LearnRInforms[l].TheInfo)
				if !goodSend {
					//Failed to send text; log this in session and elsewhere
					theUserSession.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
					//Collect message
					message := ""
					for j := 0; j < len(resultMessages); j++ {
						message = message + resultMessages[j] + " "
					}
					logWriter(message)
					theUserSession.LogInfo = append(theUserSession.LogInfo, message)
					//Update our UserSession Map
					UserSessionActiveMap[theUserSession.LocalSessID] = theUserSession
				} else {
					//Text successfully sent; log this and put in session info
					theUserSession.TheSession.TextsSent = append(theUserSession.TheSession.TextsSent)
					theUserSession.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
					//Collect message
					message := ""
					for j := 0; j < len(resultMessages); j++ {
						message = message + resultMessages[j] + " "
					}
					logWriter(message)
					theUserSession.LogInfo = append(theUserSession.LogInfo, message)
					//Update our UserSession Map
					UserSessionActiveMap[theUserSession.LocalSessID] = theUserSession
				}
				//Account for time delays of the next LearnRInforms
				if theUserSession.TheLearnR.LearnRInforms[l].ShouldWait {
					time.Sleep(time.Second * time.Duration(theUserSession.TheLearnR.LearnRInforms[l].WaitTime))
				}
			}
		}
		/* Done sending all texts for this LearnR. We can put this on the UserSessClose Chan and
		begin sending the results to User/CRUD DB */
		//Update Session
		theTimeNow := time.Now()
		theUserSession.TheSession.Ongoing = false //Stop the session from continuing
		theUserSession.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		//Add session information to session DB
		fmt.Printf("DEBUG: Adding learnrSession to learnRSession DB\n")
		//Update the LearnRInfo for this LearnR with our updated Session added to it
		fmt.Printf("DEBUG: Adding our leanrrSession to our LearnRInform and updating it")
		//Removing Map placement of this UserSession
		delete(UserSessionActiveMap, theUserSession.LocalSessID)
		//Add this to our Results Chan
		userSessCloseChan <- a
		fmt.Printf("We are now done texting this User\n")
	}
}

func getRandomID() int {
	finalID := 0 //The final, unique ID to return to the food/user
	randInt := 0 //The random integer added onto ID

	randIntString := "" //The integer built through a string...
	min, max := 0, 9    //The min and Max value for our randInt

	for i := 0; i < 12; i++ {
		randInt = rand.Intn(max-min) + min
		randIntString = randIntString + strconv.Itoa(randInt)
	}

	finalID, err := strconv.Atoi(randIntString)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}

	return finalID
}

func sendText(textOrder int, toNumString string, fromNumString string, textBody string) (bool, []string) {
	goodSend, resultMessages := true, []string{}
	msgData := url.Values{}
	msgData.Set("To", "+"+toNumString)
	msgData.Set("From", "+"+fromNumString)
	msgData.Set("Body", textBody)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSID, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//Get that request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Define the response we expect to recieve
	type TwilioGetBack struct {
		Code     int    `json:"code"`
		Message  string `json:"message"`
		MoreInfo string `json:"more_info"`
		Status   int    `json:"status"`
	}

	resp, theErr := client.Do(req.WithContext(ctx))
	if (resp.StatusCode >= 200 && resp.StatusCode < 300) && theErr == nil {
		theMessage := ""
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			theErr := "There was an error reading a response for Twilio: " + err.Error()
			fmt.Println(theErr)
			goodSend = false
			resultMessages = append(resultMessages, theErr)
		}
		var returnedMessage TwilioGetBack
		json.Unmarshal(body, &returnedMessage)
		theMessage = "Successful message sent"
		resultMessages = append(resultMessages, theMessage)
	} else {
		theErr := "Bad response code recieved after sending text: " + strconv.Itoa(resp.StatusCode) + "\nError: " + theErr.Error()
		fmt.Println(theErr)
		goodSend = false
		resultMessages = append(resultMessages, theErr)
	}

	//Close this response, just in case
	resp.Body.Close()

	return goodSend, resultMessages
}
