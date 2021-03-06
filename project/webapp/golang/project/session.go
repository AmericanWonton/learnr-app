package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const sessionLength int = 400 //Length of sessions

//Here's our session struct
type theSession struct {
	username     string
	lastActivity time.Time
}

//Session Database info
var dbUsers = map[string]User{}          // user ID, user
var dbSessions = map[string]theSession{} // session ID, session
var dbSessionsCleaned time.Time

//Create Login Session ID
func createSessionID(w http.ResponseWriter, req *http.Request, theUser User) {
	dbUsers[theUser.UserName] = theUser
	// create session
	uuidWithHyphen := uuid.New().String()

	cookie := &http.Cookie{
		Name:  "session",
		Value: uuidWithHyphen,
	}
	cookie.MaxAge = sessionLength
	http.SetCookie(w, cookie)
	dbSessions[cookie.Value] = theSession{theUser.UserName, time.Now()}
}

//Gets the User from the current session
func getUser(w http.ResponseWriter, req *http.Request) User {
	// get cookie
	cookie, err := req.Cookie("session")
	//If there is no session cookie, create a new session cookie
	if err != nil {
		uuidWithHyphen := uuid.New().String()
		cookie = &http.Cookie{
			Name:  "session",
			Value: uuidWithHyphen,
		}
	}
	//Set the cookie age to the max length again.
	cookie.MaxAge = sessionLength
	http.SetCookie(w, cookie) //Set the cookie to our grabbed cookie,(or new cookie)

	// if the user exists already, get user
	var theUser User
	if session, ok := dbSessions[cookie.Value]; ok {
		session.lastActivity = time.Now()
		dbSessions[cookie.Value] = session
		theUser = dbUsers[session.username]
	}
	return theUser
}

//Checks to see if this User already has a cookie and is logged in
func alreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("session")
	if err != nil {
		return false //If there is an error getting the cookie, return false
	}
	//if session is found, we update the session with the newest time since activity!
	session, ok := dbSessions[cookie.Value]
	if ok {
		session.lastActivity = time.Now()
		dbSessions[cookie.Value] = session
	}
	/* Check to see if the Username exists from this Session Username. If not, we return false. */
	_, ok = dbUsers[session.username]
	// refresh session
	cookie.MaxAge = sessionLength
	http.SetCookie(w, cookie)
	return ok
}

//Logs User out of session by removing cookie
func logUserOut(w http.ResponseWriter, r *http.Request) {
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}

	//Send a response back to Ajax after session is made
	type SuccessMSG struct {
		Message    string `json:"Message"`
		SuccessNum int    `json:"SuccessNum"`
	}
	theSuccMessage := SuccessMSG{Message: "Successful Session delete", SuccessNum: 0}

	//Marshal the user data into our type
	var userSessionUser User
	json.Unmarshal(bs, &userSessionUser)

	//Remove this User from the map
	delete(dbUsers, userSessionUser.UserName)
	//Remove from Session map
	sessionID := ""
	for key, element := range dbSessions {
		if strings.Contains(userSessionUser.UserName, element.username) {
			sessionID = key
		}
	}
	delete(dbSessions, sessionID)
	//Return JSON
	theJSONMessage, err := json.Marshal(theSuccMessage)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
	}
	fmt.Fprint(w, string(theJSONMessage))
}
