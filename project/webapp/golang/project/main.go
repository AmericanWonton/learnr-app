package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

/* TEMPLATE DEFINITION */
var template1 *template.Template

//initial functions when starting the app
func init() {
	//Initialize our web page templates
	template1 = template.Must(template.ParseGlob("./static/templates/*"))
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
