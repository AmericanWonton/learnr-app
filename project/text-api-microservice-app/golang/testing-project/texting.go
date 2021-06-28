package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

/* DEBUG ping values */
var INITIALLEARNRSEND string = "http://localhost:3000/initialLearnRStart"

/* Credentials for Twilio...THESE NEED TO BE READ IN AT SOME POINT */
var accountSID string
var authToken string
var urlStr string = "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json"

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
		goodAdd, message := callAddLearnRSession(newLearnRSession)
		if !goodAdd {
			theErr := "Failure to get random API in session: " + message
			theSuccMessage.Message = theErr
			logWriter(theErr)
		} else {
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
}

func buttFunc(theSession UserSession) {
	learnSessChannel = make(chan UserSession, 3)
	learnSessResultChannel = make(chan UserSession, 3)
	fmt.Printf("DEBUG: Adding to Channel...\n")
	learnSessChannel <- theSession
	go learnRSession(learnSessChannel, learnSessResultChannel)

	//Log results of our channel jobs being completed
	for a := 0; a <= len(learnSessResultChannel); a++ {
		aUserSess := <-learnSessResultChannel
		aMessage := "We are done with this learnRSess: " + strconv.Itoa(aUserSess.TheSession.ID) + " for this LearnR: \n" +
			aUserSess.TheSession.LearnRName + "\n"
		logWriter(aMessage)
		fmt.Println(aMessage)
	}
	defer close(learnSessChannel)       //Close the channel when needed
	defer close(learnSessResultChannel) //Close this channel when needed
}

