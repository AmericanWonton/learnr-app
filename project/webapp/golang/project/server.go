package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var mongoCrudURL string
var textAPIURL string

//Handles all requests coming in
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	http.Handle("/favicon.ico", http.NotFoundHandler()) //For missing FavIcon
	//Serve our pages
	myRouter.HandleFunc("/", index)              //Serve index page
	myRouter.HandleFunc("/login", login)         //Serve login page
	myRouter.HandleFunc("/signup", signup)       //Serve signup page
	myRouter.HandleFunc("/mainpage", mainpage)   //Serve main page
	myRouter.HandleFunc("/learnmore", learnmore) //Serve learnmore page
	myRouter.HandleFunc("/sendhelp", sendhelp)   //Serve the sendhelp page
	myRouter.HandleFunc("/learnr", learnr)       //Serve the learnr page
	myRouter.HandleFunc("/makeorg", makeorg)     //Serve the learnr page
	//Used for handling emails
	myRouter.HandleFunc("/emailMe", emailMe).Methods("POST") //Used for email Sending from Users
	//Used for session work
	myRouter.HandleFunc("/logUserOut", logUserOut).Methods("POST") //Remove our cookie after logging out user
	//Serve our Validation API
	myRouter.HandleFunc("/checkUsername", checkUsername).Methods("POST")             //Check Username
	myRouter.HandleFunc("/checkLearnROrgNames", checkLearnROrgNames).Methods("POST") //Check LearnROrg Name
	myRouter.HandleFunc("/checkLearnRNames", checkLearnRNames).Methods("POST")       //Check Learnr Name
	myRouter.HandleFunc("/checkEmail", checkEmail).Methods("POST")                   //Check Check Email
	myRouter.HandleFunc("/checkOrgAbout", checkOrgAbout).Methods("POST")             //Check LearnOrg About
	myRouter.HandleFunc("/createLearnROrg", createLearnROrg).Methods("POST")         //Create a LearnR Org
	myRouter.HandleFunc("/createLearnR", createLearnR).Methods("POST")               //Create a LearnR
	myRouter.HandleFunc("/canLogin", canLogin).Methods("POST")                       //Check User Login
	myRouter.HandleFunc("/createUser", createUser).Methods("POST")                   //Create User
	myRouter.HandleFunc("/canSendLearnR", canSendLearnR).Methods("POST")             //Send LearnR
	myRouter.HandleFunc("/searchLearnRs", searchLearnRs).Methods("POST")             //Send LearnR
	//Used for Learnr functions
	myRouter.HandleFunc("/giveAllLearnrDisplay", giveAllLearnrDisplay).Methods("GET") //Get Learnrs for display
	//Serve our static files
	myRouter.Handle("/", http.FileServer(http.Dir("./static")))
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

//Loads in the initial text API and MongoCrud URLS
func loadInMicroServiceURL() {
	//Check to see if ENV Creds are available first
	_, ok := os.LookupEnv("CRUD_URL")
	if !ok {
		message := "This ENV Variable is not present: " + "CRUD_URL"
		panic(message)
	}
	_, ok2 := os.LookupEnv("TEXT_API")
	if !ok2 {
		message := "This ENV Variable is not present: " + "TEXT_API"
		panic(message)
	}

	mongoCrudURL = os.Getenv("CRUD_URL")
	textAPIURL = os.Getenv("TEXT_API")
}
