package main

import (
	"fmt"
	"log"
	"net/http"
)

/* Both are used for usernames below */
var allUsernames []string
var usernameMap map[string]bool

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
	fmt.Printf("DEBUG: Need to collect Username information and pass it into this template\n")
	/* Execute template, handle error */
	err1 := template1.ExecuteTemplate(w, "signup.gohtml", nil)
	HandleError(w, err1)
}

// Handle Errors passing templates
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		fmt.Printf("We had an error loading this template: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}