func learnRSession(userSessionChan <-chan UserSession, userSessCloseChan chan<- UserSession) {
	fmt.Printf("DEBUG: Hey, we've started up learnRSession\n")
	for a := range userSessionChan {
		fmt.Printf("DEBUG: Running through this learnRSession: %v\n", a.TheLearnR.Name)
		//Get this UserSession off the channel and into a good variable
		theUserSession := a
		//Add this User Session to our map of phone numbers
		UserSessPhoneMap[theUserSession.PersonPhoneNum] = theUserSession.LocalSessID
		/* First send the introduction Texts to this User; we will give them a
		160 second break period in order to type STOP; if they do not, our LearnR continues, until
		they enter STOP at any time */
		continueLearnR := true                                            //This will determine if we can send the rest of our texts
		UserSessionActiveMap[theUserSession.LocalSessID] = theUserSession //Add this User Session to our map
		//Send the Introduction Text
		introMessage := "Hello " + theUserSession.PersonName + ", " + theUserSession.TheUserName + " wanted to help educate you on " +
			"something important to them."
		theTimeNow := time.Now()
		goodSend, resultMessages := sendText(-3, theUserSession.PersonPhoneNum, theUserSession.TheLearnR.PhoneNums[0],
			introMessage)
		if !goodSend || !UserSessionActiveMap[theUserSession.LocalSessID].Active {
			//Intro text failed...LearnR may not be active
			continueLearnR = false
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
			//Send the second text with our Users message
			time.Sleep(time.Second * 10) //Small wait
			introMessage = "\"" + theUserSession.IntroductionSaying + "\""
			goodSend, resultMessages := sendText(-2, theUserSession.PersonPhoneNum, theUserSession.TheLearnR.PhoneNums[0],
				introMessage)
			if !goodSend || !UserSessionActiveMap[theUserSession.LocalSessID].Active {
				//Intro text failed...LearnR may not be active
				continueLearnR = false
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
				//Send final message asking User to confirm/deny
				//Send the second text with our Users message
				time.Sleep(time.Second * 10) //Small wait
				introMessage = "At their request, we'll begin sending a LearnR to help them explain.\n\n" +
					"If you'd like to stop, please text STOP back to this number at any time."
				goodSend, resultMessages := sendText(-1, theUserSession.PersonPhoneNum, theUserSession.TheLearnR.PhoneNums[0],
					introMessage)
				if !goodSend || !UserSessionActiveMap[theUserSession.LocalSessID].Active {
					//Intro text failed...LearnR may not be active
					continueLearnR = false
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
			}
		}
		time.Sleep(time.Second * 10) //Small wait
		/* Hopefully we've sent our first three texts successfully, (continueLearnR == true).
		If not, log the failure remove this Session/Phone Map recording and update our DB*/
		if !continueLearnR || !UserSessionActiveMap[theUserSession.LocalSessID].Active {
			//Failed sending messages...ending session
			theMessage := "LearnR Session ending; did not have success sending texts or User terminated session..."
			fmt.Println(theMessage)
			logWriter(theMessage)
			theUserSession.Active = false
			UserSessionActiveMap[theUserSession.LocalSessID] = theUserSession
			theTimeNow = time.Now()
			theUserSession.TheSession.Ongoing = false //Stop the session from continuing
			theUserSession.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
			//Add session information to session DB
			fmt.Printf("DEBUG: Adding learnrSession to learnRSession DB\n")
			//Update the LearnRInfo for this LearnR with our updated Session added to it
			fmt.Printf("DEBUG: Adding our leanrrSession to our LearnRInform and updating it")
			//Removing Map placement of this UserSession
			delete(UserSessionActiveMap, theUserSession.LocalSessID)
			//Removing Phone Num from active UserSession
			delete(UserSessPhoneMap, theUserSession.PersonPhoneNum)
			//Add this to our Results Chan
			userSessCloseChan <- a
			fmt.Printf("We are now done texting this User\n")
		} else {
			//First three messages sent, getting ready to send the rest of the messages...
			//Start sending texts on this session
			for l := 0; l < len(theUserSession.TheLearnR.LearnRInforms); l++ {
				//Do this only if the 'Ongoing' variable in the LearnRSession is true
				if theUserSession.TheSession.Ongoing && UserSessionActiveMap[theUserSession.LocalSessID].Active {
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
						theUserSession.TheSession.TextsSent = append(theUserSession.TheSession.TextsSent, theUserSession.TheLearnR.LearnRInforms[l])
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
			theUserSession.Active = false
			UserSessionActiveMap[theUserSession.LocalSessID] = theUserSession
			theTimeNow = time.Now()
			theUserSession.TheSession.Ongoing = false //Stop the session from continuing
			theUserSession.TheSession.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
			//Add session information to session DB
			fmt.Printf("DEBUG: Adding learnrSession to learnRSession DB\n")
			//Update the LearnRInfo for this LearnR with our updated Session added to it
			fmt.Printf("DEBUG: Adding our leanrrSession to our LearnRInform and updating it")
			//Removing Map placement of this UserSession
			delete(UserSessionActiveMap, theUserSession.LocalSessID)
			//Removing Phone Num from active UserSession
			delete(UserSessPhoneMap, theUserSession.PersonPhoneNum)
			//Add this to our Results Chan
			userSessCloseChan <- a
			fmt.Printf("We are now done texting this User\n")
		}
	}
}

func conductLearnRSession(theLearnRUserSess UserSession) {
	fmt.Printf("DEBUG: Running through this learnRSession: %v\n", theLearnRUserSess.TheLearnR.Name)
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
	fmt.Printf("DEBUG: Starting with first text\n")
	goodSend, resultMessages := sendText(-3, theLearnRUserSess.PersonPhoneNum, theLearnRUserSess.TheLearnR.PhoneNums[0],
		introMessage)
	fmt.Printf("DEBUG: ending with first text\n")
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
		fmt.Printf("DEBUG: Adding learnrSession to learnRSession DB\n")
		//Update the LearnRInfo for this LearnR with our updated Session added to it
		fmt.Printf("DEBUG: Adding our leanrrSession to our LearnRInform and updating it")
		//Removing Map placement of this UserSession
		delete(UserSessionActiveMap, theLearnRUserSess.LocalSessID)
		//Removing Phone Num from active UserSession
		delete(UserSessPhoneMap, theLearnRUserSess.PersonPhoneNum)
		fmt.Printf("DEBUG: We are now done texting this User\n")
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
		fmt.Printf("DEBUG: Adding learnrSession to learnRSession DB\n")
		//Update the LearnRInfo for this LearnR with our updated Session added to it
		fmt.Printf("DEBUG: Adding our leanrrSession to our LearnRInform and updating it")
		//Removing Map placement of this UserSession
		delete(UserSessionActiveMap, theLearnRUserSess.LocalSessID)
		//Removing Phone Num from active UserSession
		delete(UserSessPhoneMap, theLearnRUserSess.PersonPhoneNum)
		fmt.Printf("We are now done texting this User\n")
	}
}

func sendText(textOrder int, toNumString string, fromNumString string, textBody string) (bool, []string) {
	goodSend, resultMessages := true, []string{}
	msgData := url.Values{}
	msgData.Set("To", "+"+toNumString)
	msgData.Set("From", "+"+fromNumString)
	msgData.Set("Body", textBody)
	msgDataReader := *strings.NewReader(msgData.Encode())

	fmt.Printf("DEBUG: Here is our msgDataReader: %v\n", msgDataReader)

	req, err := http.NewRequest("POST", urlStr, &msgDataReader)
	if err != nil {
		fmt.Printf("Error making newRequest: %v\n", err.Error())
	}
	fmt.Printf("DEBUG: Here is our account SID: %v\nAnd here is our AuthToken: %v\n", accountSID, authToken)
	req.SetBasicAuth(accountSID, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	fmt.Printf("DEBUG: Starting the request here\n")

	type TwilioGetBack struct {
		Code     int    `json:"Code"`
		Message  string `json:"Message"`
		MoreInfo string `json:"MoreInfo"`
		Status   int    `json:"Status"`
	}
	//Get that request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, theErr := http.DefaultClient.Do(req.WithContext(ctx))
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

	fmt.Printf("DEBUG: Got to the response here\n")
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
