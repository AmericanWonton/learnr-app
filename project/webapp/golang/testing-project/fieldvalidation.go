package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/* Used for API Calls */
const GETRANDOMID string = "http://localhost:4000/randomIDCreationAPI"
const ADDUSERURL string = "http://localhost:4000/addUser"
const UPDATEURL string = "http://localhost:4000/updateUser"
const ADDLEARNRORGURL string = "http://localhost:4000/addLearnOrg"
const GETUSERLOGIN string = "http://localhost:4000/userLogin"

/* Used for LearnR/LearnR Org creation */
var allLearnROrgNames []string
var learnOrgMapNames map[string]bool

/* Used for LearnR creation */
var allLearnRNames []string
var learnrMap map[string]bool

/* DEFINED SLURS */
var slurs []string = []string{}

//This gets the slur words we check against in our username and
//text messages
func getbadWords() {
	file, err := os.Open("security/badphrases.txt")

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

	slurs = text
}

//Checks the Usernames after every keystroke
func checkUsername(w http.ResponseWriter, r *http.Request) {
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	sbs := string(bs)

	if len(sbs) <= 0 {
		fmt.Fprint(w, "TooShort")
	} else if len(sbs) > 20 {
		fmt.Fprint(w, "TooLong")
	} else if containsLanguage(sbs) {
		fmt.Fprint(w, "ContainsLanguage")
	} else {
		fmt.Fprint(w, usernameMap[sbs])
	}
}

//Checks the LearnROrg Names after every keystroke
func checkLearnROrgNames(w http.ResponseWriter, r *http.Request) {
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	sbs := string(bs)

	if len(sbs) <= 0 {
		fmt.Fprint(w, "TooShort")
	} else if len(sbs) > 25 {
		fmt.Fprint(w, "TooLong")
	} else if containsLanguage(sbs) {
		fmt.Fprint(w, "ContainsLanguage")
	} else {
		fmt.Fprint(w, learnOrgMapNames[sbs])
	}
}

//Checks the LearnR Names after every key stroke
func checkLearnRNames(w http.ResponseWriter, r *http.Request) {
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	sbs := string(bs)

	if len(sbs) <= 0 {
		fmt.Fprint(w, "TooShort")
	} else if len(sbs) > 40 {
		fmt.Fprint(w, "TooLong")
	} else if containsLanguage(sbs) {
		fmt.Fprint(w, "ContainsLanguage")
	} else {
		fmt.Fprint(w, learnrMap[sbs])
	}
}

//This checks the 'about LearnROrg' section after every keystroke. Also works for LearnR About
func checkOrgAbout(w http.ResponseWriter, r *http.Request) {
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	sbs := string(bs)

	if len(sbs) <= 0 {
		fmt.Fprint(w, "TooShort")
	} else if len(sbs) > 400 {
		fmt.Fprint(w, "TooLong")
	} else if containsLanguage(sbs) {
		fmt.Fprint(w, "ContainsLanguage")
	} else {
		fmt.Fprint(w, "okay")
	}
}

//Checks to see if the Username contains language
func containsLanguage(theText string) bool {
	hasLanguage := false
	textLower := strings.ToLower(theText)
	for i := 0; i < len(slurs); i++ {
		if strings.Contains(textLower, slurs[i]) {
			hasLanguage = true
			return hasLanguage
		}
	}
	return hasLanguage
}

