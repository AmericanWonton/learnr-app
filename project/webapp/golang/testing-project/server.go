package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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
	//Used for session work
	myRouter.HandleFunc("/logUserOut", logUserOut).Methods("POST") //Remove our cookie after logging out user
	//Serve our Validation API
	myRouter.HandleFunc("/checkUsername", checkUsername).Methods("POST")             //Check Username
	myRouter.HandleFunc("/checkLearnROrgNames", checkLearnROrgNames).Methods("POST") //Check LearnROrg Name
	myRouter.HandleFunc("/canLogin", canLogin).Methods("POST")                       //Check User Login
	myRouter.HandleFunc("/createUser", createUser).Methods("POST")                   //Create User
	//Serve our static files
	myRouter.Handle("/", http.FileServer(http.Dir("./static")))
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
