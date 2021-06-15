package main

import (
	"bufio"
	b64 "encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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

//initial functions when starting the app
func init() {
	usernameMap = make(map[string]bool) //Clear all Usernames when loading so no problems are caused
	learnOrgMap = make(map[string]bool) //Clear all Org Names when loading so no problems are caused
	learnrMap = make(map[string]bool)   //Clear all Learnr Names when loading so no problems are caused
	//Initialize our web page templates
	template1 = template.Must(template.New("").Funcs(funcMap).ParseGlob("./static/templates/*"))
	//Initialize Mongo Creds
	getCredsMongo()
	//Initialize our bad phrases
	getbadWords()
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

	//Mongo Connect
	mongoClient = connectDB()
	defer mongoClient.Disconnect(theContext) //Disconnect in 10 seconds if you can't connect

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
