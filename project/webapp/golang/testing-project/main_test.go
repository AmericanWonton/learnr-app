package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

//Declare struct we are expecting
type TestSearchR struct {
	TheCases   []int  `json:"TheCases"`
	TheTag     string `json:"TheTag"`
	LearnRName string `json:"LearnRName"`
	EntryFrom  int    `json:"EntryFrom"`
	EntryTo    int    `json:"EntryTo"`
}

type LearnRSearchCrudCreate struct {
	LearnRSearches      TestSearchR
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
	ExpectedLearnRID    map[string]int
}

var LearnRSearchCrudCreators []LearnRSearchCrudCreate

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

//This is setup values declared for testing
func setup() {
	fmt.Printf("Setting up test values...\n")
	/* Start by connecting to Mongo client */
	getCredsMongo()        //Get mongo creds
	createCreateUserCrud() //Add our User Crud testing values for Create
	createUserReadCrud()   //Add our User Crud testing values for Reading
	createUserUpdateCrud() // Add our User Crud testing values for updating
	createUserDeleteCrud() //Add our User Crud testing values for deleting
	createUserLogin()      //Create creds for logging Users in
	/* Add values for LearnR Org test cases */
	createCreateLearnOrgCrud()
	createLearnOrgReadCrud()
	createLearnOrgUpdateCrud()
	createLearnOrgDeleteCrud()
	/* Add values for LearnR test cases */
	createCreateLearnrCrud()
	createLearnrReadCrud()
	createLearnrUpdateCrud()
	createLearnrDeleteCrud()
	createLearnRSpecialGet()
	/* Add values for LearnRInfo test cases */
	createCreateLearnrInfoCrud()
	createLearnrInfoReadCrud()
	createLearnrInfoUpdateCrud()
	createLearnrInfoDeleteCrud()
	/* Add values for LearnrSession test cases */
	createCreateLearnrSessionCrud()
	createLearnrSessionReadCrud()
	createLearnrSessionUpdateCrud()
	createLearnrSessionDeleteCrud()
	/* Add values for LearnRInform test cases */
	createCreateLearnRInformCrud()
	createLearnRInformReadCrud()
	createLearnRInformUpdateCrud()
	createLearnRInformDeleteCrud()
	/* Add values for special LearnR Search */
	createSpecialLearnRSearch()
}

//This is shutdown values/actions for testing
func shutdown() {
	fmt.Printf("Setting up shutdown values/functions...\n")
}

/* Test DIRECTORY EXAMPLE */
func TestReadFile(t *testing.T) {
	data, err := ioutil.ReadFile("test-data/test.data")
	if err != nil {
		t.Fatal("Could not open file:\n" + err.Error())
	}
	if string(data) != "hello world from test.data" {
		t.Fatal("String contents do not match expected")
	}
}

/* Test logwrite example */
func TestLogWriter(t *testing.T) {
	/* Test read */
	_, err := ioutil.ReadFile("logging/weblog.txt")
	if err != nil {
		t.Fatal("Could not open file:\n" + err.Error())
	}
	/* Test logwriter write */
	logWriter("This is a test message")
}

/* Test init example */
func Testinit(t *testing.T) {

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

/* Test handle routes */

//Get a random ID we can use for any struct
func TestRandomID(t *testing.T) {
	//Call our crudOperations Microservice in order to get our Usernames
	req, err := http.Get(GETRANDOMID)
	if err != nil {
		theErr := "There was an error getting Usernames in loadUsernames: " + err.Error()
		t.Fatal(theErr)
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "There was an error getting a response for Usernames in loadUsernames: " + err.Error()
		t.Fatal(theErr)
	}

	//Marshal the response into a type we can read
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
		RandomID   int      `json:"RandomID"`
	}
	var returnedMessage ReturnMessage
	json.Unmarshal(body, &returnedMessage)

	defer req.Body.Close()

	//Assign our map variable to the map varialbe and see if it's okay
	if returnedMessage.SuccOrFail != 0 {
		errString := ""
		for l := 0; l < len(returnedMessage.TheErr); l++ {
			errString = errString + returnedMessage.TheErr[l]
		}
		t.Fatal("Had an error getting map: " + errString)
	} else {
		//fmt.Printf("Here is our random ID: %v\n", returnedMessage.RandomID)
	}
}

