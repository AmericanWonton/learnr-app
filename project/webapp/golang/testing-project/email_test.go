package main

import (
	"strconv"
	"testing"
	"time"
)

//EmailCrud Create
type EmailVerifCrudCreate struct {
	TheEmailVerif       EmailVerify
	ExpectedNum         int
	ExpectedBool        bool
	ExpectedStringArray []string
}

var emailVerifCrudCreateResults []EmailVerifCrudCreate

//EmailVerifCrud Read
type EmailVerifCrudRead struct {
	ID                  int
	ExpectedNum         int
	ExpectedBool        bool
	ExpectedStringArray []string
}

var emailVerifCrudReadResults []EmailVerifCrudRead

//EmailVerif Crud Update
type EmailVerifCrudUpdate struct {
	TheEmailVerify      EmailVerify
	ExpectedNum         int
	ExpectedBool        bool
	ExpectedStringArray []string
}

var emailVerifCrudUpdateResults []EmailVerifCrudUpdate

//EmailVerif CRUD Delete
type EmailVerifCrudDelete struct {
	ID                  int
	ExpectedNum         int
	ExpectedBool        bool
	ExpectedStringArray []string
}

var emailVerifCrudDeleteResults []EmailVerifCrudDelete

//This creates our Crud Testing cases for Creating Users
func createCreateEmailVerifCrud() {
	//Good Value Crud Create
	emailVerifCrudCreateResults = append(emailVerifCrudCreateResults, EmailVerifCrudCreate{EmailVerify{
		Username: "bigSkeezy",
		Email:    "johnTron95@gmail.com",
		ID:       654321,
		TimeMade: time.Now(),
		Active:   true,
	}, 0, true, []string{"EmailVerif successfully added in addEmailVerif"}})
	//Empty Value Crud
	emailVerifCrudCreateResults = append(emailVerifCrudCreateResults, EmailVerifCrudCreate{EmailVerify{}, 1,
		false, []string{"Error adding EmailVerif in addEmailVerif", "Error reading the request"}})
	//Value with Zero value
	emailVerifCrudCreateResults = append(emailVerifCrudCreateResults, EmailVerifCrudCreate{EmailVerify{ID: 0}, 1,
		false, []string{"Error adding EmailVerify in addEmailVerify", "Error reading the request"}})
	//Value with negative ID value
	emailVerifCrudCreateResults = append(emailVerifCrudCreateResults, EmailVerifCrudCreate{EmailVerify{ID: -1}, 1,
		false, []string{"Error adding EmailVerify in addEmailVerify", "Error reading the request"}})
}

//This creates our CRUD Testing cases for Reading Users
func createEmailVerifReadCrud() {
	//Good Value Crud Read
	emailVerifCrudReadResults = append(emailVerifCrudReadResults, EmailVerifCrudRead{654321,
		0, true, []string{"Email Verificaiton successfully read in getEmail Verificaiton"}})
	//Bad Value CRUD Read
	emailVerifCrudReadResults = append(emailVerifCrudReadResults, EmailVerifCrudRead{0, 1,
		false, []string{"Error reading Email in getEmail Verificaiton", "Error reading the request"}})
	//Not seen ID
	emailVerifCrudReadResults = append(emailVerifCrudReadResults, EmailVerifCrudRead{4000000, 1,
		false, []string{"Error reading Email in getEmail Verificaiton", "Error reading the request"}})
	//Another not seen ID
	emailVerifCrudReadResults = append(emailVerifCrudReadResults, EmailVerifCrudRead{-1, 1,
		false, []string{"Error reading Email in getEmail Verificaiton", "Error reading the request"}})
}

//This creates our CRUD Update cases for Updating Users
func createEmailVerifUpdateCrud() {
	//Good Value Crud Create
	emailVerifCrudUpdateResults = append(emailVerifCrudUpdateResults, EmailVerifCrudUpdate{EmailVerify{
		Username: "bigSkeezy2",
		Email:    "johnTron95@gmail.com",
		ID:       654321,
		TimeMade: time.Now(),
		Active:   true,
	}, 0, true, []string{"Email verify successfully added in updateEmailVerify"}})
	//Bad Non-Existent ID
	emailVerifCrudUpdateResults = append(emailVerifCrudUpdateResults, EmailVerifCrudUpdate{EmailVerify{
		Username: "bigSkeezy",
		Email:    "johnTron95@gmail.com",
		ID:       5556667778888,
		TimeMade: time.Now(),
		Active:   true,
	}, 1, false, []string{"Error updating Email Verify", "Error reading the request"}})
	//Bad Empty Value Crud
	emailVerifCrudUpdateResults = append(emailVerifCrudUpdateResults, EmailVerifCrudUpdate{EmailVerify{}, 1,
		false, []string{"Error updating Email Verify", "Error reading the request"}})
}

//This creates our CRUD Delete Cases for deleting Users
func createEmailVerifDeleteCrud() {
	//Good value Crud Read
	emailVerifCrudDeleteResults = append(emailVerifCrudDeleteResults, EmailVerifCrudDelete{654321, 0,
		true, []string{"Email Verification successfully deleted in deleteEmailVerification"}})
	//Bad value CRUD Read
	emailVerifCrudDeleteResults = append(emailVerifCrudDeleteResults, EmailVerifCrudDelete{0, 1,
		false, []string{"Error deleting Email Verificaiton in deleteEmailVerification", "Error reading the request"}})
	//Not seen ID
	emailVerifCrudDeleteResults = append(emailVerifCrudDeleteResults, EmailVerifCrudDelete{4000000, 1,
		false, []string{"Error deleting Email Verificaiton in deleteEmailVerification", "Error reading the request"}})
	//Another not seen ID
	emailVerifCrudDeleteResults = append(emailVerifCrudDeleteResults, EmailVerifCrudDelete{-1, 1,
		false, []string{"Error deleting Email Verificaiton in deleteEmailVerification", "Error reading the request"}})
}

/* Testing for Email Verification */
func TestEmailVerificationAdd(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range emailVerifCrudCreateResults {
		success, message := addEmailVerif(test.TheEmailVerif)
		if success != test.ExpectedBool {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message)
		}
		testNum = testNum + 1
	}
}

//Test for updating EmailVerification
func TestEmailVerificationUpdate(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range emailVerifCrudUpdateResults {
		success, message := updateEmailVerify(test.TheEmailVerify)
		if success != test.ExpectedBool {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message)
		}
		testNum = testNum + 1 //Increment this number for testing
	}
}

//Test for Reading EmailVerification
func TestEmailVerificationRead(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range emailVerifCrudReadResults {
		success, message, _ := getEmailVerify(test.ID)
		if success != test.ExpectedBool {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message + " ")
		}
		testNum = testNum + 1 //Increment this number for testing
	}
}

//Test for Deleting EmailVerification
func TestEmailVerificationDelete(t *testing.T) {
	testNum := 0 //Used for incrementing
	for _, test := range emailVerifCrudDeleteResults {
		success, message := deleteEmailVerif(test.ID)
		if success != test.ExpectedBool {
			t.Fatal("Failed at this step: " + strconv.Itoa(testNum) + " :" + message)
		}
		testNum = testNum + 1 //Increment this number for testing
	}
}
