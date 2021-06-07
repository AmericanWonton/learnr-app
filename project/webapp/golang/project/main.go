package main

import (
	"fmt"
	"math/rand"
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

func main() {
	fmt.Printf("DEBUG: Hello, we are in func main\n") //Debug statement
	rand.Seed(time.Now().UTC().UnixNano())            //Randomly Seed

	//Handle our incoming web requests
	handleRequests()
}
