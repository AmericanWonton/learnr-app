package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/* Both are used for usernames below */
var allUsernames []string
var usernameMap map[string]bool

/* Used for emails */
var emailMap map[string]bool

/* Used for displaying Learners */
var displayLearnrs []Learnr

var GETALLUSERNAMESURL string
var GETALLLEARNRORGURL string
var GETALLLEARNORGUSERADMIN string
var GETALLLEARNRURL string

//ViewData
type UserViewData struct {
	TheUser          User        `json:"TheUser"`          //The User
	Username         string      `json:"Username"`         //The Username
	Password         string      `json:"Password"`         //The Password
	Firstname        string      `json:"Firstname"`        //The First name
	Lastname         string      `json:"Lastname"`         //The Last name
	PhoneNums        []string    `json:"PhoneNums"`        //The Phone numbers
	UserID           int         `json:"UserID"`           //The UserID
	Email            []string    `json:"Email"`            //The Emails
	Whoare           string      `json:"Whoare"`           //Who is this person
	AdminOrgs        []int       `json:"AdminOrgs"`        //List of admin orgs
	OrgMember        []int       `json:"OrgMember"`        //List of organizations this Member is apart of
	AdminOrgList     []LearnrOrg `json:"AdminOrgList"`     //List of organization objects this User is Admin of(used on SOME pages)
	Banned           bool        `json:"Banned"`           //If the User is banned, we display nothing
	OrganizedLearnRs []Learnr    `json:"OrganizedLearnRs"` //An array of Learnrs with User input for ordering
	DateCreated      string      `json:"DateCreated"`      //Date this User was created
	DateUpdated      string      `json:"DateUpdated"`      //Date this User was updated
	MessageDisplay   int         `json:"MessageDisplay"`   //This is IF we need a message displayed
}

//Define pagehandler variables to Crud Microservice
func definePageHandlerVariables() {
	GETALLUSERNAMESURL = mongoCrudURL + "/giveAllUsernames"
	GETALLLEARNRORGURL = mongoCrudURL + "/giveAllLearnROrg"
	GETALLLEARNORGUSERADMIN = mongoCrudURL + "/getLearnOrgAdminOf"
	GETALLLEARNRURL = mongoCrudURL + "/giveAllLearnr"
}

//Handles the Index requests; Ask User if they're legal here
func index(w http.ResponseWriter, r *http.Request) {
	/* REdirect, Index not needed */
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

//Handles login/ page
func login(w http.ResponseWriter, r *http.Request) {
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "login.gohtml", nil)
	HandleError(w, err1)
}

//Handles the signup page
func signup(w http.ResponseWriter, r *http.Request) {
	usernameMap = loadUsernames() //Load all usernames
	emailMap = loadEmails()       //Load all Emails
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "signup.gohtml", nil)
	HandleError(w, err1)
}

//Handles the mainpage
func mainpage(w http.ResponseWriter, r *http.Request) {
	//Erase the learnrs loaded
	displayLearnrs = nil
	aUser := getUser(w, r)
	theLearnRs, goodGet, message := getSpecialLearnRs([]int{0, 1, 1, 1}, "", "", 0, 0)
	if !goodGet {
		logWriter("Issue getting Learnrs for this page: " + message)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	displayLearnrs = theLearnRs //Set Learnrs for display
	//Redirect User if they are not logged in
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	vd := UserViewData{
		TheUser:          aUser,
		Username:         aUser.UserName,
		Password:         aUser.Password,
		Firstname:        aUser.Firstname,
		Lastname:         aUser.Lastname,
		PhoneNums:        aUser.PhoneNums,
		UserID:           aUser.UserID,
		Email:            aUser.Email,
		Whoare:           aUser.Whoare,
		AdminOrgs:        aUser.AdminOrgs,
		OrgMember:        aUser.OrgMember,
		Banned:           aUser.Banned,
		OrganizedLearnRs: theLearnRs,
		DateCreated:      aUser.DateCreated,
		DateUpdated:      aUser.DateUpdated,
		MessageDisplay:   0,
	}
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
	HandleError(w, err1)
}

//Handles the learnmore page
func learnmore(w http.ResponseWriter, r *http.Request) {
	aUser := getUser(w, r)
	//Redirect User if they are not logged in
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	vd := UserViewData{
		TheUser:        aUser,
		Username:       aUser.UserName,
		UserID:         aUser.UserID,
		MessageDisplay: 0,
		Banned:         aUser.Banned,
	}
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "learnmore.gohtml", vd)
	HandleError(w, err1)
}

//Handles the sendhelp page
func sendhelp(w http.ResponseWriter, r *http.Request) {
	aUser := getUser(w, r)
	//Redirect User if they are not logged in
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	vd := UserViewData{
		TheUser:        aUser,
		Username:       aUser.UserName,
		UserID:         aUser.UserID,
		MessageDisplay: 0,
		Banned:         aUser.Banned,
	}
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "sendhelp.gohtml", vd)
	HandleError(w, err1)
}

