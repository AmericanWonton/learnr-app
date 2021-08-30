package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"gopkg.in/mgo.v2/bson"
)

/* Both are used for usernames below */
var allUsernames []string
var usernameMap map[string]bool

/* Used for emails */
var emailMap map[string]bool

/* Used for displaying Learners */
var displayLearnrs []Learnr
var newDisplay int = 0 //Used if we need to show a returned group of Users

const GETALLUSERNAMESURL string = "http://localhost:4000/giveAllUsernames"
const GETALLLEARNRORGURL string = "http://localhost:4000/giveAllLearnROrg"
const GETALLLEARNORGUSERADMIN string = "http://localhost:4000/getLearnOrgAdminOf"
const GETALLLEARNRURL string = "http://localhost:4000/giveAllLearnr"

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
	NewSearchLearnR  int         `json:"NewSearchLearnR"`  //Determines if new LearnRs are being displayed in Agular. 0 == no
	OrganizedLearnRs []Learnr    `json:"OrganizedLearnRs"` //An array of Learnrs with User input for ordering
	DateCreated      string      `json:"DateCreated"`      //Date this User was created
	DateUpdated      string      `json:"DateUpdated"`      //Date this User was updated
	MessageDisplay   int         `json:"MessageDisplay"`   //This is IF we need a message displayed
	UserMessage      string      `json:"UserMessage"`      //The Message displayed to our User
	ActionDisplay    int         `json:"ActionDisplay"`    //A condition for displaying various things to our Users
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
	if newDisplay == 1 {
		aUser := getUser(w, r)
		//Redirect User if they are not logged in
		if !alreadyLoggedIn(w, r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		/* Set view data based on new learnrs we searched for */
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
			NewSearchLearnR:  newDisplay,
			OrganizedLearnRs: displayLearnrs,
			DateCreated:      aUser.DateCreated,
			DateUpdated:      aUser.DateUpdated,
			MessageDisplay:   0,
		}
		/* Execute template, handle error */
		err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
		HandleError(w, err1)
	} else {
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
			NewSearchLearnR:  newDisplay,
			OrganizedLearnRs: theLearnRs,
			DateCreated:      aUser.DateCreated,
			DateUpdated:      aUser.DateUpdated,
			MessageDisplay:   0,
		}
		/* Execute template, handle error */
		err1 := template1.ExecuteTemplate(w, "mainpage.gohtml", vd)
		HandleError(w, err1)
	}
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

