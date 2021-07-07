package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

//LearnR Org
type LearnrOrg struct {
	OrgID       int      `json:"OrgID"` //Unique ID of this organization
	Name        string   `json:"Name"`  //Name of this organization
	OrgGoals    []string //A list of goals for this organization
	UserList    []int    //All the Users in this organization
	AdminList   []int    //A list of all the Admins in this organization,(UserIDs)
	LearnrList  []int    //A list of all learnr ints in this organization
	DateCreated string   `json:"DateCreated"`
	DateUpdated string   `json:"DateUpdated"`
}

//LearnR
type Learnr struct {
	ID            int             `json:"ID"`            //ID of this LearnR
	InfoID        int             `json:"InfoID"`        //Links to the LearnRInfo object which holds data
	OrgID         int             `json:"OrgID"`         //Which organization does this belong to
	Name          string          `json:"Name"`          //Name of this LearnR
	Tags          []string        `json:"Tags"`          //Tags that describe this LearnR
	Description   []string        `json:"Description"`   //Description of this LearnR
	PhoneNums     []string        `json:"PhoneNums"`     //Phone Nums attatched to this LearnR
	LearnRInforms []LearnRInforms `json:"LearnRInforms"` //What we'll text to our Users
	Active        bool            `json:"Active"`        //Whether this LearnR is still active
	DateCreated   string          `json:"DateCreated"`
	DateUpdated   string          `json:"DateUpdated"`
}

//LearnRInfo
type LearnrInfo struct {
	ID               int             `json:"ID"`               //ID of this LearnR Info
	LearnRID         int             `json:"LearnRID"`         //The LearnR ID related to this info
	AllSessions      []LearnRSession `json:"AllSessions"`      //An array of all the sessions
	FinishedSessions []LearnRSession `json:"FinishedSessions"` //An array of complete sessions only
	DateCreated      string          `json:"DateCreated"`
	DateUpdated      string          `json:"DateUpdated"`
}

//LearnRSession
type LearnRSession struct {
	ID               int             `json:"ID"`               //ID of this session
	LearnRID         int             `json:"LearnRID"`         //ID of this LearnR
	LearnRName       string          `json:"LearnRName"`       //Name of this LearnR
	TheLearnR        Learnr          `json:"TheLearnR"`        //The actual LearnR
	TheUser          User            `json:"TheUser"`          //Who is the User that sent this LearnR to someone?
	TargetUserNumber string          `json:"TargetUserNumber"` //User this session started to
	Ongoing          bool            `json:"Ongoing"`          //Is this session ongoing? Determined by time
	TextsSent        []LearnRInforms `json:"TextsSent"`        //All the Informs our program sent to User
	UserResponses    []string        `json:"UserResponses"`    //All the text responses sent by the User
	DateCreated      string          `json:"DateCreated"`
	DateUpdated      string          `json:"DateUpdated"`
}

//LearnRInforms
type LearnRInforms struct {
	ID          int    `json:"ID"`         //ID of this Inform
	Name        string `json:"Name"`       //Name of this Inform
	LearnRID    int    `json:"LearnRID"`   //ID of the LearnR this belongs to
	LearnRName  string `json:"LearnRName"` //Name this LearnR belongs to
	Order       int    `json:"Order"`      //The Order in the LearnR this will be
	TheInfo     string `json:"TheInfo"`    //What you want to say to someone
	ShouldWait  bool   `json:"ShouldWait"` //Should this info wait for User Response?
	WaitTime    int    `json:"WaitTime"`   //How much time should User be given to read this text?
	DateCreated string `json:"DateCreated"`
	DateUpdated string `json:"DateUpdated"`
}

/* BEGINNING LEARNR CRUD OPERATIONS */

