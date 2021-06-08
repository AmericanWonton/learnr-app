package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

/* TEMPLATE DEFINITION */
var template1 *template.Template

//initial functions when starting the app
func init() {
	//Initialize our web page templates
	template1 = template.Must(template.ParseGlob("./static/templates/*"))
}

//Used for writing logs
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
	fmt.Printf("Beginning unit golang tests...\n ")
	j := Calculate(2)
	fmt.Printf("DEBUG: here is j: %v\n", j)
	logWriter("Test log message in main")
}

func Calculate(x int) (result int) {
	result = x + 2
	return result
}

func Add(x, y int) int {
	return x + y
}