/* Create User, if everything checks */
func createUser(w http.ResponseWriter, r *http.Request) {
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

	//Marshal it into our type
	var theUser User
	json.Unmarshal(bs, &theUser)

	// get form values
	username := theUser.UserName
	password := theUser.Password
	firstname := theUser.Firstname
	lastname := theUser.Lastname
	phonenums := theUser.PhoneNums
	email := theUser.Email
	whoare := theUser.Whoare

	/* perform Crud API here to insert the new User */
	//Declare new User to insert
	//Begin to add User to Mongo
	bsString := []byte(password)                  //Encode Password
	encodedString := hex.EncodeToString(bsString) //Encode Password Pt2
	theTimeNow := time.Now()

	/* First call to random ID API */
	goodIDGet, message, randomid := randomAPICall()
	if goodIDGet {
		newUser := User{
			UserName:    username,
			Password:    encodedString,
			Firstname:   firstname,
			Lastname:    lastname,
			PhoneNums:   phonenums,
			UserID:      randomid, //DEBUG VALUE
			Email:       email,
			Whoare:      whoare,
			AdminOrgs:   []int{},
			OrgMember:   []int{},
			Banned:      false,
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Attempt User Insert
		goodAdd, message := callAddUser(newUser)
		if goodAdd {
			theSuccMessage.Message = message
			theSuccMessage.SuccessNum = 0
		} else {
			theSuccMessage.Message = message
			theSuccMessage.SuccessNum = 1
		}
	} else {
		//Couldn't get random Numb
		theSuccMessage.Message = message
		theSuccMessage.SuccessNum = 1
	}
	/* Send the response back to Ajax */
	theJSONMessage, err := json.Marshal(theSuccMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in createUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* Tests to see if User can log in */
func canLogin(w http.ResponseWriter, r *http.Request) {
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	//Send a response back to Ajax after session is made
	type SuccessMSG struct {
		Message    string `json:"Message"`
		SuccessNum int    `json:"SuccessNum"`
	}
	theSuccMessage := SuccessMSG{}

	//Declare DataType from Ajax
	type LoginData struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}

	//Marshal the user data into our type
	var dataForLogin LoginData
	json.Unmarshal(bs, &dataForLogin)

	/* Call our CRUD API to see if password and Username are correct */
	goodLogin, message, returnedUser := callUserLogin(dataForLogin.Username, dataForLogin.Password)
	if goodLogin {
		theSuccMessage.SuccessNum = 0
		theSuccMessage.Message = "Successful User login"
		//Create User Session ID
		createSessionID(w, r, returnedUser)
	} else {
		theSuccMessage.SuccessNum = 1
		theSuccMessage.Message = "Failed User Login; Username/Password might not match"
		logWriter(message)
	}

	//Return JSON
	theJSONMessage, err := json.Marshal(theSuccMessage)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* Creates LearnrOrg, if everything checks */
func createLearnROrg(w http.ResponseWriter, r *http.Request) {
	//Declare Ajax return statements to be sent back
	type SuccessMSG struct {
		Message    string `json:"Message"`
		SuccessNum int    `json:"SuccessNum"`
	}
	theSuccMessage := SuccessMSG{
		Message:    "LearnR Organization created successfully",
		SuccessNum: 0,
	}

	//Declare struct we are expecting
	type SendJSON struct {
		TheLearnOrg LearnrOrg `json:"TheLearnOrg"`
		OurUser     User      `json:"OurUser"`
	}
	//Get the byte slice from the request
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	//Marshal it into our type
	var ourJSON SendJSON
	json.Unmarshal(bs, &ourJSON)

	fmt.Printf("DEBUG: Here is our UserArray before: %v\n", ourJSON.OurUser.AdminOrgs)

	/* perform Crud API here to insert the new User */
	theTimeNow := time.Now()

	/* First call to random ID API */
	goodIDGet, message, randomid := randomAPICall()
	if goodIDGet {
		newLearnROrg := LearnrOrg{
			OrgID:       randomid,
			Name:        ourJSON.TheLearnOrg.Name,
			OrgGoals:    ourJSON.TheLearnOrg.OrgGoals,
			UserList:    []int{ourJSON.OurUser.UserID},
			AdminList:   []int{ourJSON.OurUser.UserID},
			LearnrList:  []int{},
			DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Attempt User Insert
		goodAdd, message := calladdLearnOrg(newLearnROrg)
		if goodAdd {
			//Need to update our User as well
			updatedUser := User{
				UserName:    ourJSON.OurUser.UserName,
				Password:    ourJSON.OurUser.Password,
				Firstname:   ourJSON.OurUser.Firstname,
				Lastname:    ourJSON.OurUser.Lastname,
				PhoneNums:   ourJSON.OurUser.PhoneNums,
				UserID:      ourJSON.OurUser.UserID,
				Email:       ourJSON.OurUser.Email,
				Whoare:      ourJSON.OurUser.Whoare,
				AdminOrgs:   ourJSON.OurUser.AdminOrgs,
				OrgMember:   ourJSON.OurUser.OrgMember,
				Banned:      ourJSON.OurUser.Banned,
				DateCreated: ourJSON.OurUser.DateCreated,
				DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
			}
			updatedUser.AdminOrgs = append(updatedUser.AdminOrgs, randomid)
			updatedUser.OrgMember = append(updatedUser.OrgMember, randomid)
			fmt.Printf("DEBUG: Here is our updated User: %v\n", updatedUser)
			goodAdd2, message2 := callUpdateUser(updatedUser)

			if goodAdd2 {
				//Update our User Session too
				dbUsers[updatedUser.UserName] = updatedUser
				theSuccMessage.Message = message2
				theSuccMessage.SuccessNum = 0
			} else {
				theSuccMessage.Message = message2
				theSuccMessage.SuccessNum = 1
			}
		} else {
			theSuccMessage.Message = message
			theSuccMessage.SuccessNum = 1
		}
	} else {
		//Couldn't get random Numb
		theSuccMessage.Message = message
		theSuccMessage.SuccessNum = 1
	}
	/* Send the response back to Ajax */
	theJSONMessage, err := json.Marshal(theSuccMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in createUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* Creates a LearnR, if everything checks. Also
creates a learnRInfo to keep track of this LearnRs information overtime.
Finally, we'll update the LearnR Org with this new ID */
func createLearnR(w http.ResponseWriter, r *http.Request) {
	//Declare Ajax return statements to be sent back
	type SuccessMSG struct {
		Message    string `json:"Message"`
		SuccessNum int    `json:"SuccessNum"`
	}
	theSuccMessage := SuccessMSG{
		Message:    "LearnR Organization created successfully",
		SuccessNum: 0,
	}

	//Declare struct we are expecting
	type SendJSON struct {
		TheLearnr Learnr `json:"TheLearnr"`
		OurUser   User   `json:"OurUser"`
	}
	//Get the byte slice from the request
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	//Marshal it into our type
	var ourJSON SendJSON
	json.Unmarshal(bs, &ourJSON)

	/* Create basic learnr for data entry */
	theTimeNow := time.Now()
	goodIDGet, message, randomid := randomAPICall()
	if goodIDGet {
		theLearnr := Learnr{
			ID:            randomid,
			InfoID:        0,
			OrgID:         ourJSON.TheLearnr.OrgID,
			Name:          ourJSON.TheLearnr.Name,
			Tags:          ourJSON.TheLearnr.Tags,
			Description:   ourJSON.TheLearnr.Description,
			PhoneNums:     []string{},
			LearnRInforms: ourJSON.TheLearnr.LearnRInforms,
			Active:        true,
			DateCreated:   theTimeNow.Format("2006-01-02 15:04:05"),
			DateUpdated:   theTimeNow.Format("2006-01-02 15:04:05"),
		}
		//Create LearnrInfo for this LearnR
		goodIDGet, message, randomid := randomAPICall()
		if goodIDGet {
			//Create LearnR and add CRUD it to our DB
			theLearnRInfo := LearnrInfo{
				ID:               randomid,
				LearnRID:         theLearnr.ID,
				AllSessions:      []LearnRSession{},
				FinishedSessions: []LearnRSession{},
				DateCreated:      theTimeNow.Format("2006-01-02 15:04:05"),
				DateUpdated:      theTimeNow.Format("2006-01-02 15:04:05"),
			}
			goodAdd, message := callAddLearnrInfo(theLearnRInfo)
			if goodAdd {
				//LearnRInfo added successfully, add it to our LearnR
				theLearnr.InfoID = theLearnRInfo.ID
				//Fix LearnRInforms to have unique, correct values
				goodFixing := true //Determines if LearnRInforms successfully updated
				for n := 0; n < len(theLearnr.LearnRInforms); n++ {
					goodIDGet, _, randomid := randomAPICall()
					if goodIDGet {
						theLearnr.LearnRInforms[n].ID = randomid
						theLearnr.LearnRInforms[n].LearnRName = theLearnr.Name
						theLearnr.LearnRInforms[n].LearnRID = theLearnr.ID
						theLearnr.LearnRInforms[n].DateCreated = theTimeNow.Format("2006-01-02 15:04:05")
						theLearnr.LearnRInforms[n].DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
						theLearnr.LearnRInforms[n].Name = "LearnrInfo" + strconv.Itoa(n)
						theLearnr.LearnRInforms[n].Order = n
						goodAdd, _ := callAddLearnRInform(theLearnr.LearnRInforms[n])
						if !goodAdd {
							//Could not add to DB
							goodFixing = false
							break
						}
					} else {
						goodFixing = false
						break
					}
				}
				if goodFixing {
					//The Learnr has all values fixed/created, we can now add it to DB
					addLearnR, message := callAddLearnR(theLearnr)
					if addLearnR {
						//LearnR Added to DB; need to update the ORG it's under with our new ID
						theLearnOrgs := loadLearnROrgArray(ourJSON.OurUser)
						finalFixing := true    //Will determine if our learnorgs are updated correctly
						foundLearnOrg := false //Determines if learnorg is found and updated
						for j := 0; j < len(theLearnOrgs); j++ {
							if theLearnOrgs[j].OrgID == theLearnr.OrgID {
								updatedLearnROrg := theLearnOrgs[j]
								updatedLearnROrg.LearnrList = append(updatedLearnROrg.LearnrList, theLearnr.ID)
								updatedLearnROrg.DateUpdated = theTimeNow.Format("2006-01-02 15:04:05")
								goodUpdate, _ := callUpdateLearnOrg(updatedLearnROrg)
								if !goodUpdate {
									finalFixing = false
									break
								}
								foundLearnOrg = true
							}
						}
						if finalFixing && foundLearnOrg {
							//Return success
							theSuccMessage.SuccessNum = 0
							theSuccMessage.Message = "LearnR successfully added and all organizations updated"
						} else {
							theSuccMessage.SuccessNum = 1
							theSuccMessage.Message = "Failed to add LearnR to DB: "
						}
					} else {
						theSuccMessage.SuccessNum = 1
						theSuccMessage.Message = "Failed to add LearnR to DB: " + message
					}
				} else {
					theSuccMessage.SuccessNum = 1
					theSuccMessage.Message = "Failed to get proper ID: " + message
				}
			} else {
				theSuccMessage.SuccessNum = 1
				theSuccMessage.Message = "Failed to get proper ID: " + message
			}
		} else {
			theSuccMessage.SuccessNum = 1
			theSuccMessage.Message = "Failed to get proper ID: " + message
		}
	} else {
		theSuccMessage.SuccessNum = 1
		theSuccMessage.Message = "Failed to get proper ID: " + message
	}

	/* Send the response back to Ajax */
	theJSONMessage, err := json.Marshal(theSuccMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in createUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* Gets a random API after calling our random API */
func randomAPICall() (bool, string, int) {
	goodGet, message, finalInt := true, "", 0
	//Call our crudOperations Microservice in order to get our Usernames
	//Create a context for timing out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", GETRANDOMID, nil)
	if err != nil {
		theErr := "There was an error getting Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
		goodGet, message = false, theErr
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		goodGet, message = false, "Wrong response code gotten; failed to create random ID: "+strconv.Itoa(resp.StatusCode)
	} else if err != nil {
		theErr := "Had an error getting good random ID: " + err.Error()
		logWriter(theErr)
		goodGet, message = false, theErr
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
		goodGet, message = false, theErr
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
		goodGet, message = false, errString
	} else {
		finalInt = returnedMessage.RandomID
	}

	return goodGet, message, finalInt
}

/* Checks to see if this User is an admin. This is called from our gotemplate,
to see if User can create a learnR. If 0, they are an admin */
func isAdmin(aUser User) int {
	isAdmin := 1

	if len(aUser.AdminOrgs) > 0 {
		isAdmin = 0
	}

	return isAdmin
}

/* Calls to our Text API to see if the LearnR has started sending */
func canSendLearnR(w http.ResponseWriter, r *http.Request) {
	//Declare Ajax return statements to be sent back
	type SuccessMSG struct {
		Message    string `json:"Message"`
		SuccessNum int    `json:"SuccessNum"`
	}
	theSuccMessage := SuccessMSG{
		Message:    "LearnR sent successfully",
		SuccessNum: 0,
	}

	//Declare struct we are expecting
	type OurJSON struct {
		TheUser        User       `json:"TheUser"`
		TheLearnR      Learnr     `json:"TheLearnR"`
		TheLearnRInfo  LearnrInfo `json:"TheLearnRInfo"`
		PersonName     string     `json:"PersonName"`
		PersonPhoneNum string     `json:"PersonPhoneNum"`
		Introduction   string     `json:"Introduction"`
	}
	//Get the byte slice from the request
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	//Marshal it into our type
	var ourJSON OurJSON
	json.Unmarshal(bs, &ourJSON)

	/* Perform initial checks to see if the phone numbers, person name, and introduction are good */
	goodInitialCheck := true
	if (len(ourJSON.Introduction) < 1) || (len(ourJSON.Introduction) > 120) {
		goodInitialCheck = false
	}
	if (len(ourJSON.PersonName) < 1) || (len(ourJSON.PersonName) > 20) {
		goodInitialCheck = false
	}
	if (len(ourJSON.PersonPhoneNum) < 1) || (len(ourJSON.PersonPhoneNum) > 11) {
		goodInitialCheck = false
	}
	if (strings.Contains(ourJSON.PersonPhoneNum, "-")) || (strings.Contains(ourJSON.PersonPhoneNum, "+")) ||
		(strings.Contains(ourJSON.PersonPhoneNum, " ")) || (strings.Contains(ourJSON.PersonPhoneNum, ".")) ||
		(strings.Contains(ourJSON.PersonPhoneNum, ",")) {
		goodInitialCheck = false
	}

	if !goodInitialCheck {
		fmt.Printf("User has the wrong field: Intro: %v\nPersonName: %v\nPersonPhoneNum: %v\n", ourJSON.Introduction,
			ourJSON.PersonName, ourJSON.PersonPhoneNum)
		theSuccMessage.SuccessNum = 1
		theSuccMessage.Message = "Incorrect fields! Please look again"
	} else {
		/* Get the LearnRInfo assossiated with this LearnR */
		goodGet, result, theLearnRInfo := callReadLearnrInfo(ourJSON.TheLearnR.InfoID)
		if !goodGet {
			theErr := "Could not get proper LearnR information! " + result
			logWriter(theErr)
			fmt.Println(theErr)
			theSuccMessage.Message = theErr
			theSuccMessage.SuccessNum = 1
		} else {
			ourJSON.TheLearnRInfo = theLearnRInfo //Add LearnRInfo to JSON
			/* Check to see that our other values aren't nulled; this will cause a bad session if they are so... */
			if !(ourJSON.TheUser.UserID >= 1) || !(len(ourJSON.PersonName) >= 1) || !(len(ourJSON.PersonPhoneNum) >= 1 && len(ourJSON.PersonPhoneNum) <= 11) ||
				!(len(ourJSON.Introduction) >= 1) || !(ourJSON.TheLearnR.ID >= 1) {
				fmt.Printf("DEBUG: %v\n", ourJSON.TheUser.UserID)
				fmt.Printf("DEBUG: %v\n", ourJSON.PersonName)
				fmt.Printf("DEBUG: %v\n", ourJSON.PersonPhoneNum)
				fmt.Printf("DEBUG: %v\n", ourJSON.Introduction)
				fmt.Printf("DEBUG: %v\n", ourJSON.TheLearnR.ID)
				theErr := "Invalid values entered for this LearnR. Sending failed!"
				logWriter(theErr)
				fmt.Println(theErr)
				theSuccMessage.Message = theErr
				theSuccMessage.SuccessNum = 1
			} else {
				//Good check, go see if LearnR can be sent/started
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				/* 2. Marshal test case to JSON expect */
				theJSONMessage, err := json.Marshal(ourJSON)
				if err != nil {
					theErr := "Could not marshal JSON: " + err.Error()
					logWriter(theErr)
					fmt.Println(theErr)
					theSuccMessage.Message = theErr
					theSuccMessage.SuccessNum = 1
				}
				/* 3. Create Post to JSON */
				pingLocation := textAPIURL + "/initialLearnRStart"
				logWriter(string(theJSONMessage)) //Debug
				payload := strings.NewReader(string(theJSONMessage))
				req, err := http.NewRequest("POST", pingLocation, payload)
				if err != nil {
					theErr := "Error making request to Text API: " + err.Error()
					logWriter(theErr)
					fmt.Println(theErr)
					theSuccMessage.Message = theErr
					theSuccMessage.SuccessNum = 1
				}
				req.Header.Add("Content-Type", "application/json")
				defer req.Body.Close()
				/* 4. Get response from Post */
				resp, err := http.DefaultClient.Do(req.WithContext(ctx))
				if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
					theErr := "Failed response from initialLearnRStart: " + strconv.Itoa(resp.StatusCode)
					logWriter(theErr)
					theSuccMessage.Message = theErr
					theSuccMessage.SuccessNum = 1
				} else if err != nil {
					theErr := "Failed response from initialLearnRStart: " + strconv.Itoa(resp.StatusCode) + " " + err.Error()
					logWriter(theErr)
					theSuccMessage.Message = theErr
					theSuccMessage.SuccessNum = 1
				}
				defer resp.Body.Close()
				//Declare message we expect to see returned
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					theErr := "There was an error reading response from initialLearnRStart " + err.Error()
					logWriter(theErr)
					theSuccMessage.Message = theErr
					theSuccMessage.SuccessNum = 1
				}
				type TheSuccessMsg struct {
					Message    string `json:"Message"`
					SuccessNum int    `json:"SuccessNum"`
				}
				var returnedMessage TheSuccessMsg
				json.Unmarshal(body, &returnedMessage)
				/* 5. Evaluate response in returnedMessage */
				if returnedMessage.SuccessNum != 0 {
					theSuccMessage.SuccessNum = 1
					theSuccMessage.Message = "Failed to start convo with User!"
				} else {
					theSuccMessage.Message = "Text convo started with " + ourJSON.PersonName
				}
			}
		}
	}

	/* Send the response back to Ajax */
	theJSONMessage, err := json.Marshal(theSuccMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in createUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* Calls our CRUD API to narrow our search down */
func searchLearnRs(w http.ResponseWriter, r *http.Request) {
	//Declare Ajax return statements to be sent back
	type SuccessMSG struct {
		Message    string `json:"Message"`
		SuccessNum int    `json:"SuccessNum"`
	}
	theSuccMessage := SuccessMSG{
		Message:    "LearnR sent successfully",
		SuccessNum: 0,
	}

	//Declare struct we are expecting
	type SearchJSON struct {
		TheNameInput string `json:"TheNameInput"`
		TheTagInput  string `json:"TheTagInput"`
	}
	//Get the byte slice from the request
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	//Marshal it into our type
	var searchJSON SearchJSON
	json.Unmarshal(bs, &searchJSON)

	/* Send the response back to Ajax */
	theJSONMessage, err := json.Marshal(theSuccMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in createUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}
