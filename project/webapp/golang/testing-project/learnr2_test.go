package main

import (
	"strconv"
	"testing"
	"time"
)

/* LEARNINFO CRUD SECTION */

type LearnRInfoCrudCreate struct {
	LearnrInfo          LearnrInfo
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var LearnRInfoCrudCreateResults []LearnRInfoCrudCreate

//LearnRInfo Crud Read
type LearnRInfoCrudRead struct {
	ID                  int
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var LearnRInfoCrudReadResults []LearnRInfoCrudRead

//LearnRInfo Crud Update
type LearnRInfoCrudUpdate struct {
	TheLearnrInfo       LearnrInfo
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var LearnRInfoCrudUpdateResults []LearnRInfoCrudUpdate

//LearnrInfo CRUD Delete
type LearnrInfoCrudDelete struct {
	ID                  int
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var LearnRInfoCrudDeleteResults []LearnrInfoCrudDelete

func createCreateLearnrInfoCrud() {
	theTimeNow := time.Now() //Used for creating time later
	//Good Crud Create
	LearnRInfoCrudCreateResults = append(LearnRInfoCrudCreateResults, LearnRInfoCrudCreate{LearnrInfo{
		ID:               1111,
		LearnRID:         1234,
		AllSessions:      []LearnRSession{},
		FinishedSessions: []LearnRSession{},
		DateCreated:      theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated:      theTimeNow.Format("2006-01-02 15:04:05"),
	}, 0, true, []string{"Learnr successfully added in addlearnr"}})
	//Empty Crud
	LearnRInfoCrudCreateResults = append(LearnRInfoCrudCreateResults, LearnRInfoCrudCreate{LearnrInfo{}, 1, false,
		[]string{"Error adding LearnRInfo in addLearnrInfo", "Error adding Learnr in addLearnrInfo in crudoperations API"}})
	// with Zero value
	LearnRInfoCrudCreateResults = append(LearnRInfoCrudCreateResults, LearnRInfoCrudCreate{LearnrInfo{ID: 0}, 1, false,
		[]string{"Error adding LearnRInfo in addLearnrInfo", "Error adding Learnr in addLearnrInfo in crudoperations API"}})
	// with negative OrgID value
	LearnRInfoCrudCreateResults = append(LearnRInfoCrudCreateResults, LearnRInfoCrudCreate{LearnrInfo{ID: -1}, 1, false,
		[]string{"Error adding LearnRInfo in addLearnrInfo", "Error adding Learnr in addLearnrInfo in crudoperations API"}})
}

//This creates our CRUD Testing cases for Reading LearnrInfo
func createLearnrInfoReadCrud() {
	//Good Crud Read
	LearnRInfoCrudReadResults = append(LearnRInfoCrudReadResults, LearnRInfoCrudRead{1111, 0, true,
		[]string{"LearnrInfo successfully read in getLearnrInfo"}})
	//Bad CRUD Read
	LearnRInfoCrudReadResults = append(LearnRInfoCrudReadResults, LearnRInfoCrudRead{0, 1, false,
		[]string{"Error adding LearnrInfo in addLearnrInfo", "Error reading the request"}})
	//Not seen ID
	LearnRInfoCrudReadResults = append(LearnRInfoCrudReadResults, LearnRInfoCrudRead{4000000, 1, false,
		[]string{"Error adding LearnrInfo in addLearnrInfo", "Error reading the request"}})
	//Another not seen ID
	LearnRInfoCrudReadResults = append(LearnRInfoCrudReadResults, LearnRInfoCrudRead{-1, 1, false,
		[]string{"Error adding LearnrInfo in addLearnrInfo", "Error reading the request"}})
}

//This creates our CRUD Update cases for Updating LearnrInfo
func createLearnrInfoUpdateCrud() {
	theTimeNow := time.Now() //Used for creating time later
	//Good Crud Create
	LearnRInfoCrudUpdateResults = append(LearnRInfoCrudUpdateResults, LearnRInfoCrudUpdate{LearnrInfo{
		ID:               1111,
		LearnRID:         4444,
		AllSessions:      []LearnRSession{},
		FinishedSessions: []LearnRSession{},
		DateCreated:      theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated:      theTimeNow.Format("2006-01-02 15:04:05"),
	}, 0, true, []string{"LearnRInfo successfully updated in addLearnRInfo"}})
	//Bad Non-Existent ID
	LearnRInfoCrudUpdateResults = append(LearnRInfoCrudUpdateResults, LearnRInfoCrudUpdate{LearnrInfo{
		ID:          400000,
		LearnRID:    4444,
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}, 1, false, []string{"Error updating LearnRInfo in updateLearnRInfo", "Error reading the request"}})
	//Bad Empty LearnOrg Crud
	LearnRInfoCrudUpdateResults = append(LearnRInfoCrudUpdateResults, LearnRInfoCrudUpdate{LearnrInfo{}, 1, false,
		[]string{"Error updating LearnRInfo in updateLearnRInfo", "Error reading the request"}})
}

//This creates our CRUD Delete Cases for deleting LearnrInfo
func createLearnrInfoDeleteCrud() {
	//Good Crud Read
	LearnRInfoCrudDeleteResults = append(LearnRInfoCrudDeleteResults, LearnrInfoCrudDelete{1111, 0, true,
		[]string{"LearnrInfo successfully deleted in deleteLearnrInfo"}})
	//Bad CRUD Read
	LearnRInfoCrudDeleteResults = append(LearnRInfoCrudDeleteResults, LearnrInfoCrudDelete{0, 1, false,
		[]string{"Error deleting LearnrInfo in deleteLearnrInfo", "Error reading the request"}})
	//Not seen ID
	LearnRInfoCrudDeleteResults = append(LearnRInfoCrudDeleteResults, LearnrInfoCrudDelete{4000000, 1, false,
		[]string{"Error deleting LearnrInfo in deleteLearnrInfo", "Error reading the request"}})
	//Another not seen ID
	LearnRInfoCrudDeleteResults = append(LearnRInfoCrudDeleteResults, LearnrInfoCrudDelete{-1, 1, false,
		[]string{"Error deleting LearnrInfo in deleteLearnrInfo", "Error reading the request"}})
}

func TestLearnrInfoAdd(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range LearnRInfoCrudCreateResults {
		success, message := callAddLearnrInfo(test.LearnrInfo)
		if success != test.ExpectedTruth {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message)
		}
		testNum = testNum + 1
	}
}

//Test for updating LearnRInfo
func TestLearnrInfoUpdate(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range LearnRInfoCrudUpdateResults {
		success, message := callUpdateLearnrInfo(test.TheLearnrInfo)
		if success != test.ExpectedTruth {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message)
		}
		testNum = testNum + 1 //Increment this number for testing
	}
}

//Test for Reading LearnRInfo
func TestLearnrInfoRead(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range LearnRInfoCrudReadResults {
		success, message, _ := callReadLearnrInfo(test.ID)
		if success != test.ExpectedTruth {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message + " ")
		}
		testNum = testNum + 1 //Increment this number for testing
	}
}

//Test for Deleting LearnRInfo
func TestLearnrInfoDelete(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range LearnRInfoCrudDeleteResults {
		success, message := callDeleteLearnrInfo(test.ID)
		if success != test.ExpectedTruth {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message)
		}
		testNum = testNum + 1 //Increment this number for testing
	}
}