//Handles the learnr page
func learnr(w http.ResponseWriter, r *http.Request) {
	learnrMap = loadLearnrs() //Get all our LearnR names for validation
	aUser := getUser(w, r)
	theAdminOrgs := loadLearnROrgArray(aUser)
	//Redirect User if they are not logged in
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	vd := UserViewData{
		TheUser:        aUser,
		Username:       aUser.UserName,
		UserID:         aUser.UserID,
		PhoneNums:      aUser.PhoneNums,
		Email:          aUser.Email,
		AdminOrgs:      aUser.AdminOrgs,
		MessageDisplay: 0,
		AdminOrgList:   theAdminOrgs,
		Banned:         aUser.Banned,
	}
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "learnr.gohtml", vd)
	HandleError(w, err1)
}

//Handles the makeorg page
func makeorg(w http.ResponseWriter, r *http.Request) {
	learnOrgMapNames = loadLearnROrgs()
	aUser := getUser(w, r)
	//Redirect User if they are not logged in
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	vd := UserViewData{
		TheUser:        aUser,
		Username:       aUser.UserName,
		UserID:         aUser.UserID,
		PhoneNums:      aUser.PhoneNums,
		Email:          aUser.Email,
		AdminOrgs:      aUser.AdminOrgs,
		MessageDisplay: 0,
		AdminOrgList:   []LearnrOrg{},
		Banned:         aUser.Banned,
	}
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "makeorg.gohtml", vd)
	HandleError(w, err1)
}

