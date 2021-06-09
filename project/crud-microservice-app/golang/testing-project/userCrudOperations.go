package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	fmt.Printf("DEBUG: We've connected to the MongoDB\n")

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

	if canCrud {
		user_collection := mongoClient.Database("learnR").Collection("users") //Here's our collection
		collectedUsers := []interface{}{postedUser}
		//Insert Our Data
		_, err2 := user_collection.InsertMany(context.TODO(), collectedUsers)

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
	if canCrud {
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
			theErr := "User successfully deleted in deleteUser in crudoperations: " + string(bs)
			fmt.Printf("DEBUG: %v . Here is the amount deleted:%v\n", theErr, bulkWrite.DeletedCount)
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 0
		}
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
		updatedInfo, err := userCollection.UpdateOne(theContext, theFilter, updatedDocument)

		if err != nil {
			theErr := "Error writing update User in updateUser in crudoperations: " + err.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			theErr := "User successfully updated in updateUser in crudoperations: " + string(bs)
			fmt.Printf("DEBUG: %v. Here is the update results: %v\n", theErr, updatedInfo.ModifiedCount)
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 0
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
			fmt.Printf("DEBUG: %v\n", theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 1
			logWriter(theErr)
		} else {
			theErr := "There was an error returning results for this Users, :" + err.Error()
			fmt.Printf("DEBUG: %v\n", theErr)
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

/* VALIDATION API ENDING */
