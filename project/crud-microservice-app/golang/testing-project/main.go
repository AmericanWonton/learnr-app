package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

//initial functions when starting the app
func init() {
	//Get Environment variables
	loadInMicroServiceURL()
	//Initialize Mongo Creds
	getCredsMongo()
}

func logWriter(logMessage string) {
	//Logging info

	wd, _ := os.Getwd()
	logDir := filepath.Join(wd, "logging", "crudapilog.txt ")
	logFile, err := os.OpenFile(logDir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

	defer logFile.Close()

	if err != nil {
		fmt.Println("Failed opening log file")
	}

	log.SetOutput(logFile)

	log.Println(logMessage)
}

func main() {
	fmt.Printf("DEBUG: Hello, we are listing with CRUD API\n") //Debug statement
	rand.Seed(time.Now().UTC().UnixNano())                     //Randomly Seed

	//Mongo Connect
	mongoClient = connectDB()
	defer mongoClient.Disconnect(theContext) //Disconnect in 10 seconds if you can't connect

	//Handle our incoming web requests
	handleRequests()
}

//Get mongo creds
func getCredsMongo() {
	//Check to see if ENV Creds are available first
	_, ok := os.LookupEnv("MONGO_URL")
	if !ok {
		message := "This ENV Variable is not present: " + "MONGO_URL"
		panic(message)
	}

	mongoURI = os.Getenv("MONGO_URL")
}
