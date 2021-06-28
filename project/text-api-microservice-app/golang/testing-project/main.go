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
	"time"
)

/* DEFINED SLURS */
var slurs []string = []string{}

//Used for writing log messages
func logWriter(logMessage string) {
	//Logging info

	wd, _ := os.Getwd()
	logDir := filepath.Join(wd, "logging", "textapilog.txt")
	logFile, err := os.OpenFile(logDir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

	defer logFile.Close()

	if err != nil {
		fmt.Println("Failed opening log file")
	}

	log.SetOutput(logFile)

	log.Println(logMessage)
}

//Initial functions to run
func init() {
	UserSessionActiveMap = make(map[int]UserSession) //Make Map not crazy
	UserSessPhoneMap = make(map[string]int)
	StopText = make(map[string]string)
	//Initialize Mongo Creds
	getCredsMongo()
	//Initialize our bad phrases
	getbadWords()
	//Initialize our Twilio Cresd
	getTwilioCreds()
	//Get our stop text values
	fillStopText()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Randomly Seed

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
		logWriter("Trouble opening mongo credentials: " + err.Error())
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

	return finishedMongoString
}

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

/* fill our stop map */
func fillStopText() {
	StopText["stop"] = "stop"
	StopText["stp"] = "stp"
	StopText["stahp"] = "stahp"
	StopText["starhp"] = "starhp"
	StopText["STOP"] = "STOP"
}
