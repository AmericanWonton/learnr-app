package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Handles all requests coming in
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	logWriter("We are now handling requests")
	//Serve our User Crud API
	myRouter.HandleFunc("/addUser", addUser).Methods("POST") //Check Username
	//Serve our static files
	log.Fatal(http.ListenAndServe(":4000", myRouter))
}
