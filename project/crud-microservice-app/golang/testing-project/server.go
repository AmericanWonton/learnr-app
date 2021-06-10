package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Handles all requests coming in
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	fmt.Printf("Handling requests on 4000...\n")

	logWriter("We are now handling requests")
	//Serve our User Crud API
	myRouter.HandleFunc("/addUser", addUser).Methods("POST")       //Add a User
	myRouter.HandleFunc("/deleteUser", deleteUser).Methods("POST") //Delete a User
	myRouter.HandleFunc("/updateUser", updateUser).Methods("POST") //Update a User
	myRouter.HandleFunc("/getUser", getUser).Methods("POST")       //Get User
	//Serve our validation APIs
	myRouter.HandleFunc("/giveAllUsernames", giveAllUsernames).Methods("GET") //Get all our Usernames
	//Serve our static files
	log.Fatal(http.ListenAndServe(":4000", myRouter))
}
