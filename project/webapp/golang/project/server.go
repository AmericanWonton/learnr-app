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
	myRouter.HandleFunc("/", index) //Serve index page
	//Serve our static files
	myRouter.Handle("/", http.FileServer(http.Dir("./static")))
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":3000", myRouter))
}