/* LEARNRSESSION CRUD SECTION */

type LearnRSessionCrudCreate struct {
	LearnRSession       LearnRSession
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var LearnRSessionCrudCreateResults []LearnRSessionCrudCreate

//LearnRSession Crud Read
type LearnRSessionCrudRead struct {
	ID                  int
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var LearnRSessionCrudReadResults []LearnRSessionCrudRead

//LearnRSession Crud Update
type LearnRSessionCrudUpdate struct {
	TheLearnrSession    LearnRSession
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var LearnRSessionCrudUpdateResults []LearnRSessionCrudUpdate

//LearnrSession CRUD Delete
type LearnrSessionCrudDelete struct {
	ID                  int
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var LearnrSessionCrudDeleteResults []LearnrSessionCrudDelete

func createCreateLearnrSessionCrud() {
	theTimeNow := time.Now() //Used for creating time later
	//Good Crud Create
	LearnRSessionCrudCreateResults = append(LearnRSessionCrudCreateResults, LearnRSessionCrudCreate{LearnRSession{
		ID:               1111,
		LearnRID:         1234,
		LearnRName:       "Test LearnRSession",
		TheLearnR:        Learnr{},
		TheUser:          User{},
		TargetUserNumber: "3143695167",
		Ongoing:          true,
		TextsSent:        []LearnRInforms{},
		UserResponses:    []string{"Cool", "Don't care"},
		DateCreated:      theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated:      theTimeNow.Format("2006-01-02 15:04:05"),
	}, 0, true, []string{"LearnrSession successfully added in addlearnrSession"}})
	//Empty Crud
	LearnRSessionCrudCreateResults = append(LearnRSessionCrudCreateResults, LearnRSessionCrudCreate{LearnRSession{}, 1, false,
		[]string{"Error adding LearnRSession in addLearnSession", "Error adding Learnr in addLearnrSesion in crudoperations API"}})
	// with Zero value
	LearnRSessionCrudCreateResults = append(LearnRSessionCrudCreateResults, LearnRSessionCrudCreate{LearnRSession{ID: 0}, 1, false,
		[]string{"Error adding LearnRSession in addLearnSession", "Error adding Learnr in addLearnrSesion in crudoperations API"}})
	// with negative OrgID value
	LearnRSessionCrudCreateResults = append(LearnRSessionCrudCreateResults, LearnRSessionCrudCreate{LearnRSession{ID: -1}, 1, false,
		[]string{"Error adding LearnRSession in addLearnSession", "Error adding Learnr in addLearnrSesion in crudoperations API"}})
}

//This creates our CRUD Testing cases for Reading LearnrSession
func createLearnrSessionReadCrud() {
	//Good Crud Read
	LearnRSessionCrudReadResults = append(LearnRSessionCrudReadResults, LearnRSessionCrudRead{1111, 0, true,
		[]string{"LearnrSession successfully read in getLearnrSession"}})
	//Bad CRUD Read
	LearnRSessionCrudReadResults = append(LearnRSessionCrudReadResults, LearnRSessionCrudRead{0, 1, false,
		[]string{"Error adding LearnrSession in addLearnrSession", "Error reading the request"}})
	//Not seen ID
	LearnRSessionCrudReadResults = append(LearnRSessionCrudReadResults, LearnRSessionCrudRead{4000000, 1, false,
		[]string{"Error adding LearnrSession in addLearnrSession", "Error reading the request"}})
	//Another not seen ID
	LearnRSessionCrudReadResults = append(LearnRSessionCrudReadResults, LearnRSessionCrudRead{-1, 1, false,
		[]string{"Error adding LearnrSession in addLearnrSession", "Error reading the request"}})
}

//This creates our CRUD Update cases for Updating LearnrSession
func createLearnrSessionUpdateCrud() {
	theTimeNow := time.Now() //Used for creating time later
	//Good Crud Create
	LearnRSessionCrudUpdateResults = append(LearnRSessionCrudUpdateResults, LearnRSessionCrudUpdate{LearnRSession{
		ID:               1111,
		LearnRID:         4444,
		LearnRName:       "Name updated",
		TheLearnR:        Learnr{},
		TheUser:          User{},
		TargetUserNumber: "44455666",
		Ongoing:          false,
		TextsSent:        []LearnRInforms{},
		UserResponses:    []string{"boof"},
		DateCreated:      theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated:      theTimeNow.Format("2006-01-02 15:04:05"),
	}, 0, true, []string{"LearnRInfo successfully updated in addLearnRInfo"}})
	//Bad Non-Existent ID
	LearnRSessionCrudUpdateResults = append(LearnRSessionCrudUpdateResults, LearnRSessionCrudUpdate{LearnRSession{
		ID:          400000,
		LearnRID:    4444,
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}, 1, false, []string{"Error updating LearnRSession in updateLearnRSession", "Error reading the request"}})
	//Bad Empty Crud
	LearnRSessionCrudUpdateResults = append(LearnRSessionCrudUpdateResults, LearnRSessionCrudUpdate{LearnRSession{}, 1, false,
		[]string{"Error updating LearnRSession in updateLearnRSession", "Error reading the request"}})
}

//This creates our CRUD Delete Cases for deleting LearnrSession
func createLearnrSessionDeleteCrud() {
	//Good Crud Read
	LearnrSessionCrudDeleteResults = append(LearnrSessionCrudDeleteResults, LearnrSessionCrudDelete{1111, 0, true,
		[]string{"LearnrSession successfully deleted in deleteLearnrSession"}})
	//Bad CRUD Read
	LearnrSessionCrudDeleteResults = append(LearnrSessionCrudDeleteResults, LearnrSessionCrudDelete{0, 1, false,
		[]string{"Error deleting LearnrSession in deleteLearnrSession", "Error reading the request"}})
	//Not seen ID
	LearnrSessionCrudDeleteResults = append(LearnrSessionCrudDeleteResults, LearnrSessionCrudDelete{4000000, 1, false,
		[]string{"Error deleting LearnrSession in deleteLearnrSession", "Error reading the request"}})
	//Another not seen ID
	LearnrSessionCrudDeleteResults = append(LearnrSessionCrudDeleteResults, LearnrSessionCrudDelete{-1, 1, false,
		[]string{"Error deleting LearnrSession in deleteLearnrSession", "Error reading the request"}})
}

func TestLearnrSessionAdd(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range LearnRSessionCrudCreateResults {
		success, message := callAddLearnRSession(test.LearnRSession)
		if success != test.ExpectedTruth {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message)
		}
		testNum = testNum + 1
	}
}

//Test for updating LearnrSession
func TestLearnrSessionUpdate(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range LearnRSessionCrudUpdateResults {
		success, message := callUpdateLearnRSession(test.TheLearnrSession)
		if success != test.ExpectedTruth {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message)
		}
		testNum = testNum + 1 //Increment this number for testing
	}
}

//Test for Reading LearnrSession
func TestLearnrSessionRead(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range LearnRSessionCrudReadResults {
		success, message, _ := callReadLearnRSession(test.ID)
		if success != test.ExpectedTruth {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message + " ")
		}
		testNum = testNum + 1 //Increment this number for testing
	}
}

//Test for Deleting LearnrSession
func TestLearnrSessionDelete(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range LearnrSessionCrudDeleteResults {
		success, message := callDeleteLearnRSession(test.ID)
		if success != test.ExpectedTruth {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message)
		}
		testNum = testNum + 1 //Increment this number for testing
	}
}