// Handle Errors passing templates
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		fmt.Printf("We had an error loading this template: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

//Calls 'giveAllUsernames' to run a mongo query to get all Usernames, then puts it in a map to return
func loadUsernames() map[string]bool {
	mapOusernameToReturn := make(map[string]bool) //Username to load our values into
	//Call our crudOperations Microservice in order to get our Usernames
	//Create a context for timing out
	resp, err := http.Get(GETALLUSERNAMESURL)
	if err != nil {
		theErr := "There was an error getting Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	//Marshal the response into a type we can read
	type ReturnMessage struct {
		TheErr          []string        `json:"TheErr"`
		ResultMsg       []string        `json:"ResultMsg"`
		SuccOrFail      int             `json:"SuccOrFail"`
		ReturnedUserMap map[string]bool `json:"ReturnedUserMap"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)

	//Assign our map variable to the map varialbe and see if it's okay
	if returnedMessage.SuccOrFail != 0 {
		errString := ""
		for l := 0; l < len(returnedMessage.TheErr); l++ {
			errString = errString + returnedMessage.TheErr[l]
		}
		logWriter(errString)
		fmt.Println(errString)
	} else {
		mapOusernameToReturn = returnedMessage.ReturnedUserMap
	}

	return mapOusernameToReturn
}

//Calls giveAllEmails to run a mongo query to get all Emails, then puts it in a map to return
func loadEmails() map[string]bool {
	mapOemailToReturn := make(map[string]bool) //Email to load our values into
	//Call our crudOperations Microservice in order to get our Emails
	//Create a context for timing out
	req, err := http.Get(GETUSEREMAILS)
	if err != nil {
		theErr := "There was an error getting Emails in loadEmails: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	if req.StatusCode >= 300 || req.StatusCode <= 199 {
		theErr := "There was an error reaching out to get Email API: " + strconv.Itoa(req.StatusCode)
		fmt.Println(theErr)
		logWriter(theErr)
	} else if err != nil {
		theErr := "Error from response to loadEmails: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "There was an error getting a response for Emails in loadEmails: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	//Marshal the response into a type we can read
	type ReturnMessage struct {
		TheErr           []string        `json:"TheErr"`
		ResultMsg        []string        `json:"ResultMsg"`
		SuccOrFail       int             `json:"SuccOrFail"`
		ReturnedEmailMap map[string]bool `json:"ReturnedEmailMap"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)

	//Assign our map variable to the map varialbe and see if it's okay
	if returnedMessage.SuccOrFail != 0 {
		errString := ""
		for l := 0; l < len(returnedMessage.TheErr); l++ {
			errString = errString + returnedMessage.TheErr[l]
		}
		logWriter(errString)
		fmt.Println(errString)
	} else {
		mapOemailToReturn = returnedMessage.ReturnedEmailMap
	}

	return mapOemailToReturn
}

//Calls 'giveAllLearnROrgs' to run a mongo query to get all LearnROrgs, then puts in a map to return
func loadLearnROrgs() map[string]bool {
	mapOLearnOrgsToReturn := make(map[string]bool) //LearnROrg map to load our values into
	//Call our crudOperations Microservice in order to get our Org Names
	//Create a context for timing out
	resp, err := http.Get(GETALLLEARNRORGURL)
	if err != nil {
		theErr := "There was an error getting LearnROrgs in loadLearnROrgs: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for LearnROrgs in loadLearnROrgs: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	//Marshal the response into a type we can read
	type ReturnMessage struct {
		TheErr             []string        `json:"TheErr"`
		ResultMsg          []string        `json:"ResultMsg"`
		SuccOrFail         int             `json:"SuccOrFail"`
		ReturnedOrgNameMap map[string]bool `json:"ReturnedOrgNameMap"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)

	//Assign our map variable to the map varialbe and see if it's okay
	if returnedMessage.SuccOrFail != 0 {
		errString := ""
		for l := 0; l < len(returnedMessage.TheErr); l++ {
			errString = errString + returnedMessage.TheErr[l]
		}
		logWriter(errString)
		fmt.Println(errString)
	} else {
		mapOLearnOrgsToReturn = returnedMessage.ReturnedOrgNameMap
	}
	return mapOLearnOrgsToReturn
}

//Calls 'giveAllLearnRs' to run a mongo query to get all LearnR, then put in a map to return
func loadLearnrs() map[string]bool {
	mapOLearnrsToReturn := make(map[string]bool) //LearnROrg map to load our values into
	//Call our crudOperations Microservice in order to get our Org Names
	//Create a context for timing out
	resp, err := http.Get(GETALLLEARNRURL)
	if err != nil {
		theErr := "There was an error getting Learnrs in loadLearnR: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for LearnR in loadLearnRs: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	//Marshal the response into a type we can read
	type ReturnMessage struct {
		TheErr              []string        `json:"TheErr"`
		ResultMsg           []string        `json:"ResultMsg"`
		SuccOrFail          int             `json:"SuccOrFail"`
		ReturnedLearnRNames map[string]bool `json:"ReturnedLearnRNames"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)

	//Assign our map variable to the map varialbe and see if it's okay
	if returnedMessage.SuccOrFail != 0 {
		errString := ""
		for l := 0; l < len(returnedMessage.TheErr); l++ {
			errString = errString + returnedMessage.TheErr[l]
		}
		logWriter(errString)
		fmt.Println(errString)
	} else {
		mapOLearnrsToReturn = returnedMessage.ReturnedLearnRNames
	}
	return mapOLearnrsToReturn
}

//Calls 'giveLearnROrgs' to run a mongo query to get all LearnOrgs this User is Admin of
func loadLearnROrgArray(aUser User) []LearnrOrg {
	type TheAdminOrgs struct {
		TheIDS []int `json:"TheIDS"`
	}
	theID := TheAdminOrgs{TheIDS: aUser.AdminOrgs}
	theJSONMessage, err := json.Marshal(theID)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		log.Fatal(err)
	}
	payload := strings.NewReader(string(theJSONMessage))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequest("POST", GETALLLEARNORGUSERADMIN, payload)
	if err != nil {
		theErr := "There was an error getting LearnROrgs in loadLearnROrgs: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "There was an error reaching out to loadLearnROrg API: " + strconv.Itoa(resp.StatusCode)
		fmt.Println(theErr)
		logWriter(theErr)
	} else if err != nil {
		theErr := "Error from response to loadLearnROrg: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for LearnROrgs in loadLearnROrgs: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	//Marshal the response into a type we can read
	type TheReturnMessage struct {
		TheErr            []string    `json:"TheErr"`
		ResultMsg         []string    `json:"ResultMsg"`
		SuccOrFail        int         `json:"SuccOrFail"`
		ReturnedLearnOrgs []LearnrOrg `json:"ReturnedLearnOrgs"`
	}
	var returnedMessage TheReturnMessage
	json.Unmarshal(body, &returnedMessage)

	arrayOReturn := returnedMessage.ReturnedLearnOrgs

	return arrayOReturn
}

//Called from Ajax; gives all the learnrs to Javascript for display
func giveAllLearnrDisplay(w http.ResponseWriter, r *http.Request) {
	//Declare Ajax return statements to be sent back
	type SuccessMSG struct {
		Message           string   `json:"Message"`
		SuccessNum        int      `json:"SuccessNum"`
		TheDisplayLearnrs []Learnr `json:"TheDisplayLearnrs"`
	}
	theSuccMessage := SuccessMSG{
		Message:           "Got all Learnrs",
		SuccessNum:        0,
		TheDisplayLearnrs: displayLearnrs,
	}

	//fmt.Printf("DEBUG: Here is our learnr display we are returning: %v\n", theSuccMessage.TheDisplayLearnrs)
	/* Send the response back to Ajax */
	theJSONMessage, err := json.Marshal(theSuccMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in giveAllLearnrDisplay: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}
