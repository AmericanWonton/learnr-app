package main

import (
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

//Response Text from what our Users sent
type ResponseText struct {
	NumMedia            string `json:"NumMedia"`
	FromCountry         string `json:"FromCountry"`
	SmsStatus           string `json:"SmsStatus"`
	ApiVersion          string `json:"ApiVersion"`
	Body                string `json:"Body"`
	To                  string `json:"To"`
	ToCity              string `json:"ToCity"`
	FromCity            string `json:"FromCity"`
	AccountSid          string `json:"AccountSid"`
	ToState             string `json:"ToState"`
	From                string `json:"From"`
	MessagingServiceSid string `json:"MessagingServiceSid"`
	SmsSid              string `json:"SmsSid"`
	MessageSid          string `json:"MessageSid"`
	FromState           string `json:"FromState"`
	ToZip               string `json:"ToZip"`
	ToCountry           string `json:"ToCountry"`
	FromZip             string `json:"FromZip"`
	NumSegments         string `json:"NumSegments"`
	OptOutType          string `json:"OptOutType"`
}

//Here is our waitgroup
var wg sync.WaitGroup

const ALLOTTEDLEARNRTIME = 900

/* DEBUG ping values */

/* Credentials for Twilio...THESE NEED TO BE READ IN AT SOME POINT */
var accountSID string
var authToken string
var urlStr string

//Used to record active sessions with our LearnRs
type UserSession struct {
	Active             bool          `json:"Active"`
	StartTime          time.Time     `json:"StartTime"`
	EndTime            time.Time     `json:"EndTime"`
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

/* A map of our active sessions */
var UserSessionActiveMap map[int]UserSession

/* A Map of phone numbers that LINK us to those active sessions, (on the random session id) */
var UserSessPhoneMap map[string]int

/* A map of all the keywords for Users telling us to stop texting them */
var StopText map[string]string

/* Get our creds for twilio */
func getTwilioCreds() {
	//Check to see if ENV Creds are available first
	_, ok := os.LookupEnv("TWILIO_ACCNTID")
	if !ok {
		message := "This ENV Variable is not present: " + "TWILIO_ACCNTID"
		panic(message)
	}
	_, ok2 := os.LookupEnv("TWILIO_AUTHTOKEN")
	if !ok2 {
		message := "This ENV Variable is not present: " + "TWILIO_AUTHTOKEN"
		panic(message)
	}

	accountSID = os.Getenv("TWILIO_ACCNTID")
	authToken = os.Getenv("TWILIO_AUTHTOKEN")
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
			StartTime:          time.Now(),
			EndTime:            time.Time{},
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
		time.Sleep(time.Second * 1) //DEBUG
		go func() {
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
		}()
		time.Sleep(time.Second * 2) //Debug
		go conductLearnRSession(newUserSession)
	}
}

func conductLearnRSession(theLearnRUserSess UserSession) {
	/* Start the timer for this LearnRSession */
	theLearnRUserSess.StartTime = time.Now()
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
			fmt.Println(message)
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
	time.Sleep(time.Second * 120) //Small wait
	/* Hopefully we've sent our first three texts successfully, (continueLearnR == true).
	If not, log the failure remove this Session/Phone Map recording and update our DB*/
	if !continueLearnR || !UserSessionActiveMap[theLearnRUserSess.LocalSessID].Active {
		//Failed sending messages...ending session
		theMessage := "LearnR Session ending; did not have success sending texts or User terminated session..."
		fmt.Println(theMessage)
		logWriter(theMessage)
		theLearnRUserSess.Active = false
		theLearnRUserSess.EndTime = time.Now()
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
		wg.Wait()
		//Removing Map placement of this UserSession
		delete(UserSessionActiveMap, theLearnRUserSess.LocalSessID)
		//Removing Phone Num from active UserSession
		delete(UserSessPhoneMap, theLearnRUserSess.PersonPhoneNum)
		fmt.Printf("This LearnR Session has been ended: %v for this LearnR: %v\n", theLearnRUserSess.LocalSessID, theLearnRUserSess.TheLearnR.Name)
	} else {
		//First three messages sent, getting ready to send the rest of the messages...
		//Start sending texts on this session
		for l := 0; l < len(theLearnRUserSess.TheLearnR.LearnRInforms); l++ {
			//Do this only if the 'Ongoing' variable in the LearnRSession is true
			if theLearnRUserSess.TheSession.Ongoing && UserSessionActiveMap[theLearnRUserSess.LocalSessID].Active {
				theTimeNow := time.Now()
				fmt.Printf("Sending this LearnR Text now...%v\n", theLearnRUserSess.TheLearnR.LearnRInforms[l])
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
				//Check to see if this LearnR has been going on for too long...
				timeDuration := time.Since(theLearnRUserSess.StartTime).Seconds()
				if timeDuration >= ALLOTTEDLEARNRTIME {
					//LearnR has been going on for too long...killing this session
					theMessage := "Session for this LearnR,(" + theLearnRUserSess.TheLearnR.Name + ") has been going on for too long,(" +
						strconv.Itoa(int(timeDuration)) + "). Killing User Session: " + strconv.Itoa(theLearnRUserSess.LocalSessID)
					theLearnRUserSess.LogInfo = append(theLearnRUserSess.LogInfo, theMessage)
					theLearnRUserSess.EndTime = time.Now()
					fmt.Println(theMessage)
					logWriter(theMessage)
					break
				}
			}
		}
		/* Done sending all texts for this LearnR. We can put this on the UserSessClose Chan and
		begin sending the results to User/CRUD DB */
		theLearnRUserSess.EndTime = time.Now()
		//Update Session
		theLearnRUserSess.Active = false
		theLearnRUserSess.TheSession.Ongoing = false //Stop the session from continuing
		theLearnRUserSess.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		UserSessionActiveMap[theLearnRUserSess.LocalSessID] = theLearnRUserSess
		theTimeNow = time.Now()
		//Add session information to session DB
		wg.Add(1)
		go fastAddLearnRSession(theLearnRUserSess.TheSession)
		//Update the LearnRInfo for this LearnR with our updated Session added to info
		theLearnRUserSess.TheLearnRInfo.AllSessions = append(theLearnRUserSess.TheLearnRInfo.AllSessions, theLearnRUserSess.TheSession)
		theLearnRUserSess.TheLearnRInfo.FinishedSessions = append(theLearnRUserSess.TheLearnRInfo.FinishedSessions, theLearnRUserSess.TheSession)
		theLearnRUserSess.TheLearnRInfo.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		wg.Add(1)
		go fastUpdateLearnRInform(theLearnRUserSess.TheLearnRInfo)
		wg.Wait() //Need to make sure we can exit this function properly
		//Removing Map placement of this UserSession
		delete(UserSessionActiveMap, theLearnRUserSess.LocalSessID)
		//Removing Phone Num from active UserSession
		delete(UserSessPhoneMap, theLearnRUserSess.PersonPhoneNum)
		//Print success
		fmt.Printf("This LearnR Session has now ended: %v for this LearnR: %v\n", theLearnRUserSess.LocalSessID, theLearnRUserSess.TheLearnR.Name)
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
		themessage := "Good response gotten: " + returnedMessage.Body
		resultMessages = append(resultMessages, themessage)
	} else {
		fmt.Println(resp.Status)
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("There's an error reading the response: %v\n", err.Error())
		} else {
			fmt.Printf("Here is the response: %v\n", string(b))
		}
		goodSend = false
		//Check to see if we have a 'black listed' number; end session if we can find it in our array
		type TwilioErrResponse struct {
			Code     int    `json:"code"`
			Message  string `json:"message"`
			MoreInfo string `json:"more_info"`
			Status   int    `json:"status"`
		}
		var twilErrorResponse TwilioErrResponse
		json.Unmarshal(b, &twilErrorResponse)
		if strings.Contains(strings.ToLower(twilErrorResponse.Message), strings.ToLower("blacklist")) {
			blacklistErr := "This number is on the 'blacklist' rule: " + toNumString + "...Stopping session for number"
			resultMessages = append(resultMessages, blacklistErr)
			fmt.Println(blacklistErr)
			/* stopping phone session for this number */
			fromTextedNum := strings.ReplaceAll(toNumString, "+", "")
			goodPhoneGet, userSessID := phoneSessionMapOk(fromTextedNum)
			if !goodPhoneGet {
				err := "Phone number, (" + toNumString + "), not found in session"
				logWriter(err)
				fmt.Println(err)
				fmt.Printf("DEBUG: Here is our map: %v\n and here is our user sess map: %v\n", UserSessPhoneMap, UserSessionActiveMap)
				resultMessages = append(resultMessages, err)
			} else {
				//Phone Session still active; need to cancel it
				okayuserSess, theUserSess := getSessionMapOk(userSessID)
				if !okayuserSess {
					msg := "Could not find an active User Session"
					logWriter(msg)
					fmt.Println(msg)
					resultMessages = append(resultMessages, msg)
				} else {
					//Found active User Session; determine what needs to be done with it
					responseNone := ResponseText{
						From: toNumString,
						Body: "STOP",
					}
					actOnTextSent(theUserSess, responseNone)
				}
			}
		} else if strings.Contains(strings.ToLower(twilErrorResponse.Message), strings.ToLower("not a valid")) {
			//Not a valid phone number...ending session
			notPhoneErr := "This number, " + toNumString + " is not a valid phone number. Stopping session for number"
			resultMessages = append(resultMessages, notPhoneErr)
			fmt.Println(notPhoneErr)
			/* stopping phone session for this number */
			fromTextedNum := strings.ReplaceAll(toNumString, "+", "")
			goodPhoneGet, userSessID := phoneSessionMapOk(fromTextedNum)
			if !goodPhoneGet {
				err := "Phone number, (" + toNumString + "), not found in session"
				logWriter(err)
				fmt.Println(err)
				fmt.Printf("DEBUG: Here is our map: %v\n and here is our user sess map: %v\n", UserSessPhoneMap, UserSessionActiveMap)
				resultMessages = append(resultMessages, err)
			} else {
				//Phone Session still active; need to cancel it
				okayuserSess, theUserSess := getSessionMapOk(userSessID)
				if !okayuserSess {
					msg := "Could not find an active User Session"
					logWriter(msg)
					fmt.Println(msg)
					resultMessages = append(resultMessages, msg)
				} else {
					//Found active User Session; determine what needs to be done with it
					responseNone := ResponseText{
						From: toNumString,
						Body: "STOP",
					}
					actOnTextSent(theUserSess, responseNone)
				}
			}
		}

		//Check to see if theErr is nil
		if theErr == nil {
			resultMessages = append(resultMessages, "The error returned from response is nil...")
		} else {
			resultMessages = append(resultMessages, theErr.Error())
		}
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

//Listens to all webhooks and determines next actions based on what is sent from form
func textWebhook(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //Parses our form first

	theReturnText := ResponseText{
		NumMedia:            r.FormValue("NumMedia"),
		FromCountry:         r.FormValue("FromCountry"),
		SmsStatus:           r.FormValue("SmsStatus"),
		ApiVersion:          r.FormValue("ApiVersion"),
		Body:                r.FormValue("Body"),
		To:                  r.FormValue("To"),
		ToCity:              r.FormValue("ToCity"),
		FromCity:            r.FormValue("FromCity"),
		AccountSid:          r.FormValue("AccountSid"),
		ToState:             r.FormValue("ToState"),
		From:                r.FormValue("From"),
		MessagingServiceSid: r.FormValue("MessagingServiceSid"),
		SmsSid:              r.FormValue("SmsSid"),
		MessageSid:          r.FormValue("MessageSid"),
		FromState:           r.FormValue("FromState"),
		ToZip:               r.FormValue("ToZip"),
		ToCountry:           r.FormValue("ToCountry"),
		FromZip:             r.FormValue("FromZip"),
		NumSegments:         r.FormValue("NumSegments"),
		OptOutType:          r.FormValue("OptOutType"),
	}

	/* Evaluate response;
	1. Check to see if the number this came from is in our map of session numbers
	2. Look up this number in our map and see if the session is still active,(or not over session time)
	3. Determine if this number has texted a keyword to stop the session
	*/
	fromTextedNum := strings.ReplaceAll(theReturnText.From, "+", "")
	okayPhone, userSessID := phoneSessionMapOk(fromTextedNum)
	if !okayPhone {
		err := "Phone number, (" + theReturnText.From + "), not found in session"
		logWriter(err)
		fmt.Println(err)
	} else {
		//Phone number exists, see if session is there and active
		okayuserSess, theUserSess := getSessionMapOk(userSessID)
		if !okayuserSess {
			msg := "Could not find an active User Session"
			logWriter(msg)
			fmt.Println(msg)
		} else {
			//Found active User Session; determine what needs to be done with it
			actOnTextSent(theUserSess, theReturnText)
		}
	}
}

//Gets our userSession ID from a phone number that texted the webhook
func phoneSessionMapOk(theNumber string) (bool, int) {
	fmt.Printf("DEBUG: Here is what our UserSessPhoneMap looks like...%v\n", UserSessPhoneMap)
	fromTextedNum := strings.ReplaceAll(theNumber, "+", "") //Get rid of plus just in case
	if theSessNum, ok := UserSessPhoneMap[fromTextedNum]; ok {
		return ok, theSessNum
	} else {
		return ok, theSessNum
	}
}

//Get a User Session from a sessionID
func getSessionMapOk(theSessID int) (bool, UserSession) {
	if theUserSess, ok := UserSessionActiveMap[theSessID]; ok {
		//Determine if this UserSess is still active and not past it's time limit
		goodSess := true
		duration := time.Since(theUserSess.StartTime).Seconds()
		if (int(duration) >= ALLOTTEDLEARNRTIME) || (!UserSessionActiveMap[theSessID].Active) {
			goodSess = false
		}

		if !goodSess {
			return goodSess, theUserSess
		} else {
			return ok, theUserSess
		}
	} else {
		return ok, theUserSess
	}
}

//Determine the text sent and act accordingly
func actOnTextSent(theUserSession UserSession, theUserResponse ResponseText) {
	theTimeNow := time.Now()
	//Set Body to stop if the userresponse contains the twilio 'OptOutType' as STOP
	if strings.Contains(strings.ToLower(theUserResponse.OptOutType), strings.ToLower("stop")) {
		theUserResponse.Body = "stop"
	}
	//Determine if the text is STOP or stop adjacent
	if val, ok := StopText[theUserResponse.Body]; ok {
		//It's a stop word, stop the session
		//Session update
		theUserSession.Active = false
		theUserSession.TheSession.Ongoing = false
		theUserSession.EndTime = time.Now()
		theUserSession.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		UserSessionActiveMap[theUserSession.LocalSessID] = theUserSession
	} else {
		//Not a stop; add to UserRSession
		//Session Update
		theUserSession.TheSession.UserResponses = append(theUserSession.TheSession.UserResponses, val)
		theUserSession.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
		//User Session Update
		logInfo := "User responsded with the following text: " + theUserResponse.Body + "\nTime: " + theUserSession.TheSession.DateUpdated
		theUserSession.LogInfo = append(theUserSession.LogInfo, logInfo)
		logWriter(logInfo)
		fmt.Println(logInfo)
		//Update the map
		UserSessionActiveMap[theUserSession.LocalSessID] = theUserSession
	}
}

/* TEST FUNCTIONS */
func httpTakerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Printf("DEBUG: Got to beginning of the function\n")
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

	fmt.Printf("Got to beginning of the taker...\n")
	//Get the byte slice from the request
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	type OurJSON struct {
		StringVal string `json:"StringVal"`
	}
	//Marshal it into our type
	var theJSON OurJSON
	json.Unmarshal(bs, &theJSON)

	/* Send the response back to Ajax */
	go func() {
		fmt.Printf("DEBUG: Got to go func\n")
		theJSONMessage, err := json.Marshal(theSuccMessage)
		//Send the response back
		if err != nil {
			errIs := "Error formatting JSON for return in testFunc: " + err.Error()
			panic(errIs)
		}
		theInt, theErr := fmt.Fprint(w, string(theJSONMessage))
		if theErr != nil {
			panic("Error writing back to initialLearnRStart: " + theErr.Error() + " " + strconv.Itoa(theInt))
		}
		flusher.Flush()
		fmt.Printf("DEBUG: We flushed \n")
	}()
	time.Sleep(time.Second * 1)
	go specialGoRoutine()
	fmt.Printf("Done with this function\n")
}

func specialGoRoutine() {
	fmt.Printf("Hey uhhhh, welcome to the random num generator...\n")
	for l := 0; l < 4; l++ {
		time.Sleep(time.Second * 5)
		rand.Seed(time.Now().UnixNano())
		min := 10
		max := 30
		fmt.Println(rand.Intn(max-min+1) + min)
	}
}