func createSpecialLearnRSearch() {
	var theMap map[string]int //Declare map for initialization later
	//Empty Search (GOOD)
	/* Warning, will need to adjust this based on ALL the learnrs in our DB */
	theMap = map[string]int{"855367233056": 855367233056,
		"478483273602": 478483273602,
		"65286261652":  65286261652,
		"645741771884": 645741771884,
		"340870637511": 340870637511,
		"814584281060": 814584281060}
	LearnRSearchCrudCreators = append(LearnRSearchCrudCreators, LearnRSearchCrudCreate{
		LearnRSearches: TestSearchR{TheCases: []int{0, 1, 1, 1},
			TheTag:     "",
			LearnRName: "",
			EntryFrom:  0,
			EntryTo:    0},
		ExpectedNum:         0,
		ExpectedTruth:       true,
		ExpectedStringArray: []string{"Nice"},
		ExpectedLearnRID:    theMap,
	})
	theMap = map[string]int{"855367233056": 855367233056}
	//Single Tag Search
	LearnRSearchCrudCreators = append(LearnRSearchCrudCreators, LearnRSearchCrudCreate{
		LearnRSearches: TestSearchR{TheCases: []int{0, 1, 0, 1},
			TheTag:     "Twitter",
			LearnRName: "",
			EntryFrom:  0,
			EntryTo:    0},
		ExpectedNum:         0,
		ExpectedTruth:       true,
		ExpectedStringArray: []string{"Nice"},
		ExpectedLearnRID:    theMap,
	})
	theMap = map[string]int{"855367233056": 855367233056, "645741771884": 645741771884}
	//Multiple Tag Search
	LearnRSearchCrudCreators = append(LearnRSearchCrudCreators, LearnRSearchCrudCreate{
		LearnRSearches: TestSearchR{TheCases: []int{0, 1, 0, 1},
			TheTag:     "Blue",
			LearnRName: "",
			EntryFrom:  0,
			EntryTo:    0},
		ExpectedNum:         0,
		ExpectedTruth:       true,
		ExpectedStringArray: []string{"Nice"},
		ExpectedLearnRID:    theMap,
	})
	//Multiple Name Search
	theMap = map[string]int{"65286261652": 65286261652, "645741771884": 645741771884}
	LearnRSearchCrudCreators = append(LearnRSearchCrudCreators, LearnRSearchCrudCreate{
		LearnRSearches: TestSearchR{TheCases: []int{0, 0, 1, 1},
			TheTag:     "",
			LearnRName: "the",
			EntryFrom:  0,
			EntryTo:    0},
		ExpectedNum:         0,
		ExpectedTruth:       true,
		ExpectedStringArray: []string{"Nice"},
		ExpectedLearnRID:    theMap,
	})
	//Single Name Search
	theMap = map[string]int{"855367233056": 855367233056}
	LearnRSearchCrudCreators = append(LearnRSearchCrudCreators, LearnRSearchCrudCreate{
		LearnRSearches: TestSearchR{TheCases: []int{0, 0, 1, 1},
			TheTag:     "",
			LearnRName: "Twitter",
			EntryFrom:  0,
			EntryTo:    0},
		ExpectedNum:         0,
		ExpectedTruth:       true,
		ExpectedStringArray: []string{"Nice"},
		ExpectedLearnRID:    theMap,
	})
	//Tag and Name Search
	theMap = map[string]int{"65286261652": 65286261652, "645741771884": 645741771884}
	LearnRSearchCrudCreators = append(LearnRSearchCrudCreators, LearnRSearchCrudCreate{
		LearnRSearches: TestSearchR{TheCases: []int{0, 0, 0, 1},
			TheTag:     "special",
			LearnRName: "the",
			EntryFrom:  0,
			EntryTo:    0},
		ExpectedNum:         0,
		ExpectedTruth:       true,
		ExpectedStringArray: []string{"Nice"},
		ExpectedLearnRID:    theMap,
	})
	//Return Nothing Search
	theMap = map[string]int{}
	LearnRSearchCrudCreators = append(LearnRSearchCrudCreators, LearnRSearchCrudCreate{
		LearnRSearches: TestSearchR{TheCases: []int{0, 0, 0, 1},
			TheTag:     "Insane-test-value",
			LearnRName: "Another-Insane-Test-Value",
			EntryFrom:  0,
			EntryTo:    0},
		ExpectedNum:         0,
		ExpectedTruth:       true,
		ExpectedStringArray: []string{"Nothing"},
		ExpectedLearnRID:    theMap,
	})
}

/* Test Special LearnRSearch */
func TestSpecialLearnRSearch(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range LearnRSearchCrudCreators {
		//Go to special LearRs
		arrayOLearnRs, goodGet, message := getSpecialLearnRs(test.LearnRSearches.TheCases,
			test.LearnRSearches.TheTag, test.LearnRSearches.LearnRName,
			test.LearnRSearches.EntryFrom, test.LearnRSearches.EntryTo)
		if !goodGet {
			t.Fatal("Error from getSpecialLearnRs: " + message)
		}
		//Check to see if returned Array matches test case size
		if len(arrayOLearnRs) != len(test.ExpectedLearnRID) {
			t.Fatal("Test case " + strconv.Itoa(testNum) + " failed, " + "Our array of LearnRs returned does not match our test case: " + strconv.Itoa(len(arrayOLearnRs)) +
				"TestNum: " + strconv.Itoa(testNum) + "\n The test case was: " + test.LearnRSearches.LearnRName + "   " +
				test.LearnRSearches.TheTag)
		}
		for n := 0; n < len(arrayOLearnRs); n++ {
			if _, ok := test.ExpectedLearnRID[strconv.Itoa(arrayOLearnRs[n].ID)]; ok {
				//We have the ID we need
			} else {
				t.Fatal("The LearnRID is not found in our returned LearnRs: " + strconv.Itoa(arrayOLearnRs[n].ID) +
					strconv.Itoa(testNum))
			}
		}
		testNum = testNum + 1
	}
}
