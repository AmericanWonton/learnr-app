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

//Admin token for making actions
type AdminToken struct {
	AdminID string `json:"AdminID"`
	UserID  int    `json:"UserID"`
}

/* Used for API Calls */
const GETRANDOMID string = "http://localhost:4000/randomIDCreationAPI"
const ADDUSERURL string = "http://localhost:4000/addUser"
const ADDLEARNRORGURL string = "http://localhost:4000/addLearnOrg"
const GETUSERLOGIN string = "http://localhost:4000/userLogin"

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
func checkUsername(w http.ResponseWriter, req *http.Request) {
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
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
