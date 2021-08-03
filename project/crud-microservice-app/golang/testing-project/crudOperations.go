package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

//Mongo DB Declarations
var mongoClient *mongo.Client

var theContext context.Context
var mongoURI string //Connection string loaded

/* App/Data type declarations for our application */
// Desc: This person uses our app
type User struct {
	UserName    string   `json:"UserName"`
	Password    string   `json:"Password"`
	Firstname   string   `json:"Firstname"`
	Lastname    string   `json:"Lastname"`
	PhoneNums   []string `json:"PhoneNums"`
	UserID      int      `json:"UserID"`
	Email       []string `json:"Email"`
	Whoare      string   `json:"Whoare"`
	AdminOrgs   []int    `json:"AdminOrgs"`
	OrgMember   []int    `json:"OrgMember"`
	Banned      bool     `json:"Banned"`
	DateCreated string   `json:"DateCreated"`
	DateUpdated string   `json:"DateUpdated"`
}

//This is used for email verification
type EmailVerify struct {
	Username string    `json:"Username"`
	Email    string    `json:"Email"`
	ID       int       `json:"ID"`
	TimeMade time.Time `json:"TimeMade"`
	Active   bool      `json:"Active"`
}

//This gets the client to connect to our DB
func connectDB() *mongo.Client {
	//Setup Mongo connection to Atlas Cluster
	theClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("Errored getting mongo client: %v\n", err)
		log.Fatal(err)
	}
	theContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = theClient.Connect(theContext)
	if err != nil {
		fmt.Printf("Errored getting mongo client context: %v\n", err)
		log.Fatal(err)
	}
	//Double check to see if we've connected to the database
	err = theClient.Ping(theContext, readpref.Primary())
	if err != nil {
		fmt.Printf("Errored pinging MongoDB: %v\n", err)
		log.Fatal(err)
	}

	return theClient
}