//This adds a learnR to our DB; called from anywhere
func addLearnR(w http.ResponseWriter, req *http.Request) {
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
		theErr := "Error reading the request from learnR: " + err.Error() + "\n" + string(bs)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false //Reading failed, need to return failure
	}
	//Marshal it into our type
	var postedLearnR Learnr
	json.Unmarshal(bs, &postedLearnR)

	//Check to see if we can perform CRUD operations and we aren't passing a null LearnR
	if canCrud && postedLearnR.ID > 0 {
		theCollection := mongoClient.Database("learnR").Collection("learnr") //Here's our collection
		collectedInterface := []interface{}{postedLearnR}
		//Insert Our Data
		_, err2 := theCollection.InsertMany(theContext, collectedInterface)

		if err2 != nil {
			theErr := "Error adding LearnR in addLearnR in crudoperations API: " + err2.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			theErr := "LearnR successfully added in addlearnR in crudoperations: " + string(bs)
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 0
		}
	} else {
		theErr := "Error adding LearnR; could not perform CRUD or OrgID was bad: " + strconv.Itoa(postedLearnR.OrgID)
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

//This deletes a LearnR to our database; called from anywhere
func deleteLearnR(w http.ResponseWriter, req *http.Request) {
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
		theErr := "Error reading the request from deleteLearnR: " + err.Error() + "\n" + string(bs)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}
	//Declare JSON we're looking for
	type LearnRDelete struct {
		ID int `json:"ID"`
	}
	//Marshal it into our type
	var postedType LearnRDelete
	json.Unmarshal(bs, &postedType)

	//Delete only if we had no issues above
	if canCrud && postedType.ID > 0 {
		//Search for User and delete
		collection := mongoClient.Database("learnR").Collection("learnr") //Here's our collection
		deletes := []bson.M{
			{"id": postedType.ID},
		} //Here's our filter to look for
		deletes = append(deletes, bson.M{"id": bson.M{
			"$eq": postedType.ID,
		}}, bson.M{"id": bson.M{
			"$eq": postedType.ID,
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
			theErr := "Error writing delete learnr in deleteLear in crudoperations: " + err.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			//Check to see if delete count worked; must have deleted at least one
			resultInt := bulkWrite.DeletedCount
			if resultInt > 0 {
				theErr := "LearnR successfully deleted in deletelear in crudoperations: " + string(bs)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 0
			} else {
				theErr := "No documents deleted for this given learnR: " + strconv.Itoa(postedType.ID)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			}
		}
	} else {
		theErr := "Error, could not CRUD operate in deleteLearnR, or the number we recieved was wrong: " +
			strconv.Itoa(postedType.ID)
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

//This updates a LearnR to our database; called from anywhere
func updateLearnR(w http.ResponseWriter, req *http.Request) {
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
		theErr := "Error reading the request from updateLearnR: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Marshal it into our type
	var theTypePosted Learnr
	json.Unmarshal(bs, &theTypePosted)

	//Update LearnR if we have successfully decoded from JSON
	if canCrud {
		//Update LearnR if their LearnR ID != 0 or nil
		if theTypePosted.ID != 0 {
			//Update User
			theTimeNow := time.Now()
			collection := mongoClient.Database("learnR").Collection("learnr") //Here's our collection
			theFilter := bson.M{
				"id": bson.M{
					"$eq": theTypePosted.ID, // check if bool field has value of 'false'
				},
			}
			updatedDocument := bson.M{
				"$set": bson.M{
					"id":            theTypePosted.ID,
					"infoid":        theTypePosted.InfoID,
					"orgid":         theTypePosted.OrgID,
					"name":          theTypePosted.Name,
					"tags":          theTypePosted.Tags,
					"description":   theTypePosted.Description,
					"phonenums":     theTypePosted.PhoneNums,
					"learnrinforms": theTypePosted.LearnRInforms,
					"active":        theTypePosted.Active,
					"datecreated":   theTypePosted.DateCreated,
					"dateupdated":   theTimeNow.Format("2006-01-02 15:04:05"),
				},
			}
			updateResult, err := collection.UpdateOne(theContext, theFilter, updatedDocument)

			if err != nil {
				theErr := "Error writing update LearnR in updatelearnR in crudoperations: " + err.Error()
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			} else {
				//Check to see if anything was updated; if not, return the error
				if updateResult.ModifiedCount < 1 {
					theErr := "No document updated with this ID: " + strconv.Itoa(theTypePosted.ID)
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 1
				} else {
					theErr := "LearnR successfully updated in updateLearnR in crudoperations: " + string(bs) + "\n"
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 0
				}
			}
		} else {
			theErr := "The LearnR ID was not found: " + strconv.Itoa(theTypePosted.ID)
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

//This gets a LearnR with a certain LearnR ID
func getLearnR(w http.ResponseWriter, req *http.Request) {
	canCrud := true
	//Declare data to return
	type ReturnMessage struct {
		TheErr         []string `json:"TheErr"`
		ResultMsg      []string `json:"ResultMsg"`
		SuccOrFail     int      `json:"SuccOrFail"`
		ReturnedLearnR Learnr   `json:"ReturnedLearnR"`
	}
	theReturnMessage := ReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success

	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from getLearnR: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Decalre JSON we recieve
	type LearnRID struct {
		ID int `json:"ID"`
	}

	//Marshal it into our type
	var typePosted LearnRID
	json.Unmarshal(bs, &typePosted)

	//If we successfully decoded, (and the ID is not 0) we can get our item
	if canCrud && typePosted.ID > 0 {
		/* Find the Learnorg with the given LearnorgID */
		var itemReturned Learnr                                           //Initialize Item to be returned after Mongo query
		collection := mongoClient.Database("learnR").Collection("learnr") //Here's our collection
		theFilter := bson.M{
			"id": bson.M{
				"$eq": typePosted.ID, // check if bool field has value of 'false'
			},
		}
		findOptions := options.Find()
		find, err := collection.Find(theContext, theFilter, findOptions)
		theFind := 0 //A counter to track how many users we find
		if find.Err() != nil || err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				stringUserID := strconv.Itoa(typePosted.ID)
				returnedErr := "For " + stringUserID + ", no LearnR was returned: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedLearnR = Learnr{}
			} else {
				stringUserID := strconv.Itoa(typePosted.ID)
				returnedErr := "For " + stringUserID + ", there was a Mongo Error: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedLearnR = Learnr{}
			}
		} else {
			//Found Learnorg, decode to return
			for find.Next(theContext) {
				stringid := strconv.Itoa(typePosted.ID)
				err := find.Decode(&itemReturned)
				if err != nil {
					returnedErr := "For " + stringid +
						", there was an error decoding document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedLearnR = Learnr{}
				} else if len(itemReturned.Name) <= 1 {
					returnedErr := "For " + stringid +
						", there was an no document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedLearnR = Learnr{}
				} else {
					//Successful decode, do nothing
				}
				theFind = theFind + 1
			}
			find.Close(theContext)
		}

		if theFind <= 0 {
			//Error, return an error back and log it
			stringID := strconv.Itoa(typePosted.ID)
			returnedErr := "For " + stringID +
				", No LearnR was returned."
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 1
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
			theReturnMessage.ReturnedLearnR = Learnr{}
		} else {
			//Success, log the success and return User
			stringID := strconv.Itoa(typePosted.ID)
			returnedErr := "For " + stringID +
				", Learnr should be successfully decoded."
			//fmt.Println(returnedErr)
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 0
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ReturnedLearnR = itemReturned
		}
	} else {
		//Error, return an error back and log it
		theIDString := strconv.Itoa(typePosted.ID)
		returnedErr := "For " + theIDString +
			", No LearnR was returned. Learnorg was also not accepted: " + theIDString
		logWriter(returnedErr)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
		theReturnMessage.ReturnedLearnR = Learnr{}
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

/* This function takes parameters for LearnR sorting and returns an array */
func specialLearnRGive(w http.ResponseWriter, req *http.Request) {
	canCrud := true
	//Declare data to return
	type ReturnMessage struct {
		TheErr          []string `json:"TheErr"`
		ResultMsg       []string `json:"ResultMsg"`
		SuccOrFail      int      `json:"SuccOrFail"`
		ReturnedLearnrs []Learnr `json:"ReturnedLearnrs"`
	}
	theReturnMessage := ReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success

	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from specialLearnRGive: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Decalre JSON we recieve
	type TheSpecialCases struct {
		CaseSearch       []int  `json:"CaseSearch"`
		OrganizationName string `json:"OrganizationName"`
		Tag              string `json:"Tag"`
		LearnRName       string `json:"LearnRName"`
		EntryAmountFrom  int    `json:"EntryAmountFrom"`
		EntryAmountTo    int    `json:"EntryAmountTo"`
	}

	//Marshal it into our type
	var theitem TheSpecialCases
	json.Unmarshal(bs, &theitem)

	//Do CRUD operations if allowed
	if canCrud {
		/* Begin building crud operation based on our criteria */
		collection := mongoClient.Database("learnR").Collection("learnr") //Here's our collection
		theFilter := bson.M{}
		findOptions := options.Find()
		theFind := 0 //A counter to track how many Learnrs we find

		/* Alter the filter/findoptions based on criteria */
		if theitem.CaseSearch[0] == 0 {
			//Do nothing, just get all Learnrs
		}
		/*DEBUG: Add cases later for more criteria */
		/* Run the mongo query after fixed filter/findoptions */
		find, err := collection.Find(theContext, theFilter, findOptions)
		if find.Err() != nil || err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				returnedErr := "No documents returned; may be that there are no Learnrs yet or search was bad..."
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 0
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedLearnrs = []Learnr{}
			} else {
				returnedErr := "There was a Mongo Error: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedLearnrs = []Learnr{}
			}
		} else {
			//Found Learnr, decode to return
			for find.Next(theContext) {
				var theLearnRReturned Learnr
				err := find.Decode(&theLearnRReturned)
				if err != nil {
					returnedErr := "For " +
						", there was an error decoding document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedLearnrs = []Learnr{}
				} else if len(theLearnRReturned.Name) <= 1 {
					returnedErr := "For " +
						", there was an no document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedLearnrs = []Learnr{}
				} else {
					//Successful decode, add this to our array
					theReturnMessage.ReturnedLearnrs = append(theReturnMessage.ReturnedLearnrs, theLearnRReturned)
				}
				theFind = theFind + 1
			}
			find.Close(theContext)
		}
		//Declare results, see if we have errors
		if theReturnMessage.SuccOrFail >= 1 {
			theErr := "There are a number of errors for returning these special Learnrs..."
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		} else if len(theReturnMessage.ReturnedLearnrs) <= 0 {
			theErr := "No learnrs returned...this could be the site's first deployment with no Learnrs or a faulty search!"
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 0
		} else {
			theErr := "No issues returning Learnrs"
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.SuccOrFail = 0
		}
	} else {
		//Error, return an error back and log it
		returnedErr := "Had an issue getting special Learnrs"
		logWriter(returnedErr)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
		theReturnMessage.ReturnedLearnrs = []Learnr{}
	}

	//Format the JSON map for returning our results
	theJSONMessage, err := json.Marshal(theReturnMessage)
	//Send the response back
	if err != nil {
		errIs := "Error formatting JSON for return in giveAllLearnROrgs: " + err.Error()
		logWriter(errIs)
	}
	fmt.Fprint(w, string(theJSONMessage))
}

/* ENDING LEARNR CRUD OPERATIONS */

/* BEGINNING LEARNRINFO CRUD OPERATIONS */

func addLearnrInfo(w http.ResponseWriter, req *http.Request) {
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
		theErr := "Error reading the request from learnRiNFO: " + err.Error() + "\n" + string(bs)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false //Reading failed, need to return failure
	}
	//Marshal it into our type
	var postedType LearnrInfo
	json.Unmarshal(bs, &postedType)

	//Check to see if we can perform CRUD operations and we aren't passing a null LearnR
	if canCrud && postedType.ID > 0 {
		theCollection := mongoClient.Database("learnR").Collection("learnrinfo") //Here's our collection
		collectedInterface := []interface{}{postedType}
		//Insert Our Data
		_, err2 := theCollection.InsertMany(theContext, collectedInterface)

		if err2 != nil {
			theErr := "Error adding LearnRInfo in addLearnRInfo in crudoperations API: " + err2.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			theErr := "LearnRInfo successfully added in addlearnRInfo in crudoperations: " + string(bs)
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 0
		}
	} else {
		theErr := "Error adding LearnRInfo; could not perform CRUD or ID was bad: " + strconv.Itoa(postedType.ID)
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

//This deletes a LearnRInfo to our database; called from anywhere
func deleteLearnrInfo(w http.ResponseWriter, req *http.Request) {
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
		theErr := "Error reading the request from deleteLearnR: " + err.Error() + "\n" + string(bs)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}
	//Declare JSON we're looking for
	type LearnRInfoDelete struct {
		ID int `json:"ID"`
	}
	//Marshal it into our type
	var postedType LearnRInfoDelete
	json.Unmarshal(bs, &postedType)

	//Delete only if we had no issues above
	if canCrud && postedType.ID > 0 {
		//Search for User and delete
		collection := mongoClient.Database("learnR").Collection("learnrinfo") //Here's our collection
		deletes := []bson.M{
			{"id": postedType.ID},
		} //Here's our filter to look for
		deletes = append(deletes, bson.M{"id": bson.M{
			"$eq": postedType.ID,
		}}, bson.M{"id": bson.M{
			"$eq": postedType.ID,
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
			theErr := "Error writing delete learnrInfo in deleteLearnRInfo in crudoperations: " + err.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			//Check to see if delete count worked; must have deleted at least one
			resultInt := bulkWrite.DeletedCount
			if resultInt > 0 {
				theErr := "LearnRInfo successfully deleted in deletelearnRInfo in crudoperations: " + string(bs)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 0
			} else {
				theErr := "No documents deleted for this given learnRInfo: " + strconv.Itoa(postedType.ID)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			}
		}
	} else {
		theErr := "Error, could not CRUD operate in deleteLearnRInfo, or the number we recieved was wrong: " +
			strconv.Itoa(postedType.ID)
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

//This updates a LearnRInfo to our database; called from anywhere
func updateLearnrInfo(w http.ResponseWriter, req *http.Request) {
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
		theErr := "Error reading the request from updateLearnRInfo: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Marshal it into our type
	var theTypePosted LearnrInfo
	json.Unmarshal(bs, &theTypePosted)

	//Update item if we have successfully decoded from JSON
	if canCrud {
		//Update LearnR if their LearnR ID != 0 or nil
		if theTypePosted.ID != 0 {
			//Update User
			theTimeNow := time.Now()
			collection := mongoClient.Database("learnR").Collection("learnrinfo") //Here's our collection
			theFilter := bson.M{
				"id": bson.M{
					"$eq": theTypePosted.ID, // check if bool field has value of 'false'
				},
			}
			updatedDocument := bson.M{
				"$set": bson.M{
					"id":               theTypePosted.ID,
					"learnrid":         theTypePosted.LearnRID,
					"allsessions":      theTypePosted.AllSessions,
					"finishedsessions": theTypePosted.FinishedSessions,
					"datecreated":      theTypePosted.DateCreated,
					"dateupdated":      theTimeNow.Format("2006-01-02 15:04:05"),
				},
			}
			updateResult, err := collection.UpdateOne(theContext, theFilter, updatedDocument)

			if err != nil {
				theErr := "Error writing update LearnRInfo in updatelearnRInfo in crudoperations: " + err.Error()
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			} else {
				//Check to see if anything was updated; if not, return the error
				if updateResult.ModifiedCount < 1 {
					theErr := "No document updated with this ID: " + strconv.Itoa(theTypePosted.ID)
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 1
				} else {
					theErr := "LearnRInfo successfully updated in updateLearnRInfo in crudoperations: " + string(bs) + "\n"
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 0
				}
			}
		} else {
			theErr := "The LearnRInfo ID was not found: " + strconv.Itoa(theTypePosted.ID)
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

//This gets a LearnRInfo with a certain LearnR ID
func getLearnrInfo(w http.ResponseWriter, req *http.Request) {
	canCrud := true
	//Declare data to return
	type ReturnMessage struct {
		TheErr             []string   `json:"TheErr"`
		ResultMsg          []string   `json:"ResultMsg"`
		SuccOrFail         int        `json:"SuccOrFail"`
		ReturnedLearnRInfo LearnrInfo `json:"ReturnedLearnRInfo"`
	}
	theReturnMessage := ReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success

	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from getLearnRInfo: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Decalre JSON we recieve
	type LearnRInfoID struct {
		ID int `json:"ID"`
	}

	//Marshal it into our type
	var typePosted LearnRInfoID
	json.Unmarshal(bs, &typePosted)

	//If we successfully decoded, (and the ID is not 0) we can get our item
	if canCrud && typePosted.ID > 0 {
		/* Find the LearnRInfo with the given ID */
		var itemReturned LearnrInfo                                           //Initialize Item to be returned after Mongo query
		collection := mongoClient.Database("learnR").Collection("learnrinfo") //Here's our collection
		theFilter := bson.M{
			"id": bson.M{
				"$eq": typePosted.ID, // check if bool field has value of 'false'
			},
		}
		findOptions := options.Find()
		find, err := collection.Find(theContext, theFilter, findOptions)
		theFind := 0 //A counter to track how many users we find
		if find.Err() != nil || err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				stringUserID := strconv.Itoa(typePosted.ID)
				returnedErr := "For " + stringUserID + ", no LearnRInfo was returned: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedLearnRInfo = LearnrInfo{}
			} else {
				stringUserID := strconv.Itoa(typePosted.ID)
				returnedErr := "For " + stringUserID + ", there was a Mongo Error: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedLearnRInfo = LearnrInfo{}
			}
		} else {
			//Found Learnorg, decode to return
			for find.Next(theContext) {
				stringid := strconv.Itoa(typePosted.ID)
				err := find.Decode(&itemReturned)
				if err != nil {
					returnedErr := "For " + stringid +
						", there was an error decoding document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedLearnRInfo = LearnrInfo{}
				} else if itemReturned.ID <= 1 {
					returnedErr := "For " + stringid +
						", there was an no document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedLearnRInfo = LearnrInfo{}
				} else {
					//Successful decode, do nothing
				}
				theFind = theFind + 1
			}
			find.Close(theContext)
		}

		if theFind <= 0 {
			//Error, return an error back and log it
			stringID := strconv.Itoa(typePosted.ID)
			returnedErr := "For " + stringID +
				", No LearnRInfo was returned."
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 1
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
			theReturnMessage.ReturnedLearnRInfo = LearnrInfo{}
		} else {
			//Success, log the success and return User
			stringID := strconv.Itoa(typePosted.ID)
			returnedErr := "For " + stringID +
				", LearnrInfo should be successfully decoded."
			//fmt.Println(returnedErr)
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 0
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ReturnedLearnRInfo = itemReturned
		}
	} else {
		//Error, return an error back and log it
		theIDString := strconv.Itoa(typePosted.ID)
		returnedErr := "For " + theIDString +
			", No LearnRInfo was returned. LearnInfo was also not accepted: " + theIDString
		logWriter(returnedErr)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
		theReturnMessage.ReturnedLearnRInfo = LearnrInfo{}
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

/* ENDING LEARNRINFO CRUD OPERATIONS */

/* BEGINNING LearnRSession CRUD OPERATIONS */

func addLearnRSession(w http.ResponseWriter, req *http.Request) {
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
		theErr := "Error reading the request from learnRSession: " + err.Error() + "\n" + string(bs)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false //Reading failed, need to return failure
	}
	//Marshal it into our type
	var postedType LearnRSession
	json.Unmarshal(bs, &postedType)

	//Check to see if we can perform CRUD operations and we aren't passing a null item
	if canCrud && postedType.ID > 0 {
		theCollection := mongoClient.Database("learnR").Collection("learnrsession") //Here's our collection
		collectedInterface := []interface{}{postedType}
		//Insert Our Data
		_, err2 := theCollection.InsertMany(theContext, collectedInterface)

		if err2 != nil {
			theErr := "Error adding LearnRSession in addLearnRSession in crudoperations API: " + err2.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			theErr := "LearnRSession successfully added in addlearnRSession in crudoperations: " + string(bs)
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 0
		}
	} else {
		theErr := "Error adding LearnRSession; could not perform CRUD or ID was bad: " + strconv.Itoa(postedType.ID)
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

//This deletes a LearnRSession to our database; called from anywhere
func deleteLearnRSession(w http.ResponseWriter, req *http.Request) {
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
		theErr := "Error reading the request from deleteSession: " + err.Error() + "\n" + string(bs)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}
	//Declare JSON we're looking for
	type LearnRSessionDelete struct {
		ID int `json:"ID"`
	}
	//Marshal it into our type
	var postedType LearnRSessionDelete
	json.Unmarshal(bs, &postedType)

	//Delete only if we had no issues above
	if canCrud && postedType.ID > 0 {
		//Search for User and delete
		collection := mongoClient.Database("learnR").Collection("learnrsession") //Here's our collection
		deletes := []bson.M{
			{"id": postedType.ID},
		} //Here's our filter to look for
		deletes = append(deletes, bson.M{"id": bson.M{
			"$eq": postedType.ID,
		}}, bson.M{"id": bson.M{
			"$eq": postedType.ID,
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
			theErr := "Error writing delete learnrSession in deleteLearnRSession in crudoperations: " + err.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			//Check to see if delete count worked; must have deleted at least one
			resultInt := bulkWrite.DeletedCount
			if resultInt > 0 {
				theErr := "LearnRSession successfully deleted in deletelearnRSession in crudoperations: " + string(bs)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 0
			} else {
				theErr := "No documents deleted for this given learnRSesion: " + strconv.Itoa(postedType.ID)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			}
		}
	} else {
		theErr := "Error, could not CRUD operate in deleteLearnRSession, or the number we recieved was wrong: " +
			strconv.Itoa(postedType.ID)
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

//This updates a LearnRSession to our database; called from anywhere
func updateLearnRSession(w http.ResponseWriter, req *http.Request) {
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
		theErr := "Error reading the request from updateSession: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Marshal it into our type
	var theTypePosted LearnRSession
	json.Unmarshal(bs, &theTypePosted)

	//Update item if we have successfully decoded from JSON
	if canCrud {
		//Update LearnR if their LearnR ID != 0 or nil
		if theTypePosted.ID != 0 {
			//Update User
			theTimeNow := time.Now()
			collection := mongoClient.Database("learnR").Collection("learnrsession") //Here's our collection
			theFilter := bson.M{
				"id": bson.M{
					"$eq": theTypePosted.ID, // check if bool field has value of 'false'
				},
			}
			updatedDocument := bson.M{
				"$set": bson.M{
					"id":               theTypePosted.ID,
					"learnrid":         theTypePosted.LearnRID,
					"learnrname":       theTypePosted.LearnRName,
					"thelearnr":        theTypePosted.TheLearnR,
					"theuser":          theTypePosted.TheUser,
					"targetusernumber": theTypePosted.TargetUserNumber,
					"ongoing":          theTypePosted.Ongoing,
					"textssent":        theTypePosted.TextsSent,
					"userresponses":    theTypePosted.UserResponses,
					"datecreated":      theTypePosted.DateCreated,
					"dateupdated":      theTimeNow.Format("2006-01-02 15:04:05"),
				},
			}
			updateResult, err := collection.UpdateOne(theContext, theFilter, updatedDocument)

			if err != nil {
				theErr := "Error writing update LearnRSession in updatelearnRSession in crudoperations: " + err.Error()
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			} else {
				//Check to see if anything was updated; if not, return the error
				if updateResult.ModifiedCount < 1 {
					theErr := "No document updated with this ID: " + strconv.Itoa(theTypePosted.ID)
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 1
				} else {
					theErr := "LearnRSession successfully updated in updateLearnRSession in crudoperations: " + string(bs) + "\n"
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 0
				}
			}
		} else {
			theErr := "The LearnRSession ID was not found: " + strconv.Itoa(theTypePosted.ID)
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

//This gets a LearnRSession with a certain LearnR ID
func getLearnRSession(w http.ResponseWriter, req *http.Request) {
	canCrud := true
	//Declare data to return
	type ReturnMessage struct {
		TheErr          []string      `json:"TheErr"`
		ResultMsg       []string      `json:"ResultMsg"`
		SuccOrFail      int           `json:"SuccOrFail"`
		ReturnedSession LearnRSession `json:"ReturnedSession"`
	}
	theReturnMessage := ReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success

	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from getLearnRSession: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Decalre JSON we recieve
	type LearnRSessionID struct {
		ID int `json:"ID"`
	}

	//Marshal it into our type
	var typePosted LearnRSessionID
	json.Unmarshal(bs, &typePosted)

	//If we successfully decoded, (and the ID is not 0) we can get our item
	if canCrud && typePosted.ID > 0 {
		/* Find the LearnRInfo with the given ID */
		var itemReturned LearnRSession                                           //Initialize Item to be returned after Mongo query
		collection := mongoClient.Database("learnR").Collection("learnrsession") //Here's our collection
		theFilter := bson.M{
			"id": bson.M{
				"$eq": typePosted.ID, // check if bool field has value of 'false'
			},
		}
		findOptions := options.Find()
		find, err := collection.Find(theContext, theFilter, findOptions)
		theFind := 0 //A counter to track how many users we find
		if find.Err() != nil || err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				stringUserID := strconv.Itoa(typePosted.ID)
				returnedErr := "For " + stringUserID + ", no LearnRSession was returned: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedSession = LearnRSession{}
			} else {
				stringUserID := strconv.Itoa(typePosted.ID)
				returnedErr := "For " + stringUserID + ", there was a Mongo Error: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedSession = LearnRSession{}
			}
		} else {
			//Found Learnorg, decode to return
			for find.Next(theContext) {
				stringid := strconv.Itoa(typePosted.ID)
				err := find.Decode(&itemReturned)
				if err != nil {
					returnedErr := "For " + stringid +
						", there was an error decoding document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedSession = LearnRSession{}
				} else if itemReturned.ID <= 1 {
					returnedErr := "For " + stringid +
						", there was an no document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedSession = LearnRSession{}
				} else {
					//Successful decode, do nothing
				}
				theFind = theFind + 1
			}
			find.Close(theContext)
		}

		if theFind <= 0 {
			//Error, return an error back and log it
			stringID := strconv.Itoa(typePosted.ID)
			returnedErr := "For " + stringID +
				", No LearnRInfo was returned."
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 1
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
			theReturnMessage.ReturnedSession = LearnRSession{}
		} else {
			//Success, log the success and return Item
			stringID := strconv.Itoa(typePosted.ID)
			returnedErr := "For " + stringID +
				", LearnrSession should be successfully decoded."
			//fmt.Println(returnedErr)
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 0
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ReturnedSession = itemReturned
		}
	} else {
		//Error, return an error back and log it
		theIDString := strconv.Itoa(typePosted.ID)
		returnedErr := "For " + theIDString +
			", No LearnRSession was returned. LearnSession was also not accepted: " + theIDString
		logWriter(returnedErr)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
		theReturnMessage.ReturnedSession = LearnRSession{}
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

/*ENDING LearnRSession CRUD OPERATIONS */

/* BEGINNING LearnRInforms CRUD OPERATIONS */

func addLearnRInforms(w http.ResponseWriter, req *http.Request) {
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
		theErr := "Error reading the request from learnRInforms: " + err.Error() + "\n" + string(bs)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false //Reading failed, need to return failure
	}
	//Marshal it into our type
	var postedType LearnRInforms
	json.Unmarshal(bs, &postedType)

	//Check to see if we can perform CRUD operations and we aren't passing a null item
	if canCrud && postedType.ID > 0 {
		theCollection := mongoClient.Database("learnR").Collection("learnrinforms") //Here's our collection
		collectedInterface := []interface{}{postedType}
		//Insert Our Data
		_, err2 := theCollection.InsertMany(theContext, collectedInterface)

		if err2 != nil {
			theErr := "Error adding LearnRInforms in addLearnRInforms in crudoperations API: " + err2.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			theErr := "LearnRInforms successfully added in addlearnRInforms in crudoperations: " + string(bs)
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 0
		}
	} else {
		theErr := "Error adding LearnRInforms; could not perform CRUD or ID was bad: " + strconv.Itoa(postedType.ID)
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

//This deletes a LearnRInforms to our database; called from anywhere
func deleteLearnRInforms(w http.ResponseWriter, req *http.Request) {
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
		theErr := "Error reading the request from deleteInforms: " + err.Error() + "\n" + string(bs)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}
	//Declare JSON we're looking for
	type LearnRInformsDelete struct {
		ID int `json:"ID"`
	}
	//Marshal it into our type
	var postedType LearnRInformsDelete
	json.Unmarshal(bs, &postedType)

	//Delete only if we had no issues above
	if canCrud && postedType.ID > 0 {
		//Search for User and delete
		collection := mongoClient.Database("learnR").Collection("learnrinforms") //Here's our collection
		deletes := []bson.M{
			{"id": postedType.ID},
		} //Here's our filter to look for
		deletes = append(deletes, bson.M{"id": bson.M{
			"$eq": postedType.ID,
		}}, bson.M{"id": bson.M{
			"$eq": postedType.ID,
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
			theErr := "Error writing delete learnrInforms in deleteLearnRInforms in crudoperations: " + err.Error()
			logWriter(theErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
			theReturnMessage.SuccOrFail = 1
		} else {
			//Check to see if delete count worked; must have deleted at least one
			resultInt := bulkWrite.DeletedCount
			if resultInt > 0 {
				theErr := "LearnRInforms successfully deleted in deletelearnRInforms in crudoperations: " + string(bs)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 0
			} else {
				theErr := "No documents deleted for this given learnRInforms: " + strconv.Itoa(postedType.ID)
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			}
		}
	} else {
		theErr := "Error, could not CRUD operate in deleteLearnRInforms, or the number we recieved was wrong: " +
			strconv.Itoa(postedType.ID)
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

//This updates a LearnRInforms to our database; called from anywhere
func updateLearnRInforms(w http.ResponseWriter, req *http.Request) {
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
		theErr := "Error reading the request from updateInforms: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Marshal it into our type
	var theTypePosted LearnRInforms
	json.Unmarshal(bs, &theTypePosted)

	//Update item if we have successfully decoded from JSON
	if canCrud {
		//Update LearnR if their LearnR ID != 0 or nil
		if theTypePosted.ID != 0 {
			//Update User
			theTimeNow := time.Now()
			collection := mongoClient.Database("learnR").Collection("learnrinforms") //Here's our collection
			theFilter := bson.M{
				"id": bson.M{
					"$eq": theTypePosted.ID, // check if bool field has value of 'false'
				},
			}
			updatedDocument := bson.M{
				"$set": bson.M{
					"id":          theTypePosted.ID,
					"name":        theTypePosted.Name,
					"learnrid":    theTypePosted.LearnRID,
					"learnrname":  theTypePosted.LearnRName,
					"order":       theTypePosted.Order,
					"theinfo":     theTypePosted.TheInfo,
					"shouldwait":  theTypePosted.ShouldWait,
					"waittime":    theTypePosted.WaitTime,
					"datecreated": theTypePosted.DateCreated,
					"dateupdated": theTimeNow.Format("2006-01-02 15:04:05"),
				},
			}
			updateResult, err := collection.UpdateOne(theContext, theFilter, updatedDocument)

			if err != nil {
				theErr := "Error writing update LearnRInform in updatelearnRInform in crudoperations: " + err.Error()
				logWriter(theErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
				theReturnMessage.SuccOrFail = 1
			} else {
				//Check to see if anything was updated; if not, return the error
				if updateResult.ModifiedCount < 1 {
					theErr := "No document updated with this ID: " + strconv.Itoa(theTypePosted.ID)
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 1
				} else {
					theErr := "LearnRInform successfully updated in updateLearnRInform in crudoperations: " + string(bs) + "\n"
					logWriter(theErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
					theReturnMessage.SuccOrFail = 0
				}
			}
		} else {
			theErr := "The LearnRInform ID was not found: " + strconv.Itoa(theTypePosted.ID)
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

//This gets a LearnRInforms with a certain ID
func getLearnRInforms(w http.ResponseWriter, req *http.Request) {
	canCrud := true
	//Declare data to return
	type ReturnMessage struct {
		TheErr               []string      `json:"TheErr"`
		ResultMsg            []string      `json:"ResultMsg"`
		SuccOrFail           int           `json:"SuccOrFail"`
		ReturnedLearnRInform LearnRInforms `json:"ReturnedLearnRInform"`
	}
	theReturnMessage := ReturnMessage{}
	theReturnMessage.SuccOrFail = 0 //Initially set to success

	//Unwrap from JSON
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		theErr := "Error reading the request from getLearnRInform: " + err.Error() + "\n" + string(bs)
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, theErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, theErr)
		theReturnMessage.SuccOrFail = 1
		logWriter(theErr)
		fmt.Println(theErr)
		canCrud = false
	}

	//Decalre JSON we recieve
	type LearnRInformID struct {
		ID int `json:"ID"`
	}

	//Marshal it into our type
	var typePosted LearnRInformID
	json.Unmarshal(bs, &typePosted)

	//If we successfully decoded, (and the ID is not 0) we can get our item
	if canCrud && typePosted.ID > 0 {
		/* Find the LearnRInfo with the given ID */
		var itemReturned LearnRInforms                                           //Initialize Item to be returned after Mongo query
		collection := mongoClient.Database("learnR").Collection("learnrinforms") //Here's our collection
		theFilter := bson.M{
			"id": bson.M{
				"$eq": typePosted.ID, // check if bool field has value of 'false'
			},
		}
		findOptions := options.Find()
		find, err := collection.Find(theContext, theFilter, findOptions)
		theFind := 0 //A counter to track how many users we find
		if find.Err() != nil || err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				stringUserID := strconv.Itoa(typePosted.ID)
				returnedErr := "For " + stringUserID + ", no LearnRInform was returned: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedLearnRInform = LearnRInforms{}
			} else {
				stringUserID := strconv.Itoa(typePosted.ID)
				returnedErr := "For " + stringUserID + ", there was a Mongo Error: " + err.Error()
				fmt.Println(returnedErr)
				logWriter(returnedErr)
				theReturnMessage.SuccOrFail = 1
				theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
				theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
				theReturnMessage.ReturnedLearnRInform = LearnRInforms{}
			}
		} else {
			//Found Item, decode to return
			for find.Next(theContext) {
				stringid := strconv.Itoa(typePosted.ID)
				err := find.Decode(&itemReturned)
				if err != nil {
					returnedErr := "For " + stringid +
						", there was an error decoding document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedLearnRInform = LearnRInforms{}
				} else if itemReturned.ID <= 1 {
					returnedErr := "For " + stringid +
						", there was an no document from Mongo: " + err.Error()
					fmt.Println(returnedErr)
					logWriter(returnedErr)
					theReturnMessage.SuccOrFail = 1
					theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
					theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
					theReturnMessage.ReturnedLearnRInform = LearnRInforms{}
				} else {
					//Successful decode, do nothing
				}
				theFind = theFind + 1
			}
			find.Close(theContext)
		}

		if theFind <= 0 {
			//Error, return an error back and log it
			stringID := strconv.Itoa(typePosted.ID)
			returnedErr := "For " + stringID +
				", No LearnRInform was returned."
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 1
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
			theReturnMessage.ReturnedLearnRInform = LearnRInforms{}
		} else {
			//Success, log the success and return Item
			stringID := strconv.Itoa(typePosted.ID)
			returnedErr := "For " + stringID +
				", LearnrInform should be successfully decoded."
			//fmt.Println(returnedErr)
			logWriter(returnedErr)
			theReturnMessage.SuccOrFail = 0
			theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
			theReturnMessage.TheErr = append(theReturnMessage.TheErr, "")
			theReturnMessage.ReturnedLearnRInform = itemReturned
		}
	} else {
		//Error, return an error back and log it
		theIDString := strconv.Itoa(typePosted.ID)
		returnedErr := "For " + theIDString +
			", No LearnRInform was returned. LearnInform was also not accepted: " + theIDString
		logWriter(returnedErr)
		theReturnMessage.SuccOrFail = 1
		theReturnMessage.ResultMsg = append(theReturnMessage.ResultMsg, returnedErr)
		theReturnMessage.TheErr = append(theReturnMessage.TheErr, returnedErr)
		theReturnMessage.ReturnedLearnRInform = LearnRInforms{}
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

/* ENDING LearnRInforms CRUD OPERATIONS */
