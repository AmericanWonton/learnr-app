package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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
	//Get environment variables
	loadInMicroServiceURL()
	UserSessionActiveMap = make(map[int]UserSession) //Make Map not crazy
	UserSessPhoneMap = make(map[string]int)
	StopText = make(map[string]string)
	//Initialize our bad phrases
	getbadWords()
	//Initialize our Twilio Cresd
	getTwilioCreds()
	//Get our stop text values
	fillStopText()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) //Randomly Seed

	//Handle our incoming web requests
	handleRequests()
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
