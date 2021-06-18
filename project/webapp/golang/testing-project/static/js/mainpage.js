var learnrArray = [];

var learnrAssemble = {
    ID: 0,
    InfoID: 0,
    OrgID: 0,
    Name: "",
    Tags: [],
    Description: [],
    PhoneNums: [],
    LearnrInforms: [],
    Active: true,
    DateCreated: "",
    DateUpdated: ""
}; //Used to assemble new Learnrs through templates

var learnrInforms = {
    ID: 0,
    Name: "",
    LearnRID: 0,
    LearnRName: "",
    Order: 0,
    TheInfo: "",
    ShouldWait: true,
    WaitTime: 0,
    DateCreated: "",
    DateUpdated: ""
}; //Used for assembling LearnrInforms to add to Learnrs

/* Add the learnr to our array once it's assembled */
function sendLearnR(){
    learnrArray.push(learnrAssemble);
}

/* Debug function */
function debugLearnR(){
    for (var n = 0; n < learnrArray.length; n++) {
        console.log("DEBUG: Here is this spot in the learnr array: " + learnrArray[n].Name)
    }
}

/* Add LearnR values */
function setid(thevalue){
    learnrAssemble.ID = Number(thevalue);
}

function setinfoid(thevalue){
    learnrAssemble.InfoID = Number(thevalue);
}

function setorgid(thevalue){
    learnrAssemble.OrgID = Number(thevalue);
}

function setlearnrname(thevalue){
    learnrAssemble.Name = String(thevalue);
}

function setlearnractive(thevalue){
    learnrAssemble.Active = Boolean(thevalue);
}

function setlearnrdatecreated(thevalue){
    learnrAssemble.DateCreated = String(thevalue);
}

function setlearnrupdated(thevalue){
    learnrAssemble.DateUpdated = String(thevalue);
}

function settags(thevalue){
    learnrAssemble.Tags = String(thevalue);
}

function setdescription(thevalue){
    learnrAssemble.Description = String(thevalue);
}

function setlearnrphonenums(thevalue){
    learnrAssemble.PhoneNums = String(thevalue);
}


/* Add LearnRInform values */

function sendlearnrinformid(thevalue){
    learnrInforms.ID = Number(thevalue);
}

function setlearnrinformname(thevalue){
    learnrInforms.Name = String(thevalue);
}

function setlearnrnameinform(thevalue){
    learnrInforms.LearnRName = String(thevalue);
}

function setlearnrinformorder(thevalue){
    learnrInforms.Order = Number(thevalue);
}

function setlearnrinfo(thevalue){
    learnrInforms.TheInfo = String(thevalue);
}

function setshouldwaitinform(thevalue){
    learnrInforms.ShouldWait = Boolean(thevalue);
}

function setwaittimeinform(thevalue){
    learnrInforms.WaitTime = Number(thevalue);
}

function setdatecreatedinform(thevalue){
    learnrInforms.DateCreated = String(thevalue);
}

function setdateupdatedinform(thevalue){
    learnrInforms.DateUpdated = String(thevalue);
}

/* Send the created LearnRInform into the learnrAssemble Array 
with LearnrInforms inside it */
function sendLearnRInform(){
    learnrAssemble.LearnrInforms.push(learnrInforms);
}