//Handles the bulksend page
func bulksend(w http.ResponseWriter, r *http.Request) {
	aUser := getUser(w, r)
	//Redirect User if they are not logged in
	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	theAdminOrgs := loadLearnROrgArray(aUser) //Get the admin orgs for selection
	//Get all LearnRIds collected
	theLearnRIDs := []int{}
	for n := 0; n < len(theAdminOrgs); n++ {
		theLearnRIDs = append(theLearnRIDs, theAdminOrgs[n].LearnrList...)
	}
	theAdminLearnRs := getAdminLearnRs(theLearnRIDs)
	/* Tailor the message towards our User, based on what access/how many
	LearnROrgs they are Admins of */
	userMessage := ""
	shouldDisplay := 0
	actionDisplay := 0
	if len(theAdminLearnRs) <= 0 || theAdminLearnRs == nil {
		userMessage = "You do not have any LearnRs you are Admin of to send in Bulk! Please create a LearnR under a " +
			"LearnROrg that you are an Administrator of..."
		shouldDisplay = 1
		actionDisplay = 1
		fmt.Println(userMessage)
	}
	vd := UserViewData{
		TheUser:          aUser,
		Username:         aUser.UserName,
		UserID:           aUser.UserID,
		PhoneNums:        aUser.PhoneNums,
		Email:            aUser.Email,
		AdminOrgs:        aUser.AdminOrgs,
		MessageDisplay:   shouldDisplay,
		AdminOrgList:     theAdminOrgs,
		OrganizedLearnRs: theAdminLearnRs,
		Banned:           aUser.Banned,
		UserMessage:      userMessage,
		ActionDisplay:    actionDisplay,
	}

	/* If this is an HTTP Post, determine if we need to load this page differently*/
	if r.Method == http.MethodPost {
		/* Define stuff to return */
		type TheSuccessMsg struct {
			Message    string `json:"Message"`
			SuccessNum int    `json:"SuccessNum"`
		}
		theSuccMessage := TheSuccessMsg{
			Message:    "Bulk LearnR Successfully started",
			SuccessNum: 0,
		}

		hiddenFormValue := r.FormValue("hiddenFormValue")
		if strings.Contains(strings.ToLower(hiddenFormValue), strings.ToLower("bulk-excel")) {
			//Good value, continue working this Excel sheet
			maxSize := int64(1024000) // allow only 1MB of file size
			err := r.ParseMultipartForm(maxSize)
			if err != nil {
				theErr := "File too large. Max Size: " + strconv.Itoa(int(maxSize)) + "mb " + err.Error()
				fmt.Println(theErr)
				log.Println(err)
				vd.MessageDisplay = 1
				vd.UserMessage = theErr
				vd.ActionDisplay = 0
				theSuccMessage.SuccessNum = 1
				theSuccMessage.Message = theErr
			} else {
				//File okay, continue on
				file, fileHeader, err := r.FormFile("excel-file") //Insert name of file element here
				if err != nil {
					errMsg := "Error getting file submission: " + err.Error()
					fmt.Println(errMsg)
					vd.MessageDisplay = 1
					vd.UserMessage = errMsg
					vd.ActionDisplay = 0
					theSuccMessage.SuccessNum = 1
					theSuccMessage.Message = errMsg
				} else {
					//Good file form, moving on
					//Create path and write file on server
					hexName := bson.NewObjectId().Hex()
					fileExtension := filepath.Ext(fileHeader.Filename)
					theFileName := hexName + fileExtension
					theDir, _ := os.Getwd()
					thePath := filepath.Join(theDir, "tempFiles")
					os.MkdirAll(thePath, 0777)
					//Write file on server
					f, err := os.OpenFile(theFileName, os.O_WRONLY|os.O_CREATE, 0777)
					if err != nil {
						theErr := "Error opening Excel file: " + err.Error()
						fmt.Println(theErr)
						log.Fatal(theErr)
					}
					io.Copy(f, file)
					f.Close()
					file.Close()
					//Move file to folder
					thePath2 := filepath.Join(theDir, "tempFiles", theFileName)
					readFile, err := os.Open(theFileName)
					if err != nil {
						errMsg := "STEP 2: Error opening this file: " + err.Error()
						fmt.Println(errMsg)
						log.Fatal(errMsg)
					}
					writeToFile, err := os.Create(thePath2)
					if err != nil {
						errMsg := "STEP 3 Error creating writeToFile: " + err.Error()
						fmt.Println(errMsg)
						log.Fatal(errMsg)
					}
					//Move file Contents to folder
					_, err3 := io.Copy(writeToFile, readFile)
					if err3 != nil {
						errMsg := "PART 4 Error copying the contents of the one image to the other: " + err3.Error()
						log.Fatal(errMsg)
					}
					readFile.Close()    //Close File
					writeToFile.Close() //Close File
					//Delete created file
					removeErr := os.Remove(theFileName)
					if removeErr != nil {
						errMsg := "STEP 5 Error removing the file: " + removeErr.Error()
						fmt.Println(errMsg)
						log.Fatal(errMsg)
					}

					/* Analayze Excel sheet to determine if this Excel sheet is formattted okay */
					goodExcel, message := examineExcelSheet(thePath2, theFileName)
					if !goodExcel {
						fmt.Println(message)
						vd.MessageDisplay = 1
						vd.UserMessage = message
						vd.ActionDisplay = 0
						theSuccMessage.SuccessNum = 1
						theSuccMessage.Message = message
					} else {
						//Send this to Amazon buckets for storage
						//Create Amazon Session
						s, err := session.NewSession(&aws.Config{
							Region: aws.String("us-east-2"),
							Credentials: credentials.NewStaticCredentials(
								AWSAccessKeyId, // id
								AWSSecretKey,   // secret
								""),            // token can be left blank for now
						})
						if err != nil {
							errMsg := "STEP 6 Could not upload file. Error creating session: " + err.Error()
							vd.UserMessage = errMsg
							vd.ActionDisplay = 0
							vd.MessageDisplay = 1
							fmt.Println(errMsg)
							logWriter(errMsg)
							theSuccMessage.SuccessNum = 1
							theSuccMessage.Message = "Error uplodaing Excel file; contact Admin step 6"
						} else {
							goodAmazon, theMessage, amazonLocation := sendExcelToBucket(thePath2, hexName,
								s, file, fileHeader, aUser)
							if !goodAmazon {
								vd.UserMessage = theMessage
								vd.ActionDisplay = 0
								vd.MessageDisplay = 1
								fmt.Println(theMessage)
								theSuccMessage.SuccessNum = 1
								theSuccMessage.Message = theMessage
							} else {
								/* Good Excel sheet sending to Amazon, now we can see
								if we can get that bulk load started */
								learnRFormValue := r.FormValue("learnR")
								learnRID, _ := strconv.Atoi(learnRFormValue)
								goodSend, message := canSendBulkLearnR(aUser, amazonLocation, theFileName, learnRID)
								if !goodSend {
									errMsg := "There was an issue starting the Bulk LearnR: " + message
									vd.UserMessage = errMsg
									vd.ActionDisplay = 0
									vd.MessageDisplay = 1
									fmt.Println(errMsg)
									logWriter(errMsg)
									theSuccMessage.SuccessNum = 1
									theSuccMessage.Message = errMsg
								} else {
									//Bulk LearnR successfully started
									goodMsg := "Bulk LearnR successfully started"
									vd.UserMessage = goodMsg
									vd.ActionDisplay = 1
									vd.MessageDisplay = 1
									theSuccMessage.Message = goodMsg
								}
							}
						}
					}
				}
			}
		} else {
			theSuccMessage.SuccessNum = 0
			theSuccMessage.Message = "Error getting special values for form. Please contact Admin."
		}

		/* Send the response back to Ajax */
		theJSONMessage, err := json.Marshal(theSuccMessage)
		//Send the response back
		if err != nil {
			errIs := "Error formatting JSON for return in submitting file: " + err.Error()
			fmt.Println(errIs)
			logWriter(errIs)
		}
		fmt.Fprint(w, string(theJSONMessage))
	} else {
		/* Execute template, handle error */
		err1 := template1.ExecuteTemplate(w, "bulksend.gohtml", vd)
		HandleError(w, err1)
	}
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", GETALLUSERNAMESURL, nil)
	if err != nil {
		theErr := "There was an error getting Usernames in loadUsernames: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "There was an error reaching out to loadUsername API: " + strconv.Itoa(resp.StatusCode)
		fmt.Println(theErr)
		logWriter(theErr)
	} else if err != nil {
		theErr := "Error from response to loadUsernames: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
	}
	defer resp.Body.Close()

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", GETALLLEARNRORGURL, nil)
	if err != nil {
		theErr := "There was an error getting LearnROrgs in loadLearnROrgs: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", GETALLLEARNRURL, nil)
	if err != nil {
		theErr := "There was an error getting Learnrs in loadLearnR: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))

	if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "There was an error reaching out to loadLearnr API: " + strconv.Itoa(resp.StatusCode)
		fmt.Println(theErr)
		logWriter(theErr)
	} else if err != nil {
		theErr := "Error from response to loadLearnr: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
	}
	defer resp.Body.Close()

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

	/* Send the response back to Ajax */
	theJSONMessage, err := json.Marshal(theSuccMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in giveAllLearnrDisplay: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//Called from 'bulksend' to get all or LearnRs
func getAdminLearnRs(theLearnRIDs []int) []Learnr {
	/* Only look for admin learnRs if the passed array is larger than 0 */
	if len(theLearnRIDs) >= 1 {
		goodGet, message, returnedLearnRs := callReadLearnRArray(theLearnRIDs)
		if !goodGet {
			theErr := "Error getting array of LearnRs: " + message
			fmt.Println(theErr)
			logWriter(theErr)
		}
		return returnedLearnRs
	} else {
		return []Learnr{}
	}
}

/* This is called from mainpageapp.js to get our learnR data
to load to an Angular app */
func getLearnRAngular(w http.ResponseWriter, r *http.Request) {
	type ReturnMessage struct {
		TheErr      string   `json:"TheErr"`
		ResultMsg   string   `json:"ResultMsg"`
		SuccOrFail  int      `json:"SuccOrFail"`
		LearnRArray []Learnr `json:"LearnRArray"`
	}
	returnMessage := ReturnMessage{
		TheErr:      "",
		ResultMsg:   "LearnRs returned successfully",
		SuccOrFail:  0,
		LearnRArray: []Learnr{},
	}

	displayLearnrs = nil
	theLearnRs, goodGet, message := getSpecialLearnRs([]int{0, 1, 1, 1}, "", "", 0, 0)
	if !goodGet {
		theErr := "Issue getting Learnrs for this page: " + message
		logWriter(theErr)
		returnMessage.TheErr = theErr
		returnMessage.SuccOrFail = 1
		returnMessage.ResultMsg = theErr
		returnMessage.LearnRArray = theLearnRs
	} else {
		displayLearnrs = theLearnRs //Set Learnrs for display
		newDisplay = 1              //Set new learnrdisplay method
		returnMessage.LearnRArray = displayLearnrs
		returnMessage.SuccOrFail = 0
		returnMessage.ResultMsg = "Learnrs successfully retrieved"
	}

	/* Send the response back to Ajax */
	theJSONMessage, err := json.Marshal(returnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in giveAllLearnrDisplay: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* This is called from mainpageapp.js to get special search learnR data.
It then deletes any returned LearnRs and sets NewSearchLearnR back to false for another search*/
func getSpecialLearnRAngular(w http.ResponseWriter, r *http.Request) {
	type ReturnMessage struct {
		TheErr      string   `json:"TheErr"`
		ResultMsg   string   `json:"ResultMsg"`
		SuccOrFail  int      `json:"SuccOrFail"`
		LearnRArray []Learnr `json:"LearnRArray"`
	}
	returnMessage := ReturnMessage{
		TheErr:      "",
		ResultMsg:   "LearnRs returned successfully",
		SuccOrFail:  0,
		LearnRArray: []Learnr{},
	}
	//Assign LearnRs to return message
	returnMessage.LearnRArray = displayLearnrs
	//Delete 'searched for' LearnRs
	displayLearnrs = nil
	newDisplay = 0

	/* Send the response back to Ajax */
	theJSONMessage, err := json.Marshal(returnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in giveAllLearnrDisplay: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}
