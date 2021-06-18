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

/* This takes the learnr array we've created and begins to list it on our page.
Divs will be created, being added into 'learnrHolderDiv'*/
function addlearnRVisuals(){

    //Get our variables we need declared
    var learnrHolderDiv = document.getElementById("learnrHolderDiv");

    /* Loop through our array to create divs/other properties */
    for (var n = 0; n < learnrArray.length; n++) {
        //Create general div to hold learnr. Parent ==> learnrHolderDiv
        var resultLearnrHolder = document.createElement("div");
        resultLearnrHolder.setAttribute("id", "resultLearnrHolder" + n.toString());
        resultLearnrHolder.setAttribute("class", "resultLearnrHolder");
        resultLearnrHolder.setAttribute("name", "resultLearnrHolder" + n.toString());

        //Create Div to hold information on the LearnR. Parent ==> resultLearnrHolder
        var infolearnrHolder = document.createElement("div");
        infolearnrHolder.setAttribute("id", "infolearnrHolder" + n.toString());
        infolearnrHolder.setAttribute("class", "infolearnrHolder");
        infolearnrHolder.setAttribute("name", "infolearnrHolder" + n.toString());

        //Create Div to hold Name information for LearnR. Parent ==> infolearnrHolder
        var nameHolder = document.createElement("div");
        nameHolder.setAttribute("id", "nameHolder" + n.toString());
        nameHolder.setAttribute("class", "aInfoDiv");
        nameHolder.setAttribute("name", "nameHolder" + n.toString());
        //Create P to go inside Div for Name. Parent ==> nameHolder
        var pName = document.createElement("p");
        pName.setAttribute("id", "pName" + n.toString());
        pName.setAttribute("class", "learnRField");
        pName.setAttribute("name", "pName" + n.toString());
        pName.innerHTML = "Name: " + learnrArray[n].Name;
        //Attach this to div
        nameHolder.appendChild(pName);


        //Create Div to hold Description information for LearnR. Parent ==> infolearnrHolder
        var descriptionHolder = document.createElement("div");
        descriptionHolder.setAttribute("id", "descriptionHolder" + n.toString());
        descriptionHolder.setAttribute("class", "aInfoDiv");
        descriptionHolder.setAttribute("name", "descriptionHolder" + n.toString());
        //Create P to go inside Div for Description. Parent ==> descriptionHolder
        var theString = ""; //Used to put into inner HTML
        //Get value for description
        for (var j = 0; j < learnrArray[n].Description; j++){
            theString = theString + learnrArray[n].Description[j];
        }
        var pDescription = document.createElement("p");
        pDescription.setAttribute("id", "pDescription" + n.toString());
        pDescription.setAttribute("class", "learnRField");
        pDescription.setAttribute("name", "pDescription" + n.toString());
        pDescription.innerHTML = "Description: " + theString;
        //Attach value to div
        descriptionHolder.appendChild(pDescription);
        
        /* Add first two elements to 'infolearnrHolder' */
        infolearnrHolder.appendChild(nameHolder);
        infolearnrHolder.appendChild(descriptionHolder);

        /* Create text display to add to infolearnrHolder. ==> infolearnrHolder */
        var textDecisionHolder = document.createElement("div");
        textDecisionHolder.setAttribute("id", "textDecisionHolder" + n.toString());
        textDecisionHolder.setAttribute("class", "aBigInfoDiv");
        textDecisionHolder.setAttribute("name", "textDecisionHolder" + n.toString());

        //Make div to hold all texts for this LearnR,(will start as hidden). Parent ==> textDecisionHolder
        var allTextHolder = document.createElement("div");
        allTextHolder.setAttribute("id", "allTextHolder" + n.toString());
        allTextHolder.setAttribute("class", "aBigInfoDiv");
        allTextHolder.setAttribute("name", "allTextHolder" + n.toString());
        //Initially set to hidden; will be unhidden with 'textDropDownDiv'
        allTextHolder.style.display = "none";

        //Loop thorough texts to add text divs/texts to the allTextHolder
        for (var k = 0; k < learnrArray[n].LearnrInforms; k++) {
            console.log("DEBUG: Here is this infotext: " + learnrArray[n].LearnrInforms[k].TheInfo);
            //Create Div to hold text. Parent ==> allTextHolder
            var aTextHolder = document.createElement("div");
            aTextHolder.setAttribute("id", "aTextHolder" + n.toString() + k.toString());
            aTextHolder.setAttribute("class", "textHolder");
            aTextHolder.setAttribute("name", "aTextHolder" + n.toString() + k.toString());

            //Create P with text in it. Parent ==> aTextHolder
            var aText = document.createElement("p");
            aText.setAttribute("id", "aText" + n.toString() + k.toString());
            aText.setAttribute("class", "textFont");
            aText.setAttribute("name", "aText" + n.toString() + k.toString());
            aText.innerHTML = learnrArray[n].LearnrInforms[k].TheInfo;
            
            //Add text to div
            aTextHolder.appendChild(aText);
            //Add to allTextHolder
            allTextHolder.appendChild(aTextHolder);
        }

        //Add div for drop down. Parent ==> textDecisionHolder
        var textDropDownDiv = document.createElement("div");
        textDropDownDiv.setAttribute("id", "textDropDownDiv" + n.toString());
        textDropDownDiv.setAttribute("class", "interiorBigInfoDiv");
        textDropDownDiv.setAttribute("name", "textDropDownDiv" + n.toString());
        textDropDownDiv.style.backgroundImage = 'url(static/images/svg/downarrow.svg)'; //Set image
        //Add event listener for this button
        textDropDownDiv.addEventListener('click', function(){
            //Evaluate 'allTextHolder' to see if it's hidden
            if (allTextHolder.style.display === "none") {
                textDropDownDiv.style.backgroundImage = 'url(static/images/svg/uparrow.svg)'; //Set Image
                allTextHolder.style.display = "flex";
            } else {
                textDropDownDiv.style.backgroundImage = 'url(static/images/svg/downarrow.svg)'; //Set Image
                allTextHolder.style.display = "none";
            }
        });
        //Add this button to div first
        textDecisionHolder.appendChild(textDropDownDiv);

        //Got texts, add allTextHolder to textDecisionHolder
        textDecisionHolder.appendChild(allTextHolder);

        /* textDecisionHolder assembled, add it to infolearnrHolder */
        infolearnrHolder.appendChild(textDecisionHolder);

        /* All infolearnrHolder parts assembled. Add it to 'resultLearnrHolder' */
        resultLearnrHolder.appendChild(infolearnrHolder);

        /* All elements have been added to the learnr. Add to learnrHolderDiv */
        learnrHolderDiv.appendChild(resultLearnrHolder);
    }
}

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