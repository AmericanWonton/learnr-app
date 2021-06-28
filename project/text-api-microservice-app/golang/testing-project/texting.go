package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Here is our waitgroup
var wg sync.WaitGroup

/* DEBUG ping values */
var INITIALLEARNRSEND string = "http://localhost:3000/initialLearnRStart"

/* Credentials for Twilio...THESE NEED TO BE READ IN AT SOME POINT */
var accountSID string
var authToken string
var urlStr string

//Used to record active sessions with our LearnRs
type UserSession struct {
	Active             bool          `json:"Active"`
	LocalSessID        int           `json:"LocalSessID"`
	IntroductionSaying string        `json:"IntroductionSaying"`
	TheUserName        string        `json:"TheUserName"`
	TheUser            User          `json:"TheUser"`
	TheLearnR          Learnr        `json:"TheLearnR"`
	TheLearnRInfo      LearnrInfo    `json:"TheLearnRInfo"`
	PersonName         string        `json:"PersonName"`
	PersonPhoneNum     string        `json:"PersonPhoneNum"`
	TheSession         LearnRSession `json:"TheSession"`
	LogInfo            []string      `json:"LogInfo"`
}

//Channel for Go-Routines
var learnSessChannel chan UserSession
var learnSessResultChannel chan UserSession

/* A map of our active sessions */
var UserSessionActiveMap map[int]UserSession

/* A Map of phone numbers that LINK us to those active sessions, (on the random session id) */
var UserSessPhoneMap map[string]int

/* Get our creds for twilio */
func getTwilioCreds() {
	file, err := os.Open("security/twiliocreds.txt")

	if err != nil {
		fmt.Printf("DEBUG: Trouble opening bad word text file: %v\n", err.Error())
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	accountSID = text[0]
	authToken = text[1]
	urlStr = "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json"
}

/* Called from our webpage to initiate a learnr request to another person */
func initialLearnRStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	/* Test Flusher stuff */
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Server does not support Flusher!",
			http.StatusInternalServerError)
		return
	}

	//Declare Ajax return statements to be sent back
	type SuccessMSG struct {
		Message    string `json:"Message"`
		SuccessNum int    `json:"SuccessNum"`
	}
	theSuccMessage := SuccessMSG{
		Message:    "LearnRBegun successfully",
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
		Introduction   string     `json:"Introduction"`
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
		/* Session Added. Begin Go routine to start texting them.
		Create User Session to add onto Channel */
		newUserSession := UserSession{
			Active:             true,
			LocalSessID:        getRandomID(),
			IntroductionSaying: theJSON.Introduction,
			TheUserName:        theJSON.TheUser.UserName,
			TheUser:            theJSON.TheUser,
			TheLearnR:          theJSON.TheLearnR,
			TheLearnRInfo:      theJSON.TheLearnRInfo,
			PersonName:         theJSON.PersonName,
			PersonPhoneNum:     theJSON.PersonPhoneNum,
			TheSession:         newLearnRSession,
			LogInfo:            []string{},
		}
		//go learnRSession(learnSessChannel, learnSessResultChannel)
		/* Send the response back to Ajax */
		theJSONMessage, err := json.Marshal(theSuccMessage)
		//Send the response back
		if err != nil {
			errIs := "Error formatting JSON for return in initialLearnRStart: " + err.Error()
			logWriter(errIs)
			panic(errIs)
		}
		theInt, theErr := fmt.Fprint(w, string(theJSONMessage))
		if theErr != nil {
			logWriter("Error writing back to initialLearnRStart: " + theErr.Error() + " " + strconv.Itoa(theInt))
			panic("Error writing back to initialLearnRStart: " + theErr.Error() + " " + strconv.Itoa(theInt))
		}
		flusher.Flush()
		go conductLearnRSession(newUserSession)
	}
}

