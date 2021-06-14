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
	//Serve our LearnROrg Crud API
	myRouter.HandleFunc("/addLearnOrg", addLearnOrg).Methods("POST")       //Add a LearnROrg
	myRouter.HandleFunc("/deleteLearnOrg", deleteLearnOrg).Methods("POST") //Delete a LearnROrg
	myRouter.HandleFunc("/updateLearnOrg", updateLearnOrg).Methods("POST") //Update a LearnROrg
	myRouter.HandleFunc("/getLearnOrg", getLearnOrg).Methods("POST")       //Get LearnROrg
	//Serve our LearnR Crud API
	myRouter.HandleFunc("/addLearnR", addLearnR).Methods("POST")       //Add a LearnR
	myRouter.HandleFunc("/deleteLearnR", deleteLearnR).Methods("POST") //Delete a LearnR
	myRouter.HandleFunc("/updateLearnR", updateLearnR).Methods("POST") //Update a LearnR
	myRouter.HandleFunc("/getLearnR", getLearnR).Methods("POST")       //Get LearnR
	//Serve our LearnRInfo Crud API
	myRouter.HandleFunc("/addLearnrInfo", addLearnrInfo).Methods("POST")       //Add a LearnRInfo
	myRouter.HandleFunc("/deleteLearnrInfo", deleteLearnrInfo).Methods("POST") //Delete a LearnRInfo
	myRouter.HandleFunc("/updateLearnrInfo", updateLearnrInfo).Methods("POST") //Update a LearnRInfo
	myRouter.HandleFunc("/getLearnrInfo", getLearnrInfo).Methods("POST")       //Get LearnRInfo
	//Serve our LearnRSessions CRUD API
	myRouter.HandleFunc("/addLearnRSession", addLearnRSession).Methods("POST")       //Add a LearnRSession
	myRouter.HandleFunc("/deleteLearnRSession", deleteLearnRSession).Methods("POST") //Delete a LearnRSession
	myRouter.HandleFunc("/updateLearnRSession", updateLearnRSession).Methods("POST") //Update a LearnRSession
	myRouter.HandleFunc("/getLearnRSession", getLearnRSession).Methods("POST")       //Get LearnRSession
	//Serve our LearnRInforms CRUD API
	myRouter.HandleFunc("/addLearnRInforms", addLearnRInforms).Methods("POST")       //Add a LearnRInfo
	myRouter.HandleFunc("/deleteLearnRInforms", deleteLearnRInforms).Methods("POST") //Delete a LearnRInfo
	myRouter.HandleFunc("/updateLearnRInforms", updateLearnRInforms).Methods("POST") //Update a LearnRInfo
	myRouter.HandleFunc("/getLearnRInforms", getLearnRInforms).Methods("POST")       //Get LearnRInfo
	//Serve our validation APIs
	myRouter.HandleFunc("/giveAllUsernames", giveAllUsernames).Methods("GET")       //Get all our Usernames
	myRouter.HandleFunc("/giveAllLearnROrg", giveAllLearnROrg).Methods("GET")       //Get all our LearnROrg Names
	myRouter.HandleFunc("/randomIDCreationAPI", randomIDCreationAPI).Methods("GET") //Get a random ID
	myRouter.HandleFunc("/userLogin", userLogin).Methods("POST")                    //Get a random ID
	//Serve our static files
	log.Fatal(http.ListenAndServe(":4000", myRouter))
}
