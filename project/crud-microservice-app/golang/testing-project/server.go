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
	//Serve our LearnR Crud API
	myRouter.HandleFunc("/addLearnOrg", addLearnOrg).Methods("POST")       //Add a LearnROrg
	myRouter.HandleFunc("/deleteLearnOrg", deleteLearnOrg).Methods("POST") //Delete a LearnROrg
	myRouter.HandleFunc("/updateLearnOrg", updateLearnOrg).Methods("POST") //Update a LearnROrg
	myRouter.HandleFunc("/getLearnOrg", getLearnOrg).Methods("POST")       //Get LearnROrg
	//Serve our validation APIs
	myRouter.HandleFunc("/giveAllUsernames", giveAllUsernames).Methods("GET")       //Get all our Usernames
	myRouter.HandleFunc("/giveAllLearnROrg", giveAllLearnROrg).Methods("GET")       //Get all our LearnROrg Names
	myRouter.HandleFunc("/randomIDCreationAPI", randomIDCreationAPI).Methods("GET") //Get a random ID
	myRouter.HandleFunc("/userLogin", userLogin).Methods("POST")                    //Get a random ID
	//Serve our static files
	log.Fatal(http.ListenAndServe(":4000", myRouter))
}
