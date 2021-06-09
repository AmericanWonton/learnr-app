package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

/* Mongo DB Creds */

/* Declarative structs for our testing */

//UserCrud Create
type UserCrudCreate struct {
	TheUser             User
	ExpectedNum         int
	ExpectedStringArray []string
}

var userCrudCreateResults []UserCrudCreate

//UserCrud Read
type UserCrudRead struct {
	TheUser             User
	ExpectedNum         int
	ExpectedStringArray []string
}

var userCrudReadResults []UserCrudRead

func TestMain(m *testing.M) {
	//Build stuff for beginning of tests
	log.Println("Starting stuff in TestMain")
	fmt.Println("Starting stuff in TestMain")
	setup()
	code := m.Run()
	//Do stuff for ending of tests
	log.Println("Ending stuff in Test main")
	fmt.Println("Ending stuff in test main")
	shutdown()

	os.Exit(code)
}

//This is used for a default router we can run test http requests on
func Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	//Handle our User CRUD operations
	router.HandleFunc("/addUser", addUser).Methods("POST")
	router.HandleFunc("/deleteUser", deleteUser).Methods("POST")
	router.HandleFunc("/updateUser", updateUser).Methods("POST")
	router.HandleFunc("/getUser", getUser).Methods("POST")
	//Handle our field validation
	router.HandleFunc("/giveAllUsernames", giveAllUsernames).Methods("GET")
	return router
}

//This is setup values declared for testing
func setup() {
	fmt.Printf("Setting up test values...\n")
	/* Start by connecting to Mongo client */
	getCredsMongo()        //Get mongo creds
	createCreateUserCrud() //Add our User Crud testing values for Create
	//createReadUserCrud()   //Add our User Crud testing values for reading
}

//This creates our Crud Testing cases for Creating Users
func createCreateUserCrud() {
	theTimeNow := time.Now() //Used for creating time later
	//Good User Crud Create
	userCrudCreateResults = append(userCrudCreateResults, UserCrudCreate{User{
		UserName:    "TestUsername",
		Password:    hex.EncodeToString([]byte("testpword")),
		Firstname:   "Test",
		Lastname:    "User",
		PhoneNums:   []string{"13143228594"},
		UserID:      1111,
		Email:       []string{"jbkeller0303@gmail.com"},
		Whoare:      "I am a test User and I like to write tests",
		AdminOrgs:   []int{1111},
		OrgMember:   []int{1111},
		Banned:      false,
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}, 0, []string{"User successfully added in addUser"}})
	//Bad User Crud
	userCrudCreateResults = append(userCrudCreateResults, UserCrudCreate{User{}, 1,
		[]string{"Error adding User in addUser", "Error reading the request"}})
}

//This creates our Crud Testing cases for Reading Users
func createReadUserCrud() {
	//Good User Read Crud Create
}

//This is shutdown values/actions for testing
func shutdown() {
	fmt.Printf("Setting up shutdown values/functions...\n")
}

/* Test API Call sections */
//Test User Insert

/* Test directory read */
func TestReadFile(t *testing.T) {
	data, err := ioutil.ReadFile("test-data/test.data")
	if err != nil {
		t.Fatal("Could not open file:\n" + err.Error())
	}
	if string(data) != "hello world from test.data" {
		t.Fatal("String contents do not match expected")
	}
}

/* Test logwrite */
func TestLogWriter(t *testing.T) {
	/* Test read */
	_, err := ioutil.ReadFile("logging/crudapilog.txt")
	if err != nil {
		t.Fatal("Could not open file:\n" + err.Error())
	}
	/* Test logwriter write */
	logWriter("This is a test message")
}

/* Test HTTP Example */
func TestHTTPRequest(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "{ \"status\": \"good\" }")
	}

	r := httptest.NewRequest("GET", "http://josephkeller.me/", nil)
	w := httptest.NewRecorder()
	handler(w, r)

	resp := w.Result()
	body, theErr := ioutil.ReadAll(resp.Body)
	fmt.Printf("Here is our response code: %v\n", string(body))
	if 200 != resp.StatusCode {
		t.Fatal("Status Code not okay: " + theErr.Error())
	}
}

/* TESTING STUFF FOR OUR TESTING FILES! SHOULD BE COMMENTED OUT IF NOT NEEDED! */
func TestGiveUsernames(t *testing.T) {

	request, _ := http.NewRequest("GET", "/giveAllUsernames", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}
