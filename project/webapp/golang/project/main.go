package main

import (
	"bufio"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

/* TEMPLATE DEFINITION */
var template1 *template.Template

/* Template funcmap */
var funcMap = template.FuncMap{
	"uppercase": strings.ToUpper, //upperCase is a key we can call inside of the template html file
	"isAdmin":   isAdmin,         //Check to see if a User is an admin
}

/* Used for test calls to our MongoDB Microservice */
var TESTMONGOPOST string
var TESTMONGOGET string

//initial functions when starting the app
func init() {
	//Get Environment Variables
	loadInEmailCreds()
	loadInMicroServiceURL()
	defineAPIVariables()                     //Define variables in feildvalidation.go
	defineCrudVariables()                    //Define variables for mongoCrudOperations.go
	definePageHandlerVariables()             //Define varialbes for pagehandler.go
	usernameMap = make(map[string]bool)      //Clear all Usernames when loading so no problems are caused
	learnOrgMapNames = make(map[string]bool) //Clear all Org Names when loading so no problems are caused
	learnrMap = make(map[string]bool)        //Clear all Learnr Names when loading so no problems are caused
	emailMap = make(map[string]bool)         //Clear all Email names when loading so no problems are called
	//Initialize our web page templates
	template1 = template.Must(template.New("").Funcs(funcMap).ParseGlob("./static/templates/*"))
	//Initialize Mongo Creds
	getCredsMongo()
	//Initialize our bad phrases
	getbadWords()
	//Ping Test Crud Mongo
	//testPingMongoCRUD()
	//Initialize Emails
	OAuthGmailService()
}

func logWriter(logMessage string) {
	//Logging info

	wd, _ := os.Getwd()
	logDir := filepath.Join(wd, "logging", "weblog.txt")
	logFile, err := os.OpenFile(logDir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

	defer logFile.Close()

	if err != nil {
		fmt.Println("Failed opening log file")
	}

	log.SetOutput(logFile)

	log.Println(logMessage)
}

func main() {
	fmt.Printf("DEBUG: Hello, we are in func main\n") //Debug statement
	rand.Seed(time.Now().UTC().UnixNano())            //Randomly Seed

	//Handle our incoming web requests
	handleRequests()
}

//Get mongo creds
func getCredsMongo() {
	file, err := os.Open("security/mongocreds.txt")

	if err != nil {
		fmt.Printf("Trouble opening file for Mongo Credentials: %v\n", err.Error())
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	var text []string
	var bigMongoString64 string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	for x := 0; x < len(text); x++ {
		bigMongoString64 = bigMongoString64 + text[x]
	}

	sDecUsername, _ := b64.StdEncoding.DecodeString(text[1])
	sDecPWord, _ := b64.StdEncoding.DecodeString(text[2])
	sDB, _ := b64.StdEncoding.DecodeString(text[3])

	theUsername := string(sDecUsername)
	theUsername = strings.Replace(theUsername, "\n", "", 1)
	thePWord := string(sDecPWord)
	thePWord = strings.Replace(thePWord, "\n", "", 1)
	theDB := string(sDB)
	theDB = strings.Replace(theDB, "\n", "", 1)

	mongoURI = makeMongoString(theUsername, thePWord, theDB, text[0])
}

/* This makes the mongo string from our base64 encoded password, username,
and database name */
func makeMongoString(theUsername string, thePword string, theDB string, mongoString string) string {
	finishedMongoString := mongoString

	finishedMongoString = strings.Replace(finishedMongoString, "<USER_NAME>", theUsername, 1)
	finishedMongoString = strings.Replace(finishedMongoString, "<PASSWORD>", thePword, 1)
	finishedMongoString = strings.Replace(finishedMongoString, "<THE_DB>", theDB, 1)

	//fmt.Printf("DEBUG: Here's our decoded string: %v\n", finishedMongoString)

	return finishedMongoString
}

/* This pings our Mongo with a post and get to see if everything is up and reachable */
func testPingMongoCRUD() {
	//Set our enviornment variables
	TESTMONGOPOST = mongoCrudURL + "/testPingPost"
	TESTMONGOGET = mongoCrudURL + "/testPingGet"
	fmt.Printf("Making a ping to GET for MongoCrud: %v\n", TESTMONGOGET)
	//Start with Get
	//Call our crudOperations Microservice in order to get our Usernames
	//Create a context for timing out
	resp, err := http.Get(TESTMONGOGET)
	if err != nil {
		theErr := "Could not get successful response from mongo get"
		fmt.Println(theErr)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		theErr := "There was an error getting a response for MongoCRUD: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
	} else {
		defer resp.Body.Close()
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

	fmt.Printf("DEBUG: Here is response from Get to MongoCRUD: %v\n", returnedMessage)

	//Do Post as well
	type LoginData struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}
	testLogin := LoginData{Username: "MrUserName", Password: "MrPassword"}
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	/* 2. Marshal test case to JSON expect */
	theJSONMessage, err := json.Marshal(testLogin)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	/* 3. Create Post to JSON */
	payload := strings.NewReader(string(theJSONMessage))
	req2, err := http.NewRequest("POST", TESTMONGOPOST, payload)
	if err != nil {
		theErr := "There was an error posting User: " + err.Error()
		fmt.Println(theErr)
		logWriter(theErr)
	}
	req2.Header.Add("Content-Type", "application/json")
	/* 4. Get response from Post */
	resp2, err := http.DefaultClient.Do(req2.WithContext(ctx2))
	if err != nil {
		theErr := "Failed response from addUser: " + strconv.Itoa(resp.StatusCode) + " " + err.Error()
		logWriter(theErr)
	} else if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
		theErr := "Failed response from addUser: " + strconv.Itoa(resp.StatusCode)
		logWriter(theErr)
	}
	defer resp2.Body.Close()
	defer req2.Body.Close()
	//Declare message we expect to see returned
	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		theErr := "There was an error reading response from UserCreate " + err.Error()
		logWriter(theErr)
	}
	type ReturnMessage2 struct {
		TheErr          []string        `json:"TheErr"`
		ResultMsg       []string        `json:"ResultMsg"`
		SuccOrFail      int             `json:"SuccOrFail"`
		ReturnedUserMap map[string]bool `json:"ReturnedUserMap"`
	}
	var returnedMessage2 ReturnMessage2
	json.Unmarshal(body2, &returnedMessage2)

	fmt.Printf("DEBUG: Here is the response from CRUD URL Post: %v\n", returnedMessage2)
}
