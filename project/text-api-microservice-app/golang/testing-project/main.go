package main

import (
	"bufio"
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
	"time"
)

/* DEFINED SLURS */
var slurs []string = []string{}

/* Microservice test ping definition */
var TESTMONGOMICROPING string = "http://localhost:4000/available"

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
	//Initialize our bad phrases
	getbadWords()
	//Initialize our Twilio Cresd
	getTwilioCreds()
	//Get our stop text values
	fillStopText()
	microsUp() //Tests if our Microservices are up
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

func microsUp() {
	wg.Add(1)
	go getCrudMicro()
	wg.Wait()
}

func getCrudMicro() {
	req, err := http.Get(TESTMONGOMICROPING)
	if err != nil {
		theErr := "There was an error getting crudMicro: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
		log.Fatal(theErr)
	}
	if !strings.Contains(strings.ToLower(req.Status), "200") {
		theErr := "Issue getting the response for CRUDMicro. Response: " + req.Status
		fmt.Println(theErr)
		logWriter(theErr)
		log.Fatal(theErr)
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "There was an error getting a response for crudMicro: " + err.Error()
		logWriter(theErr)
		fmt.Println(theErr)
		log.Fatal(theErr)
	}

	//Marshal the response into a type we can read
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)

	if returnedMessage.SuccOrFail != 0 {
		theErr := "Error getting correct response from crudAPI: " + strconv.Itoa(returnedMessage.SuccOrFail)
		fmt.Println(theErr)
		logWriter(theErr)
		log.Fatal(theErr)
	}

	wg.Done() //End waitgroup
}