func conductLearnRSession(theLearnRUserSess UserSession) {
	//Add this User Session to our map of phone numbers
	UserSessPhoneMap[theLearnRUserSess.PersonPhoneNum] = theLearnRUserSess.LocalSessID
	/* First send the introduction Texts to this User; we will give them a
	160 second break period in order to type STOP; if they do not, our LearnR continues, until
	they enter STOP at any time */
	continueLearnR := true                                                  //This will determine if we can send the rest of our texts
	UserSessionActiveMap[theLearnRUserSess.LocalSessID] = theLearnRUserSess //Add this User Session to our map
	//Send the Introduction Text
	introMessage := "Hello " + theLearnRUserSess.PersonName + ", " + theLearnRUserSess.TheUserName + " wanted to help educate you on " +
		"something important to them."
	theTimeNow := time.Now()
	goodSend, resultMessages := sendText(-3, theLearnRUserSess.PersonPhoneNum, theLearnRUserSess.TheLearnR.PhoneNums[0],
		introMessage)
	if !goodSend || !UserSessionActiveMap[theLearnRUserSess.LocalSessID].Active {
		//Intro text failed...LearnR may not be active
		continueLearnR = false
		theLearnRUserSess.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		//Collect message
		message := ""
		for j := 0; j < len(resultMessages); j++ {
			message = message + resultMessages[j] + " "
		}
		logWriter(message)
		theLearnRUserSess.LogInfo = append(theLearnRUserSess.LogInfo, message)
		//Update our UserSession Map
		UserSessionActiveMap[theLearnRUserSess.LocalSessID] = theLearnRUserSess
	} else {
		//Send the second text with our Users message
		time.Sleep(time.Second * 10) //Small wait
		introMessage = "\"" + theLearnRUserSess.IntroductionSaying + "\""
		goodSend, resultMessages := sendText(-2, theLearnRUserSess.PersonPhoneNum, theLearnRUserSess.TheLearnR.PhoneNums[0],
			introMessage)
		if !goodSend || !UserSessionActiveMap[theLearnRUserSess.LocalSessID].Active {
			//Intro text failed...LearnR may not be active
			continueLearnR = false
			theLearnRUserSess.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
			//Collect message
			message := ""
			for j := 0; j < len(resultMessages); j++ {
				message = message + resultMessages[j] + " "
			}
			logWriter(message)
			theLearnRUserSess.LogInfo = append(theLearnRUserSess.LogInfo, message)
			//Update our UserSession Map
			UserSessionActiveMap[theLearnRUserSess.LocalSessID] = theLearnRUserSess
		} else {
			//Send final message asking User to confirm/deny
			//Send the second text with our Users message
			time.Sleep(time.Second * 10) //Small wait
			introMessage = "At their request, we'll begin sending a LearnR to help them explain.\n\n" +
				"If you'd like to stop, please text STOP back to this number at any time."
			goodSend, resultMessages := sendText(-1, theLearnRUserSess.PersonPhoneNum, theLearnRUserSess.TheLearnR.PhoneNums[0],
				introMessage)
			if !goodSend || !UserSessionActiveMap[theLearnRUserSess.LocalSessID].Active {
				//Intro text failed...LearnR may not be active
				continueLearnR = false
				theLearnRUserSess.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
				//Collect message
				message := ""
				for j := 0; j < len(resultMessages); j++ {
					message = message + resultMessages[j] + " "
				}
				logWriter(message)
				theLearnRUserSess.LogInfo = append(theLearnRUserSess.LogInfo, message)
				//Update our UserSession Map
				UserSessionActiveMap[theLearnRUserSess.LocalSessID] = theLearnRUserSess
			}
		}
	}
	time.Sleep(time.Second * 10) //Small wait
	/* Hopefully we've sent our first three texts successfully, (continueLearnR == true).
	If not, log the failure remove this Session/Phone Map recording and update our DB*/
	if !continueLearnR || !UserSessionActiveMap[theLearnRUserSess.LocalSessID].Active {
		//Failed sending messages...ending session
		theMessage := "LearnR Session ending; did not have success sending texts or User terminated session..."
		fmt.Println(theMessage)
		logWriter(theMessage)
		theLearnRUserSess.Active = false
		UserSessionActiveMap[theLearnRUserSess.LocalSessID] = theLearnRUserSess
		theTimeNow = time.Now()
		theLearnRUserSess.TheSession.Ongoing = false //Stop the session from continuing
		theLearnRUserSess.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		//Add session information to session DB
		wg.Add(1)
		go fastAddLearnRSession(theLearnRUserSess.TheSession)
		//Update the LearnRInfo for this LearnR with our updated Session added to it
		theLearnRUserSess.TheLearnRInfo.AllSessions = append(theLearnRUserSess.TheLearnRInfo.AllSessions, theLearnRUserSess.TheSession)
		theLearnRUserSess.TheLearnRInfo.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		wg.Add(1)
		go fastUpdateLearnRInform(theLearnRUserSess.TheLearnRInfo)
		//Removing Map placement of this UserSession
		delete(UserSessionActiveMap, theLearnRUserSess.LocalSessID)
		//Removing Phone Num from active UserSession
		delete(UserSessPhoneMap, theLearnRUserSess.PersonPhoneNum)
		wg.Wait()
	} else {
		//First three messages sent, getting ready to send the rest of the messages...
		//Start sending texts on this session
		for l := 0; l < len(theLearnRUserSess.TheLearnR.LearnRInforms); l++ {
			//Do this only if the 'Ongoing' variable in the LearnRSession is true
			if theLearnRUserSess.TheSession.Ongoing && UserSessionActiveMap[theLearnRUserSess.LocalSessID].Active {
				theTimeNow := time.Now()
				goodSend, resultMessages := sendText(l, theLearnRUserSess.PersonPhoneNum, theLearnRUserSess.TheLearnR.PhoneNums[0],
					theLearnRUserSess.TheLearnR.LearnRInforms[l].TheInfo)
				if !goodSend {
					//Failed to send text; log this in session and elsewhere
					theLearnRUserSess.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
					//Collect message
					message := ""
					for j := 0; j < len(resultMessages); j++ {
						message = message + resultMessages[j] + " "
					}
					logWriter(message)
					theLearnRUserSess.LogInfo = append(theLearnRUserSess.LogInfo, message)
					//Update our UserSession Map
					UserSessionActiveMap[theLearnRUserSess.LocalSessID] = theLearnRUserSess
				} else {
					//Text successfully sent; log this and put in session info
					theLearnRUserSess.TheSession.TextsSent = append(theLearnRUserSess.TheSession.TextsSent, theLearnRUserSess.TheLearnR.LearnRInforms[l])
					theLearnRUserSess.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
					//Collect message
					message := ""
					for j := 0; j < len(resultMessages); j++ {
						message = message + resultMessages[j] + " "
					}
					logWriter(message)
					theLearnRUserSess.LogInfo = append(theLearnRUserSess.LogInfo, message)
					//Update our UserSession Map
					UserSessionActiveMap[theLearnRUserSess.LocalSessID] = theLearnRUserSess
				}
				//Account for time delays of the next LearnRInforms
				if theLearnRUserSess.TheLearnR.LearnRInforms[l].ShouldWait {
					time.Sleep(time.Second * time.Duration(theLearnRUserSess.TheLearnR.LearnRInforms[l].WaitTime))
				}
			}
		}
		/* Done sending all texts for this LearnR. We can put this on the UserSessClose Chan and
		begin sending the results to User/CRUD DB */
		//Update Session
		theLearnRUserSess.Active = false
		UserSessionActiveMap[theLearnRUserSess.LocalSessID] = theLearnRUserSess
		theTimeNow = time.Now()
		theLearnRUserSess.TheSession.Ongoing = false //Stop the session from continuing
		theLearnRUserSess.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		//Add session information to session DB
		wg.Add(1)
		go fastAddLearnRSession(theLearnRUserSess.TheSession)
		//Update the LearnRInfo for this LearnR with our updated Session added to info
		theLearnRUserSess.TheLearnRInfo.AllSessions = append(theLearnRUserSess.TheLearnRInfo.AllSessions, theLearnRUserSess.TheSession)
		theLearnRUserSess.TheLearnRInfo.FinishedSessions = append(theLearnRUserSess.TheLearnRInfo.FinishedSessions, theLearnRUserSess.TheSession)
		theLearnRUserSess.TheLearnRInfo.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		wg.Add(1)
		go fastUpdateLearnRInform(theLearnRUserSess.TheLearnRInfo)
		//Removing Map placement of this UserSession
		delete(UserSessionActiveMap, theLearnRUserSess.LocalSessID)
		//Removing Phone Num from active UserSession
		delete(UserSessPhoneMap, theLearnRUserSess.PersonPhoneNum)
		wg.Wait() //Need to make sure we can exit this function properly
	}
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
	/*
		type TwilioGetBack struct {
			Code     int    `json:"code"`
			Message  string `json:"message"`
			MoreInfo string `json:"more_info"`
			Status   int    `json:"status"`
		}
	*/
	type TwilioResponse struct {
		Sid                 string `json:"sid"`
		DateCreated         string `json:"date created"`
		DateUpdated         string `json:"date updated"`
		DateSent            string `json:"date sent"`
		AccountSid          string `json:"account sid"`
		To                  string `json:"to"`
		From                string `json:"from"`
		MessagingServiceSid string `json:"messaging service sid"`
		Body                string `json:"body"`
		Status              string `json:"status"`
		NumSegments         string `json:"num segments"`
		NumMedia            string `json:"num media"`
		Direction           string `json:"direction"`
		APIVersion          string `json:"api version"`
		Price               string `json:"price"`
		PriceUnit           string `json:"price unit"`
		ErrorCode           string `json:"error code"`
		ErrorMessage        string `json:"error message"`
		URI                 string `json:"uri"`
		SubresourceUris     struct {
			Media string `json:"media"`
		} `json:"subresource uris"`
	}
	resp, theErr := client.Do(req.WithContext(ctx))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			theErr := "There was an error reading a response for Twilio: " + err.Error()
			fmt.Println(theErr)
		}
		var returnedMessage TwilioResponse
		json.Unmarshal(body, &returnedMessage)
		//Check for correct response obtained
		if strings.Contains(strings.ToLower(returnedMessage.Body), strings.ToLower("Sent from")) {
			//Successful text
			message := "Good text response obtained"
			resultMessages = append(resultMessages, message)
		} else {
			//Not successful response
			goodSend = false
			theErr := "Could not obtain the correct body response: " + returnedMessage.Body
			resultMessages = append(resultMessages, theErr)
		}
	} else {
		fmt.Println(resp.Status)
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("There's an error reading the response: %v\n", err.Error())
		} else {
			fmt.Printf("Here is the response: %v\n", string(b))
		}
		goodSend = false
		resultMessages = append(resultMessages, theErr.Error())
	}

	//Close this response, just in case
	resp.Body.Close()

	return goodSend, resultMessages
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