/* USER CRUD OPERATIONS BEGINNING */
//This adds a User to our database; called from anywhere
func addUser(w http.ResponseWriter, req *http.Request) {
	canCrud := true //Used to determine if we're good to try our crud operation

	//Declare data to return
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	theReturnMessage := ReturnMessage{}

	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from addUser: " + err.Error() + "\n" + string(bs)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false //Reading failed, need to return failure
	}
	//Marshal it into our type
	var postedUser User
	json.Unmarshal(bs, &postedUser)

	//Check to see if we can perform CRUD operations and we aren't passing a null User
	if canCrud && postedUser.UserID > 0 {
		user_collection := mongoClient.Database("learnR").Collection("users") //Here's our collection
		collectedUsers := []interface{}{postedUser}
		//Insert Our Data
		_, err2 := user_collection.InsertMany(theContext, collectedUsers)

		if err2 != nil {
			theErr := "Error adding User in addUser in crudoperations API: " + err2.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			theErr := "User successfully added in addUser in crudoperations: " + string(bs)
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 0
		}
	} else {
		theErr := "Error adding User; could not perform CRUD or UserID was bad: " + strconv.Itoa(postedUser.UserID)
		logWriter(theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.SuccOrFail = 1
	}

	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in addUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//This deletes a User to our database; called from anywhere
func deleteUser(w http.ResponseWriter, req *http.Request) {
	canCrud := true //Used to determine if we're good to try our crud operation

	//Declare data to return
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	theReturnMessage := ReturnMessage{}

	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from deleteUser: " + err.Error() + "\n" + string(bs)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}
	//Declare JSON we're looking for
	type UserDelete struct {
		UserID int `json:"UserID"`
	}
	//Marshal it into our type
	var postedUserID UserDelete
	json.Unmarshal(bs, &postedUserID)

	//Delete only if we had no issues above
	if canCrud && postedUserID.UserID > 0 {
		//Search for User and delete
		userCollection := mongoClient.Database("learnR").Collection("users") //Here's our collection
		deletes := []bson.M{
			{"userid": postedUserID.UserID},
		} //Here's our filter to look for
		deletes = append(deletes, bson.M{"userid": bson.M{
			"$eq": postedUserID.UserID,
		}}, bson.M{"userid": bson.M{
			"$eq": postedUserID.UserID,
		}},
		)

		// create the slice of write models
		var writes []mongo.WriteModel

		for _, del := range deletes {
			model := mongo.NewDeleteManyModel().SetFilter(del)
			writes = append(writes, model)
		}

		// run bulk write
		bulkWrite, err := userCollection.BulkWrite(theContext, writes)
		if err != nil {
			theErr := "Error writing delete User in deleteUser in crudoperations: " + err.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			//Check to see if delete count worked; must have deleted at least one
			resultInt := bulkWrite.DeletedCount
			if resultInt > 0 {
				theErr := "User successfully deleted in deleteUser in crudoperations: " + string(bs)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 0
			} else {
				theErr := "No documents deleted for this given UserID: " + strconv.Itoa(postedUserID.UserID)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			}
		}
	} else {
		theErr := "Error, could not CRUD operate in deleteUser, or the number we recieved was wrong: " + strconv.Itoa(postedUserID.UserID)
		logWriter(theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.SuccOrFail = 1
	}

	//Write the response back
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in deleteUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//This updates a User to our database; called from anywhere
func updateUser(w http.ResponseWriter, req *http.Request) {
	canCrud := true
	//Declare data to return
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	theReturnMessage := ReturnMessage{}

	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from updateUser: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Marshal it into our type
	var theUserUpdate User
	json.Unmarshal(bs, &theUserUpdate)

	//Update User if we have successfully decoded from JSON
	if canCrud {
		//Update User if their User ID != 0 or nil
		if theUserUpdate.UserID != 0 {
			//Update User
			theTimeNow := time.Now()
			userCollection := mongoClient.Database("learnR").Collection("users") //Here's our collection
			theFilter := bson.M{
				"userid": bson.M{
					"$eq": theUserUpdate.UserID, // check if bool field has value of 'false'
				},
			}
			updatedDocument := bson.M{
				"$set": bson.M{
					"username":    theUserUpdate.UserName,
					"password":    theUserUpdate.Password,
					"firstname":   theUserUpdate.Firstname,
					"lastname":    theUserUpdate.Lastname,
					"phonenums":   theUserUpdate.PhoneNums,
					"userid":      theUserUpdate.UserID,
					"email":       theUserUpdate.Email,
					"whoare":      theUserUpdate.Whoare,
					"adminorgs":   theUserUpdate.AdminOrgs,
					"orgmember":   theUserUpdate.OrgMember,
					"banned":      theUserUpdate.Banned,
					"datecreated": theUserUpdate.DateCreated,
					"dateupdated": theTimeNow.Format("2006-01-02 15:04:05"),
				},
			}
			updateResult, err := userCollection.UpdateOne(theContext, theFilter, updatedDocument)

			if err != nil {
				theErr := "Error writing update User in updateUser in crudoperations: " + err.Error()
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			} else {
				//Check to see if anything was updated; if not, return the error
				if updateResult.ModifiedCount < 1 {
					theErr := "No document updated with this ID: " + strconv.Itoa(theUserUpdate.UserID)
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 1
				} else {
					theErr := "User successfully updated in updateUser in crudoperations: " + string(bs) + "\n"
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 0
				}
			}
		} else {
			theErr := "The User ID was not found: " + strconv.Itoa(theUserUpdate.UserID)
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		}
	}

	//Send the response back
	theJSONMessage, err := json.Marshal(theReturnMessage)
	if err != nil {
		errIs := "Error formatting JSON for return in updateUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//This gets a User with a certain UserID
func getUser(w http.ResponseWriter, req *http.Request) {
	canCrud := true
	//Declare data to return
	type ReturnMessage struct {
		TheErr       []string `json:"TheErr"`
		ResultMsg    []string `json:"ResultMsg"`
		SuccOrFail   int      `json:"SuccOrFail"`
		ReturnedUser User     `json:"ReturnedUser"`
	}
	theReturnMessage := ReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success

	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from updateUser: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Decalre JSON we recieve
	type UserIDUser struct {
		TheUserID int `json:"TheUserID"`
	}

	//Marshal it into our type
	var theUserGet UserIDUser
	json.Unmarshal(bs, &theUserGet)

	//If we successfully decoded, (and the UesrID is not 0) we can get our user
	if canCrud && theUserGet.TheUserID > 0 {
		/* Find the User with the given Username */
		var theUserReturned User                                             //Initialize User to be returned after Mongo query
		userCollection := mongoClient.Database("learnR").Collection("users") //Here's our collection
		theFilter := bson.M{
			"userid": bson.M{
				"$eq": theUserGet.TheUserID, // check if bool field has value of 'false'
			},
		}
		findOptions := options.Find()
		findUser, err := userCollection.Find(theContext, theFilter, findOptions)
		theFind := 0 //A counter to track how many users we find
		if findUser.Err() != nil || err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				stringUserID := strconv.Itoa(theUserGet.TheUserID)
				returnedErr := "For " + stringUserID + ", no User was returned: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedUser = User{}
			} else {
				stringUserID := strconv.Itoa(theUserGet.TheUserID)
				returnedErr := "For " + stringUserID + ", there was a Mongo Error: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedUser = User{}
			}
		} else {
			//Found User, decode to return
			for findUser.Next(theContext) {
				stringUserID := strconv.Itoa(theUserGet.TheUserID)
				err := findUser.Decode(&theUserReturned)
				if err != nil {
					returnedErr := "For " + stringUserID +
						", there was an error decoding document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedUser = User{}
				} else if len(theUserReturned.UserName) <= 1 {
					returnedErr := "For " + stringUserID +
						", there was an no document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedUser = User{}
				} else {
					//Successful decode, do nothing
				}
				theFind = theFind + 1
			}
			findUser.Close(theContext)
		}

		if theFind <= 0 {
			//Error, return an error back and log it
			stringUserID := strconv.Itoa(theUserGet.TheUserID)
			returnedErr := "For " + stringUserID +
				", No User was returned."
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 1
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
			theReturnMessage.ReturnedUser = User{}
		} else {
			//Success, log the success and return User
			stringUserID := strconv.Itoa(theUserGet.TheUserID)
			returnedErr := "For " + stringUserID +
				", User should be successfully decoded."
			//fmt.Println(returnedErr)
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 0
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ReturnedUser = theUserReturned
		}
	} else {
		//Error, return an error back and log it
		stringUserID := strconv.Itoa(theUserGet.TheUserID)
		returnedErr := "For " + stringUserID +
			", No User was returned. UserID was also not accepted: " + stringUserID
		logWriter(returnedErr)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
		theReturnMessage.ReturnedUser = User{}
	}

	//Format the JSON map for returning our results
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in getUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* USER CRUD OPERATIONS ENDING */

/* LEARNR CRUD OPERATIONS BEGINNING */

//This adds a learnOrg to our DB; called from anywhere
func addLearnOrg(w http.ResponseWriter, req *http.Request) {
	canCrud := true //Used to determine if we're good to try our crud operation

	//Declare data to return
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	theReturnMessage := ReturnMessage{}

	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from learnOrg: " + err.Error() + "\n" + string(bs)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false //Reading failed, need to return failure
	}
	//Marshal it into our type
	var postedLearnorg LearnrOrg
	json.Unmarshal(bs, &postedLearnorg)

	//Check to see if we can perform CRUD operations and we aren't passing a null LearnOrg
	if canCrud && postedLearnorg.OrgID > 0 {
		org_collection := mongoClient.Database("learnR").Collection("learnorg") //Here's our collection
		collectedLearnOrg := []interface{}{postedLearnorg}
		//Insert Our Data
		_, err2 := org_collection.InsertMany(theContext, collectedLearnOrg)

		if err2 != nil {
			theErr := "Error adding LarnOrg in addLarnOrg in crudoperations API: " + err2.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			theErr := "LearnOrg successfully added in addlearnOrg in crudoperations: " + string(bs)
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 0
		}
	} else {
		theErr := "Error adding LearnOrg; could not perform CRUD or OrgID was bad: " + strconv.Itoa(postedLearnorg.OrgID)
		logWriter(theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.SuccOrFail = 1
	}

	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in addUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//This deletes a LearnOrg to our database; called from anywhere
func deleteLearnOrg(w http.ResponseWriter, req *http.Request) {
	canCrud := true //Used to determine if we're good to try our crud operation

	//Declare data to return
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	theReturnMessage := ReturnMessage{}

	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from deleteLearnOrg: " + err.Error() + "\n" + string(bs)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}
	//Declare JSON we're looking for
	type OrgDelete struct {
		OrgID int `json:"OrgID"`
	}
	//Marshal it into our type
	var postedLearnOrgID OrgDelete
	json.Unmarshal(bs, &postedLearnOrgID)

	//Delete only if we had no issues above
	if canCrud && postedLearnOrgID.OrgID > 0 {
		//Search for User and delete
		orgCollection := mongoClient.Database("learnR").Collection("learnorg") //Here's our collection
		deletes := []bson.M{
			{"orgid": postedLearnOrgID.OrgID},
		} //Here's our filter to look for
		deletes = append(deletes, bson.M{"orgid": bson.M{
			"$eq": postedLearnOrgID.OrgID,
		}}, bson.M{"orgid": bson.M{
			"$eq": postedLearnOrgID.OrgID,
		}},
		)

		// create the slice of write models
		var writes []mongo.WriteModel

		for _, del := range deletes {
			model := mongo.NewDeleteManyModel().SetFilter(del)
			writes = append(writes, model)
		}

		// run bulk write
		bulkWrite, err := orgCollection.BulkWrite(theContext, writes)
		if err != nil {
			theErr := "Error writing delete LearnROrg in deleteLearnorg in crudoperations: " + err.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			//Check to see if delete count worked; must have deleted at least one
			resultInt := bulkWrite.DeletedCount
			if resultInt > 0 {
				theErr := "LearnROrg successfully deleted in deletelearnorg in crudoperations: " + string(bs)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 0
			} else {
				theErr := "No documents deleted for this given learnOrg: " + strconv.Itoa(postedLearnOrgID.OrgID)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			}
		}
	} else {
		theErr := "Error, could not CRUD operate in deleteOrg, or the number we recieved was wrong: " +
			strconv.Itoa(postedLearnOrgID.OrgID)
		logWriter(theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.SuccOrFail = 1
	}

	//Write the response back
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in deleteUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//This updates a LearnOrg to our database; called from anywhere
func updateLearnOrg(w http.ResponseWriter, req *http.Request) {
	canCrud := true
	//Declare data to return
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	theReturnMessage := ReturnMessage{}

	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from updateLearnOrg: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Marshal it into our type
	var theLearnOrgUpdate LearnrOrg
	json.Unmarshal(bs, &theLearnOrgUpdate)

	//Update LearnOrg if we have successfully decoded from JSON
	if canCrud {
		//Update User if their User ID != 0 or nil
		if theLearnOrgUpdate.OrgID != 0 {
			//Update User
			theTimeNow := time.Now()
			learnOrgCollection := mongoClient.Database("learnR").Collection("learnorg") //Here's our collection
			theFilter := bson.M{
				"orgid": bson.M{
					"$eq": theLearnOrgUpdate.OrgID, // check if bool field has value of 'false'
				},
			}
			updatedDocument := bson.M{
				"$set": bson.M{
					"orgid":       theLearnOrgUpdate.OrgID,
					"name":        theLearnOrgUpdate.Name,
					"orggoals":    theLearnOrgUpdate.OrgGoals,
					"userlist":    theLearnOrgUpdate.UserList,
					"adminlist":   theLearnOrgUpdate.AdminList,
					"learnrlist":  theLearnOrgUpdate.LearnrList,
					"datecreated": theLearnOrgUpdate.DateCreated,
					"dateupdated": theTimeNow.Format("2006-01-02 15:04:05"),
				},
			}
			updateResult, err := learnOrgCollection.UpdateOne(theContext, theFilter, updatedDocument)

			if err != nil {
				theErr := "Error writing update LearnOrg in updateLearnOrg in crudoperations: " + err.Error()
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			} else {
				//Check to see if anything was updated; if not, return the error
				if updateResult.ModifiedCount < 1 {
					theErr := "No document updated with this ID: " + strconv.Itoa(theLearnOrgUpdate.OrgID)
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 1
				} else {
					theErr := "Learnorg successfully updated in updateLearnorg in crudoperations: " + string(bs) + "\n"
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 0
				}
			}
		} else {
			theErr := "The Org ID was not found: " + strconv.Itoa(theLearnOrgUpdate.OrgID)
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		}
	}

	//Send the response back
	theJSONMessage, err := json.Marshal(theReturnMessage)
	if err != nil {
		errIs := "Error formatting JSON for return in updateUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//This gets a User with a certain UserID
func getLearnOrg(w http.ResponseWriter, req *http.Request) {
	canCrud := true
	//Declare data to return
	type ReturnMessage struct {
		TheErr           []string  `json:"TheErr"`
		ResultMsg        []string  `json:"ResultMsg"`
		SuccOrFail       int       `json:"SuccOrFail"`
		ReturnedLearnOrg LearnrOrg `json:"ReturnedLearnOrg"`
	}
	theReturnMessage := ReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success

	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from updateLearnOrg: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Decalre JSON we recieve
	type LearnOrgID struct {
		TheLearnOrgID int `json:"TheLearnOrgID"`
	}

	//Marshal it into our type
	var theLearnOrgGet LearnOrgID
	json.Unmarshal(bs, &theLearnOrgGet)

	//If we successfully decoded, (and the UesrID is not 0) we can get our user
	if canCrud && theLearnOrgGet.TheLearnOrgID > 0 {
		/* Find the Learnorg with the given LearnorgID */
		var theLearnOrgReturned LearnrOrg                                           //Initialize Learnorg to be returned after Mongo query
		learnorgCollection := mongoClient.Database("learnR").Collection("learnorg") //Here's our collection
		theFilter := bson.M{
			"orgid": bson.M{
				"$eq": theLearnOrgGet.TheLearnOrgID, // check if bool field has value of 'false'
			},
		}
		findOptions := options.Find()
		findLearnOrg, err := learnorgCollection.Find(theContext, theFilter, findOptions)
		theFind := 0 //A counter to track how many users we find
		if findLearnOrg.Err() != nil || err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				stringUserID := strconv.Itoa(theLearnOrgGet.TheLearnOrgID)
				returnedErr := "For " + stringUserID + ", no Learnorg was returned: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedLearnOrg = LearnrOrg{}
			} else {
				stringUserID := strconv.Itoa(theLearnOrgGet.TheLearnOrgID)
				returnedErr := "For " + stringUserID + ", there was a Mongo Error: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedLearnOrg = LearnrOrg{}
			}
		} else {
			//Found Learnorg, decode to return
			for findLearnOrg.Next(theContext) {
				stringUserID := strconv.Itoa(theLearnOrgGet.TheLearnOrgID)
				err := findLearnOrg.Decode(&theLearnOrgReturned)
				if err != nil {
					returnedErr := "For " + stringUserID +
						", there was an error decoding document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedLearnOrg = LearnrOrg{}
				} else if len(theLearnOrgReturned.Name) <= 1 {
					returnedErr := "For " + stringUserID +
						", there was an no document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedLearnOrg = LearnrOrg{}
				} else {
					//Successful decode, do nothing
				}
				theFind = theFind + 1
			}
			findLearnOrg.Close(theContext)
		}

		if theFind <= 0 {
			//Error, return an error back and log it
			stringUserID := strconv.Itoa(theLearnOrgGet.TheLearnOrgID)
			returnedErr := "For " + stringUserID +
				", No LearnOrg was returned."
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 1
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
			theReturnMessage.ReturnedLearnOrg = LearnrOrg{}
		} else {
			//Success, log the success and return User
			stringUserID := strconv.Itoa(theLearnOrgGet.TheLearnOrgID)
			returnedErr := "For " + stringUserID +
				", Learnorg should be successfully decoded."
			//fmt.Println(returnedErr)
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 0
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ReturnedLearnOrg = theLearnOrgReturned
		}
	} else {
		//Error, return an error back and log it
		stringUserID := strconv.Itoa(theLearnOrgGet.TheLearnOrgID)
		returnedErr := "For " + stringUserID +
			", No LearnOrg was returned. Learnorg was also not accepted: " + stringUserID
		logWriter(returnedErr)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
		theReturnMessage.ReturnedLearnOrg = LearnrOrg{}
	}

	//Format the JSON map for returning our results
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in getUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//This gets all the LearnOrgs the User is Admin of
func getLearnOrgAdminOf(w http.ResponseWriter, req *http.Request) {
	canCrud := true
	//Declare data to return
	type TheReturnMessage struct {
		TheErr            []string    `json:"TheErr"`
		ResultMsg         []string    `json:"ResultMsg"`
		SuccOrFail        int         `json:"SuccOrFail"`
		ReturnedLearnOrgs []LearnrOrg `json:"ReturnedLearnOrgs"`
	}
	theReturnMessage := TheReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success

	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from getLearnOrgAdminOf: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Decalre JSON we recieve
	type TheAdminOrgs struct {
		TheIDS []int `json:"TheIDS"`
	}

	//Marshal it into our type
	var theitem TheAdminOrgs
	json.Unmarshal(bs, &theitem)

	//If we successfully decoded, (and the IDs are not 0
	if canCrud && len(theitem.TheIDS) > 0 {
		/* Call our 'getLearnROrg' function for everyID this User is Admin of, (unless the ID is 0) */
		for j := 0; j < len(theitem.TheIDS); j++ {
			goodIDGet := true
			type LearnOrgID struct {
				TheLearnOrgID int `json:"TheLearnOrgID"`
			}
			theID := LearnOrgID{TheLearnOrgID: theitem.TheIDS[j]}
			theJSONMessage, err := json.Marshal(theID)
			if err != nil {
				fmt.Println(err)
				logWriter(err.Error())
				log.Fatal(err)
				goodIDGet = false
			}
			payload := strings.NewReader(string(theJSONMessage))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			req, err := http.NewRequest("POST", "http://localhost:4000/getLearnOrg", payload)
			if err != nil {
				theErr := "There was an error getting LearnROrgs in loadLearnROrgs: " + err.Error()
				logWriter(theErr)
				fmt.Println(theErr)
				goodIDGet = false
			}
			req.Header.Add("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req.WithContext(ctx))

			if resp.StatusCode >= 300 || resp.StatusCode <= 199 {
				theErr := "There was an error reaching out to loadLearnROrg API: " + strconv.Itoa(resp.StatusCode)
				fmt.Println(theErr)
				logWriter(theErr)
				goodIDGet = false
			} else if err != nil {
				theErr := "Error from response to loadLearnROrg: " + err.Error()
				fmt.Println(theErr)
				logWriter(theErr)
				goodIDGet = false
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				theErr := "There was an error getting a response for LearnROrgs in loadLearnROrgs: " + err.Error()
				logWriter(theErr)
				fmt.Println(theErr)
				goodIDGet = false
			}

			//Marshal the response into a type we can read
			type ReturnMessage struct {
				TheErr           []string  `json:"TheErr"`
				ResultMsg        []string  `json:"ResultMsg"`
				SuccOrFail       int       `json:"SuccOrFail"`
				ReturnedLearnOrg LearnrOrg `json:"ReturnedLearnOrg"`
			}
			var returnedMessage ReturnMessage
			json.Unmarshal(body, &returnedMessage)

			//Evaluate if we can add this learnOrg to our map of returned LearnOrgs for this Admin User
			if !goodIDGet || returnedMessage.SuccOrFail != 0 || returnedMessage.ReturnedLearnOrg.OrgID == 0 {
				theErr := "Had an issue getting an ID for this Admin User and the LearnR Org: "
				for k := 0; k < len(returnedMessage.TheErr); k++ {
					theErr = theErr + returnedMessage.TheErr[k]
				}
				logWriter(theErr)
			} else {
				//Good ID, add it to return
				theReturnMessage.ReturnedLearnOrgs = append(theReturnMessage.ReturnedLearnOrgs, returnedMessage.ReturnedLearnOrg)
			}
		}
	} else {
		//Error, return an error back and log it
		returnedErr := "Can crud was not true or we had an error parsing IDs"
		logWriter(returnedErr)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
		theReturnMessage.ReturnedLearnOrgs = []LearnrOrg{}
	}

	//Format the JSON map for returning our results
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in getUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* LEARNR CRUD OPERATIONS ENDING */

/* This funciton returns a map of all OrgNames entered in our database when called,
(should be called wherever we need to create an LearnR Org) */
func giveAllLearnROrg(w http.ResponseWriter, req *http.Request) {
	//Declare data to return
	type ReturnMessage struct {
		TheErr             []string        `json:"TheErr"`
		ResultMsg          []string        `json:"ResultMsg"`
		SuccOrFail         int             `json:"SuccOrFail"`
		ReturnedOrgNameMap map[string]bool `json:"ReturnedOrgNameMap"`
	}
	theReturnMessage := ReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success

	//Declare empty map to fill and return
	orgNameMap := make(map[string]bool) //Clear Map for future use on page load

	learnROrgCollection := mongoClient.Database("learnR").Collection("learnorg") //Here's our collection

	//Query Mongo for all Users
	theFilter := bson.M{}
	findOptions := options.Find()
	currOrg, err := learnROrgCollection.Find(theContext, theFilter, findOptions)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			theErr := "No documents were returned for orgs in giveAllLearnROrgs in MongoDB: " + err.Error()
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 1
			logWriter(theErr)
		} else {
			theErr := "There was an error returning results for this Org, :" + err.Error()
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 1
			logWriter(theErr)
		}
	}
	//Loop over query results and fill User Array
	for currOrg.Next(theContext) {
		// create a value into which the single document can be decoded
		var aOrg LearnrOrg
		err := currOrg.Decode(&aOrg)
		if err != nil {
			theErr := "Error decoding Orgs in MongoDB in giveAllLearnROrgs: " + err.Error()
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 0
			logWriter(theErr)
		}
		//Fill Username map with the found Username
		orgNameMap[aOrg.Name] = true
	}
	// Close the cursor once finished
	currOrg.Close(theContext)

	//Check to see if anyusernames were returned or we have errors
	if theReturnMessage.SuccOrFail >= 1 {
		theErr := "There are a number of errors for returning these Organization Names..."
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
	} else if len(orgNameMap) <= 0 {
		theErr := "No usernames returned...this could be the site's first deployment with no Organizations!"
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
	} else {
		theErr := "No issues returning Organizations"
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 0
	}
	theReturnMessage.ReturnedOrgNameMap = orgNameMap //Add our final OrgMap

	//Format the JSON map for returning our results
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in giveAllLearnROrgs: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* This function returns a map of all LearnR names entered in our DB when called,
(should be called whenever we need to create a new LearnR)*/
func giveAllLearnr(w http.ResponseWriter, req *http.Request) {
	//Declare data to return
	type ReturnMessage struct {
		TheErr              []string        `json:"TheErr"`
		ResultMsg           []string        `json:"ResultMsg"`
		SuccOrFail          int             `json:"SuccOrFail"`
		ReturnedLearnRNames map[string]bool `json:"ReturnedLearnRNames"`
	}
	theReturnMessage := ReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success

	//Declare empty map to fill and return
	learnrNameMap := make(map[string]bool) //Clear Map for future use on page load

	collection := mongoClient.Database("learnR").Collection("learnr") //Here's our collection

	//Query Mongo for all Users
	theFilter := bson.M{}
	findOptions := options.Find()
	curr, err := collection.Find(theContext, theFilter, findOptions)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			theErr := "No documents were returned for orgs in givelearnrnames in MongoDB: " + err.Error()
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 1
			logWriter(theErr)
		} else {
			theErr := "There was an error returning results for this learnr, :" + err.Error()
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 1
			logWriter(theErr)
		}
	}
	//Loop over query results and fill User Array
	for curr.Next(theContext) {
		// create a value into which the single document can be decoded
		var item Learnr
		err := curr.Decode(&item)
		if err != nil {
			theErr := "Error decoding Learnr in MongoDB in giveAllLearnrs: " + err.Error()
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 0
			logWriter(theErr)
		}
		//Fill Username map with the found Username
		learnrNameMap[item.Name] = true
	}
	// Close the cursor once finished
	curr.Close(theContext)

	//Check to see if anyusernames were returned or we have errors
	if theReturnMessage.SuccOrFail >= 1 {
		theErr := "There are a number of errors for returning these Learnr Names..."
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
	} else if len(learnrNameMap) <= 0 {
		theErr := "No learnr returned...this could be the site's first deployment with no learnrs!"
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
	} else {
		theErr := "No issues returning learnrs"
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 0
	}
	theReturnMessage.ReturnedLearnRNames = learnrNameMap //Add our final OrgMap

	//Format the JSON map for returning our results
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in giveAllLearnROrgs: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* VALIDATION API BEGINNING */

/* This function returns a map of ALL Usernames entered in our database
when called, (should be on the index page ) */
func giveAllUsernames(w http.ResponseWriter, req *http.Request) {
	//Declare data to return
	type ReturnMessage struct {
		TheErr          []string        `json:"TheErr"`
		ResultMsg       []string        `json:"ResultMsg"`
		SuccOrFail      int             `json:"SuccOrFail"`
		ReturnedUserMap map[string]bool `json:"ReturnedUserMap"`
	}
	theReturnMessage := ReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success

	//Declare empty map to fill and return
	usernameMap := make(map[string]bool) //Clear Map for future use on page load

	userCollection := mongoClient.Database("learnR").Collection("users") //Here's our collection

	//Query Mongo for all Users
	theFilter := bson.M{}
	findOptions := options.Find()
	currUser, err := userCollection.Find(theContext, theFilter, findOptions)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			theErr := "No documents were returned for users in giveAllUsernames in MongoDB: " + err.Error()
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 1
			logWriter(theErr)
		} else {
			theErr := "There was an error returning results for this Users, :" + err.Error()
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 1
			logWriter(theErr)
		}
	}
	//Loop over query results and fill User Array
	for currUser.Next(theContext) {
		// create a value into which the single document can be decoded
		var aUser User
		err := currUser.Decode(&aUser)
		if err != nil {
			theErr := "Error decoding Users in MongoDB in giveAllUsernames: " + err.Error()
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 0
			logWriter(theErr)
		}
		//Fill Username map with the found Username
		usernameMap[aUser.UserName] = true
	}
	// Close the cursor once finished
	currUser.Close(theContext)

	//Check to see if anyusernames were returned or we have errors
	if theReturnMessage.SuccOrFail >= 1 {
		theErr := "There are a number of errors for returning these Usernames..."
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
	} else if len(usernameMap) <= 0 {
		theErr := "No usernames returned...this could be the site's first deployment with no users!"
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
	} else {
		theErr := "No issues returning Usernames"
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 0
	}
	theReturnMessage.ReturnedUserMap = usernameMap //Add our final Usermap

	//Format the JSON map for returning our results
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in giveAllUsernames: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* This function returns a map of ALL emails entered into our DB when called,
should be on the signup page */
func giveAllEmails(w http.ResponseWriter, req *http.Request) {
	//Declare data to return
	type ReturnMessage struct {
		TheErr           []string        `json:"TheErr"`
		ResultMsg        []string        `json:"ResultMsg"`
		SuccOrFail       int             `json:"SuccOrFail"`
		ReturnedEmailMap map[string]bool `json:"ReturnedEmailMap"`
	}
	theReturnMessage := ReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success

	//Declare empty map to fill and return
	emailMap := make(map[string]bool) //Clear Map for future use on page load

	userCollection := mongoClient.Database("learnR").Collection("users") //Here's our collection

	//Query Mongo for all Users
	theFilter := bson.M{}
	findOptions := options.Find()
	currUser, err := userCollection.Find(theContext, theFilter, findOptions)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			theErr := "No documents were returned for emails in giveAllEmails in MongoDB: " + err.Error()
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 1
			logWriter(theErr)
		} else {
			theErr := "There was an error returning results for these emails, :" + err.Error()
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 1
			logWriter(theErr)
		}
	}
	//Loop over query results and fill User Array
	for currUser.Next(theContext) {
		// create a value into which the single document can be decoded
		var aUser User
		err := currUser.Decode(&aUser)
		if err != nil {
			theErr := "Error decoding Users in MongoDB in giveAllEmails: " + err.Error()
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 0
			logWriter(theErr)
		}
		//Fill Email map with the found Email
		emailMap[aUser.Email[0]] = true
	}
	// Close the cursor once finished
	currUser.Close(theContext)

	//Check to see if any emails were returned or we have errors
	if theReturnMessage.SuccOrFail >= 1 {
		theErr := "There are a number of errors for returning these Emails..."
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
	} else if len(emailMap) <= 0 {
		theErr := "No emails returned...this could be the site's first deployment with no users!"
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
	} else {
		theErr := "No issues returning Emails"
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 0
	}
	theReturnMessage.ReturnedEmailMap = emailMap //Add our final Email map

	//Format the JSON map for returning our results
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in giveAllUsernames: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* This function searches with a Username and password to return a yes or no response
if the User is found; is so, we return the User, with a successful response.
If not, we return a failed response and an empty User profile */
func userLogin(w http.ResponseWriter, req *http.Request) {
	canCrud := true
	//Declare type to be returned later through JSON Response
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
		TheUser    User     `json:"TheUser"`
	}
	theResponseMessage := ReturnMessage{}
	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		logWriter(err.Error())
		canCrud = false
		theResponseMessage.ResultMsg = append(theResponseMessage.ResultMsg, "Could not get bytes")
		theResponseMessage.TheErr = append(theResponseMessage.TheErr, "Could not get bytes")
		theResponseMessage.SuccOrFail = 1
	}

	type LoginData struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}

	//Marshal the user data into our type
	var dataForLogin LoginData
	json.Unmarshal(bs, &dataForLogin)

	//Check for null values; exit program if password or Username are empty
	if canCrud && dataForLogin.Username != "" && dataForLogin.Password != "" {
		theUserReturned := User{} //Initialize User to be returned after Mongo query
		//Query for the User, given the userID for the User
		user_collection := mongoClient.Database("learnR").Collection("users") //Here's our collection
		theFilter := bson.M{
			"username": bson.M{
				"$eq": dataForLogin.Username, // check if bool field has value of 'false'
			},
			"password": bson.M{
				"$eq": dataForLogin.Password,
			},
		}
		findOptions := options.Find()
		findUser, err := user_collection.Find(theContext, theFilter, findOptions)
		theFind := 0 //A counter to track how many users we find
		if findUser.Err() != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				returnedErr := "For " + dataForLogin.Username + ", no User was returned: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theResponseMessage.SuccOrFail = 1
				theResponseMessage.ResultMsg = append(theResponseMessage.ResultMsg, returnedErr)
				theResponseMessage.TheErr = append(theResponseMessage.TheErr, returnedErr)
				theResponseMessage.TheUser = User{}
			} else {
				returnedErr := "For " + dataForLogin.Username + ", there was a Mongo Error: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theResponseMessage.SuccOrFail = 1
				theResponseMessage.ResultMsg = append(theResponseMessage.ResultMsg, returnedErr)
				theResponseMessage.TheErr = append(theResponseMessage.TheErr, returnedErr)
				theResponseMessage.TheUser = User{}
			}
		} else {
			//Set initial values so the decode function dosen't freak out
			theUserReturned.UserName = ""
			theUserReturned.Password = ""
			//Found User, decode to return
			for findUser.Next(theContext) {
				err := findUser.Decode(&theUserReturned)
				if err != nil {
					returnedErr := "For " + dataForLogin.Username +
						", there was an error decoding document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theResponseMessage.SuccOrFail = 1
					theResponseMessage.ResultMsg = append(theResponseMessage.ResultMsg, returnedErr)
					theResponseMessage.TheErr = append(theResponseMessage.TheErr, returnedErr)
					theResponseMessage.TheUser = User{}
				} else if len(theUserReturned.UserName) <= 1 {
					returnedErr := "For " + dataForLogin.Username +
						", there was an no document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theResponseMessage.SuccOrFail = 1
					theResponseMessage.ResultMsg = append(theResponseMessage.ResultMsg, returnedErr)
					theResponseMessage.TheErr = append(theResponseMessage.TheErr, returnedErr)
					theResponseMessage.TheUser = User{}
				} else {
					//Successful decode, do nothing
				}
				theFind = theFind + 1
			}
			findUser.Close(theContext)
		}

		if theFind <= 0 {
			//Error, return an error back and log it
			returnedErr := "For " + dataForLogin.Username +
				", No User was returned."
			fmt.Println(returnedErr)
			logWriter(returnedErr)
			theResponseMessage.SuccOrFail = 1
			theResponseMessage.ResultMsg = append(theResponseMessage.ResultMsg, returnedErr)
			theResponseMessage.TheErr = append(theResponseMessage.TheErr, returnedErr)
			theResponseMessage.TheUser = theUserReturned
		} else {
			//Success, log the success and return User
			returnedErr := "For " + dataForLogin.Username +
				", User should be successfully decoded."
			//fmt.Println(returnedErr)
			logWriter(returnedErr)
			theResponseMessage.SuccOrFail = 0
			theResponseMessage.ResultMsg = append(theResponseMessage.ResultMsg, returnedErr)
			theResponseMessage.TheErr = append(theResponseMessage.TheErr, returnedErr)
			theResponseMessage.TheUser = theUserReturned
		}
	} else {
		theResponseMessage.ResultMsg = append(theResponseMessage.ResultMsg,
			"Can Crud was false, or nil values: "+dataForLogin.Username+" "+dataForLogin.Password)
		theResponseMessage.TheErr = append(theResponseMessage.TheErr,
			"Can Crud was false, or nil values: "+dataForLogin.Username+" "+dataForLogin.Password)
		theResponseMessage.SuccOrFail = 1
	}

	//Errors/Success are recorded, User given, send JSON back
	theJSONMessage, err := json.Marshal(theResponseMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in userLogin: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//This should give a random id value to both food groups
func randomIDCreationAPI(w http.ResponseWriter, req *http.Request) {
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
		RandomID   int      `json:"RandomID"`
	}
	theReturnMessage := ReturnMessage{}
	finalID := 0        //The final, unique ID to return to the food/user
	randInt := 0        //The random integer added onto ID
	randIntString := "" //The integer built through a string...
	min, max := 0, 9    //The min and Max value for our randInt
	foundID := false
	for !foundID {
		randInt = 0
		randIntString = ""
		//Create the random number, convert it to string
		for i := 0; i < 12; i++ {
			randInt = rand.Intn(max-min) + min
			randIntString = randIntString + strconv.Itoa(randInt)
		}
		//Once we have a string of numbers, we can convert it back to an integer
		theID, err := strconv.Atoi(randIntString)
		if err != nil {
			fmt.Printf("We got an error converting a string back to a number, %v\n", err)
			fmt.Printf("Here is randInt: %v\n and randIntString: %v\n", randInt, randIntString)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "Error converting number")
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, "Error converting number")
			fmt.Println(err)
			log.Fatal(err)
			return
		}
		//Search all our collections to see if this UserID is unique
		canExit := []bool{true, true}
		/* User collection */
		userCollection := mongoClient.Database("learnR").Collection("users") //Here's our collection
		var testAUser User
		theErr := userCollection.FindOne(theContext, bson.M{"userid": theID}).Decode(&testAUser)
		if theErr != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				canExit[0] = true
			} else {
				theErr := "There is another error getting random ID: " + err.Error()
				logWriter(theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				canExit[0] = false
				log.Fatal(theErr)
			}
		}
		/* LearnROrg collection */
		learnRCollection := mongoClient.Database("learnR").Collection("learnorg") //Here's our collection
		var testALearnROrg LearnrOrg
		theErr2 := learnRCollection.FindOne(theContext, bson.M{"orgid": theID}).Decode(&testALearnROrg)
		if theErr2 != nil {
			if strings.Contains(theErr.Error(), "no documents in result") {
				canExit[0] = true
			} else {
				theErr := "There is another error getting random ID: " + err.Error()
				logWriter(theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				canExit[1] = false
				log.Fatal(theErr)
			}
		}
		//Final check to see if we can exit this loop
		if canExit[0] && canExit[1] {
			finalID = theID
			foundID = true
			theReturnMessage.RandomID = finalID
			theReturnMessage.SuccOrFail = 0
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, "Good new random ID added")
		} else {
			foundID = false
		}
	}

	/* Return the marshaled response */
	//Send the response back
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in randomIDCreationAPI: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//This should add a verification code to our DB
func addEmailVerif(w http.ResponseWriter, req *http.Request) {
	canCrud := true //Used to determine if we're good to try our crud operation

	//Declare data to return
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	theReturnMessage := ReturnMessage{}

	//Collect JSON from Postman or wherever
	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from emailVerif: " + err.Error() + "\n" + string(bs)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false //Reading failed, need to return failure
	}
	//Marshal it into our type
	var postedEmailVerif EmailVerify
	json.Unmarshal(bs, &postedEmailVerif)

	//Check to see if we can perform CRUD operations and we aren't passing a null LearnOrg
	if canCrud && postedEmailVerif.ID > 0 {
		collection := mongoClient.Database("learnR").Collection("emailverifs") //Here's our collection
		collectionStuff := []interface{}{postedEmailVerif}
		//Insert Our Data
		_, err2 := collection.InsertMany(theContext, collectionStuff)

		if err2 != nil {
			theErr := "Error adding Emailverif in addEmailVerif in crudoperations API: " + err2.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			theErr := "Email Verification successfully added in addEmailVerif in crudoperations: " + string(bs)
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 0
		}
	} else {
		theErr := "Error adding Email Verif; could not perform CRUD or OrgID was bad: " + strconv.Itoa(postedEmailVerif.ID)
		logWriter(theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.SuccOrFail = 1
	}

	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in addUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//This gets our Email Verif
func getEmailVerif(w http.ResponseWriter, req *http.Request) {
	canCrud := true
	//Declare data to return
	type ReturnMessage struct {
		TheErr              []string    `json:"TheErr"`
		ResultMsg           []string    `json:"ResultMsg"`
		SuccOrFail          int         `json:"SuccOrFail"`
		ReturnedEmailVerify EmailVerify `json:"ReturnedEmailVerify"`
	}
	theReturnMessage := ReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success

	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from updateLearnOrg: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Decalre JSON we recieve
	type EmailVerifID struct {
		TheEmailVerifID int `json:"TheEmailVerifID"`
	}

	//Marshal it into our type
	var theEmailVerifGet EmailVerifID
	json.Unmarshal(bs, &theEmailVerifGet)

	//If we successfully decoded, (and the Email Verif is not 0) we can get our Email Verification
	if canCrud && theEmailVerifGet.TheEmailVerifID > 0 {
		/* Find the Email Verif with the given ID */
		var theEmailVerifReturned EmailVerify                                  //Initialize value to be returned
		collection := mongoClient.Database("learnR").Collection("emailverifs") //Here's our collection
		theFilter := bson.M{
			"id": bson.M{
				"$eq": theEmailVerifGet.TheEmailVerifID, // check if bool field has value of 'false'
			},
		}
		findOptions := options.Find()
		find, err := collection.Find(theContext, theFilter, findOptions)
		theFind := 0 //A counter to track how many users we find
		if find.Err() != nil || err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				stringUserID := strconv.Itoa(theEmailVerifGet.TheEmailVerifID)
				returnedErr := "For " + stringUserID + ", no Email Verify was returned: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedEmailVerify = EmailVerify{}
			} else {
				stringUserID := strconv.Itoa(theEmailVerifGet.TheEmailVerifID)
				returnedErr := "For " + stringUserID + ", there was a Mongo Error: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedEmailVerify = EmailVerify{}
			}
		} else {
			//Found EmailVerif, decode to return
			for find.Next(theContext) {
				stringUserID := strconv.Itoa(theEmailVerifGet.TheEmailVerifID)
				err := find.Decode(&theEmailVerifReturned)
				if err != nil {
					returnedErr := "For " + stringUserID +
						", there was an error decoding document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedEmailVerify = EmailVerify{}
				} else if theEmailVerifReturned.ID <= 1 {
					returnedErr := "For " + stringUserID +
						", there was an no document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedEmailVerify = EmailVerify{}
				} else {
					//Successful decode, do nothing
				}
				theFind = theFind + 1
			}
			find.Close(theContext)
		}

		if theFind <= 0 {
			//Error, return an error back and log it
			stringUserID := strconv.Itoa(theEmailVerifGet.TheEmailVerifID)
			returnedErr := "For " + stringUserID +
				", No Email Verify was returned."
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 1
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
			theReturnMessage.ReturnedEmailVerify = EmailVerify{}
		} else {
			//Success, log the success and return User
			stringUserID := strconv.Itoa(theEmailVerifGet.TheEmailVerifID)
			returnedErr := "For " + stringUserID +
				", Email Verify should be successfully decoded."
			//fmt.Println(returnedErr)
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 0
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ReturnedEmailVerify = theEmailVerifReturned
		}
	} else {
		//Error, return an error back and log it
		stringUserID := strconv.Itoa(theEmailVerifGet.TheEmailVerifID)
		returnedErr := "For " + stringUserID +
			", No Email Verify was returned. Email Verify was also not accepted: " + stringUserID
		logWriter(returnedErr)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
		theReturnMessage.ReturnedEmailVerify = EmailVerify{}
	}

	//Format the JSON map for returning our results
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in getUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//This deletes our Email Verify
func deleteEmailVerify(w http.ResponseWriter, req *http.Request) {
	canCrud := true //Used to determine if we're good to try our crud operation

	//Declare data to return
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	theReturnMessage := ReturnMessage{}

	//Get the byte slice from the request body ajax
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from deleteEmailVerify: " + err.Error() + "\n" + string(bs)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}
	//Declare JSON we're looking for
	type ObjDelete struct {
		ID int `json:"ID"`
	}
	//Marshal it into our type
	var postedEmailVerifID ObjDelete
	json.Unmarshal(bs, &postedEmailVerifID)

	//Delete only if we had no issues above
	if canCrud && postedEmailVerifID.ID > 0 {
		//Search for User and delete
		collection := mongoClient.Database("learnR").Collection("emailverifs") //Here's our collection
		deletes := []bson.M{
			{"id": postedEmailVerifID.ID},
		} //Here's our filter to look for
		deletes = append(deletes, bson.M{"id": bson.M{
			"$eq": postedEmailVerifID.ID,
		}}, bson.M{"id": bson.M{
			"$eq": postedEmailVerifID.ID,
		}},
		)

		// create the slice of write models
		var writes []mongo.WriteModel

		for _, del := range deletes {
			model := mongo.NewDeleteManyModel().SetFilter(del)
			writes = append(writes, model)
		}

		// run bulk write
		bulkWrite, err := collection.BulkWrite(theContext, writes)
		if err != nil {
			theErr := "Error writing delete Email Verification in deleteEmailVerif in crudoperations: " + err.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			//Check to see if delete count worked; must have deleted at least one
			resultInt := bulkWrite.DeletedCount
			if resultInt > 0 {
				theErr := "Email Verification successfully deleted in deleteEmail Verif in crudoperations: " + string(bs)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 0
			} else {
				theErr := "No documents deleted for this given email Verificatition: " + strconv.Itoa(postedEmailVerifID.ID)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			}
		}
	} else {
		theErr := "Error, could not CRUD operate in deleteEmail Veriifctation, or the number we recieved was wrong: " +
			strconv.Itoa(postedEmailVerifID.ID)
		logWriter(theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.SuccOrFail = 1
	}

	//Write the response back
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in deleteUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

//This updates the Email Verify
func updateEmailVerify(w http.ResponseWriter, req *http.Request) {
	canCrud := true
	//Declare data to return
	type ReturnMessage struct {
		TheErr     []string `json:"TheErr"`
		ResultMsg  []string `json:"ResultMsg"`
		SuccOrFail int      `json:"SuccOrFail"`
	}
	theReturnMessage := ReturnMessage{}

	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from updateLearnOrg: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Marshal it into our type
	var theemailVerifUpdate EmailVerify
	json.Unmarshal(bs, &theemailVerifUpdate)

	//Update email verification if we have successfully decoded from JSON
	if canCrud {
		//Update Email verif if their ID != 0 or nil
		if theemailVerifUpdate.ID != 0 {
			//Update emailverif
			collection := mongoClient.Database("learnR").Collection("emailverifs") //Here's our collection
			theFilter := bson.M{
				"id": bson.M{
					"$eq": theemailVerifUpdate.ID, // check if bool field has value of 'false'
				},
			}
			updatedDocument := bson.M{
				"$set": bson.M{
					"username": theemailVerifUpdate.Username,
					"email":    theemailVerifUpdate.Email,
					"id":       theemailVerifUpdate.ID,
					"timemade": theemailVerifUpdate.TimeMade,
					"active":   theemailVerifUpdate.Active,
				},
			}
			updateResult, err := collection.UpdateOne(theContext, theFilter, updatedDocument)

			if err != nil {
				theErr := "Error writing update Email Verif in updateEmailVerify in crudoperations: " + err.Error()
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			} else {
				//Check to see if anything was updated; if not, return the error
				if updateResult.ModifiedCount < 1 {
					theErr := "No document updated with this ID: " + strconv.Itoa(theemailVerifUpdate.ID)
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 1
				} else {
					theErr := "Email Verification successfully updated in updateEmailVerify in crudoperations: " + string(bs) + "\n"
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 0
				}
			}
		} else {
			theErr := "The Org ID was not found: " + strconv.Itoa(theemailVerifUpdate.ID)
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		}
	}

	//Send the response back
	theJSONMessage, err := json.Marshal(theReturnMessage)
	if err != nil {
		errIs := "Error formatting JSON for return in updateUser: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* VALIDATION API ENDING */