/* LEARNRINFORMS CRUD SECTION */

type LearnRInformsCrudCreate struct {
	LearnRInforms       LearnRInforms
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var LearnRInformsCrudCreateResults []LearnRInformsCrudCreate

//LearnRInform Crud Read
type LearnRInformCrudRead struct {
	ID                  int
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var LearnRInformCrudReadResults []LearnRInformCrudRead

//LearnRInform Crud Update
type LearnRInformCrudUpdate struct {
	TheLearnRInforms    LearnRInforms
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var LearnRInformCrudUpdateResults []LearnRInformCrudUpdate

//LearnRInform CRUD Delete
type LearnrInformCrudDelete struct {
	ID                  int
	ExpectedNum         int
	ExpectedTruth       bool
	ExpectedStringArray []string
}

var LearnrInformCrudDeleteResults []LearnrInformCrudDelete

func createCreateLearnRInformCrud() {
	theTimeNow := time.Now() //Used for creating time later
	//Good Crud Create
	LearnRInformsCrudCreateResults = append(LearnRInformsCrudCreateResults, LearnRInformsCrudCreate{LearnRInforms{
		ID:          1111,
		Name:        "Test LearnRInform",
		LearnRID:    1234,
		LearnRName:  "Test LearnRSession",
		Order:       1,
		TheInfo:     "Test information",
		ShouldWait:  false,
		WaitTime:    0,
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}, 0, true, []string{"LearnrInform successfully added in addLearnrInform"}})
	//Empty Crud
	LearnRInformsCrudCreateResults = append(LearnRInformsCrudCreateResults, LearnRInformsCrudCreate{LearnRInforms{}, 1, false,
		[]string{"Error adding LearnRinform in addLearnRinform", "Error adding LearnrInform in addLearnRinform in crudoperations API"}})
	// with Zero value
	LearnRInformsCrudCreateResults = append(LearnRInformsCrudCreateResults, LearnRInformsCrudCreate{LearnRInforms{ID: 0}, 1, false,
		[]string{"Error adding LearnRinform in addLearnRinform", "Error adding LearnrInform in addLearnRinform in crudoperations API"}})
	// with negative OrgID value
	LearnRInformsCrudCreateResults = append(LearnRInformsCrudCreateResults, LearnRInformsCrudCreate{LearnRInforms{ID: -1}, 1, false,
		[]string{"Error adding LearnRinform in addLearnRinform", "Error adding LearnrInform in addLearnRinform in crudoperations API"}})
}

//This creates our CRUD Testing cases for Reading LearnRInform
func createLearnRInformReadCrud() {
	//Good Crud Read
	LearnRInformCrudReadResults = append(LearnRInformCrudReadResults, LearnRInformCrudRead{1111, 0, true,
		[]string{"LearnrInform successfully read in getLearnrInform"}})
	//Bad CRUD Read
	LearnRInformCrudReadResults = append(LearnRInformCrudReadResults, LearnRInformCrudRead{0, 1, false,
		[]string{"Error adding Learnrinform in addLearnrInform", "Error reading the request"}})
	//Not seen ID
	LearnRInformCrudReadResults = append(LearnRInformCrudReadResults, LearnRInformCrudRead{4000000, 1, false,
		[]string{"Error adding Learnrinform in addLearnrInform", "Error reading the request"}})
	//Another not seen ID
	LearnRInformCrudReadResults = append(LearnRInformCrudReadResults, LearnRInformCrudRead{-1, 1, false,
		[]string{"Error adding Learnrinform in addLearnrInform", "Error reading the request"}})
}

//This creates our CRUD Update cases for Updating LearnRInform
func createLearnRInformUpdateCrud() {
	theTimeNow := time.Now() //Used for creating time later
	//Good Crud Create
	LearnRInformCrudUpdateResults = append(LearnRInformCrudUpdateResults, LearnRInformCrudUpdate{LearnRInforms{
		ID:          1111,
		Name:        "Test LearnRInform Edited",
		LearnRID:    1234,
		LearnRName:  "Test LearnRSession Edited",
		Order:       1,
		TheInfo:     "Test information but edited",
		ShouldWait:  true,
		WaitTime:    100,
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}, 0, true, []string{"LearnRInform successfully updated in addLearnRInform"}})
	//Bad Non-Existent ID
	LearnRInformCrudUpdateResults = append(LearnRInformCrudUpdateResults, LearnRInformCrudUpdate{LearnRInforms{
		ID:          400000,
		LearnRID:    4444,
		DateCreated: theTimeNow.Format("2006-01-02 15:04:05"),
		DateUpdated: theTimeNow.Format("2006-01-02 15:04:05"),
	}, 1, false, []string{"Error updating LearnRInform in updateLearnRInform", "Error reading the request"}})
	//Bad Empty Crud
	LearnRInformCrudUpdateResults = append(LearnRInformCrudUpdateResults, LearnRInformCrudUpdate{LearnRInforms{}, 1, false,
		[]string{"Error updating LearnRInform in updateLearnRInform", "Error reading the request"}})
}

//This creates our CRUD Delete Cases for deleting LearnRInform
func createLearnRInformDeleteCrud() {
	//Good Crud Read
	LearnrInformCrudDeleteResults = append(LearnrInformCrudDeleteResults, LearnrInformCrudDelete{1111, 0, true,
		[]string{"LearnrInform successfully deleted in deleteLearnrInform"}})
	//Bad CRUD Read
	LearnrInformCrudDeleteResults = append(LearnrInformCrudDeleteResults, LearnrInformCrudDelete{0, 1, false,
		[]string{"Error deleting LearnrInform in deleteLearnrInform", "Error reading the request"}})
	//Not seen ID
	LearnrInformCrudDeleteResults = append(LearnrInformCrudDeleteResults, LearnrInformCrudDelete{4000000, 1, false,
		[]string{"Error deleting LearnrInform in deleteLearnrInform", "Error reading the request"}})
	//Another not seen ID
	LearnrInformCrudDeleteResults = append(LearnrInformCrudDeleteResults, LearnrInformCrudDelete{-1, 1, false,
		[]string{"Error deleting LearnrInform in deleteLearnrInform", "Error reading the request"}})
}

func TestLearnRInformAdd(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range LearnRInformsCrudCreateResults {
		success, message := callAddLearnRInform(test.LearnRInforms)
		if success != test.ExpectedTruth {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message)
		}
		testNum = testNum + 1
	}
}

//Test for updating LearnRInform
func TestLearnRInformUpdate(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range LearnRInformCrudUpdateResults {
		success, message := callUpdateLearnRInform(test.TheLearnRInforms)
		if success != test.ExpectedTruth {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message)
		}
		testNum = testNum + 1 //Increment this number for testing
	}
}

//Test for Reading LearnRInform
func TestLearnRInformRead(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range LearnRInformCrudReadResults {
		success, message, _ := callReadLearnRInform(test.ID)
		if success != test.ExpectedTruth {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message + " ")
		}
		testNum = testNum + 1 //Increment this number for testing
	}
}

//Test for Deleting LearnRInform
func TestLearnRInformDelete(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range LearnrInformCrudDeleteResults {
		success, message := callDeleteLearnRInform(test.ID)
		if success != test.ExpectedTruth {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message)
		}
		testNum = testNum + 1 //Increment this number for testing
	}
}
