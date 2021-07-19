var app = angular.module('mymainpageApp', []);

var displayedTexts = [];

/* This takes the learnr array we've created and begins to list it on our page.
Divs will be created, being added into 'learnrHolderDiv'*/
function addlearnRVisuals(learnrArray){
    /* Loop through our array to create divs/other properties */
    for (var n = 0; n < learnrArray.length; n++) {
        visualCreator(n, learnrArray);
    }
}

/* This is called whenever User searches for a specific learnr;
we delete all the learnrs on a page and populate it with what
their search returns. */
function rePopulateLearnRs(learnrArray){
    //Get our variables we need declared
    var learnrHolderDiv = document.getElementById("learnrHolderDiv");
    //Delete all variables within
    learnrHolderDiv.innerHTML = "";
    addlearnRVisuals(learnrArray); //Repopulate with the new learnrs
}

/* This creates our divs and other learnR stuf for users to see. Called from addlearnRVisuals */
function visualCreator(intCurrently, learnrArray){
    /* Create an array of bools for our button; this will determine if we can keep it disabled
    or not. 'True' means disabled for phone num, introduction, then person name */
    var submitDisablers = [true, true, true];

    var theInt = Number(intCurrently);
    //Get our variables we need declared
    var learnrHolderDiv = document.getElementById("learnrHolderDiv");

    //Create general div to hold learnr. Parent ==> learnrHolderDiv
    var resultLearnrHolder = document.createElement("div");
    resultLearnrHolder.setAttribute("id", "resultLearnrHolder" + theInt.toString());
    resultLearnrHolder.setAttribute("class", "resultLearnrHolder");
    resultLearnrHolder.setAttribute("name", "resultLearnrHolder" + theInt.toString());

    //Create Div to hold information on the LearnR. Parent ==> resultLearnrHolder
    var infolearnrHolder = document.createElement("div");
    infolearnrHolder.setAttribute("id", "infolearnrHolder" + theInt.toString());
    infolearnrHolder.setAttribute("class", "infolearnrHolder");
    infolearnrHolder.setAttribute("name", "infolearnrHolder" + theInt.toString());

    //Create Div to hold Name information for LearnR. Parent ==> infolearnrHolder
    var nameHolder = document.createElement("div");
    nameHolder.setAttribute("id", "nameHolder" + theInt.toString());
    nameHolder.setAttribute("class", "aInfoDiv");
    nameHolder.setAttribute("name", "nameHolder" + theInt.toString());
    //Create P to go inside Div for Name. Parent ==> nameHolder
    var pName = document.createElement("p");
    pName.setAttribute("id", "pName" + theInt.toString());
    pName.setAttribute("class", "learnRField");
    pName.setAttribute("name", "pName" + theInt.toString());
    pName.innerHTML = "Name: " + learnrArray[theInt].Name;
    //Attach this to div
    nameHolder.appendChild(pName);


    //Create Div to hold Description information for LearnR. Parent ==> infolearnrHolder
    var descriptionHolder = document.createElement("div");
    descriptionHolder.setAttribute("id", "descriptionHolder" + theInt.toString());
    descriptionHolder.setAttribute("class", "aInfoDiv");
    descriptionHolder.setAttribute("name", "descriptionHolder" + theInt.toString());
    //Create P to go inside Div for Description. Parent ==> descriptionHolder
    var theString = ""; //Used to put into inner HTML
    //Get value for description
    for (var j = 0; j < learnrArray[theInt].Description.length; j++){
        theString = theString + learnrArray[theInt].Description[j];
    }
    var pDescription = document.createElement("p");
    pDescription.setAttribute("id", "pDescription" + theInt.toString());
    pDescription.setAttribute("class", "learnRField");
    pDescription.setAttribute("name", "pDescription" + theInt.toString());
    pDescription.innerHTML = "Description: " + theString;
    //Attach value to div
    descriptionHolder.appendChild(pDescription);
    
    /* Add first two elements to 'infolearnrHolder' */
    infolearnrHolder.appendChild(nameHolder);
    infolearnrHolder.appendChild(descriptionHolder);

    /* Create text display to add to infolearnrHolder. ==> infolearnrHolder */
    var textDecisionHolder = document.createElement("div");
    textDecisionHolder.setAttribute("id", "textDecisionHolder" + theInt.toString());
    textDecisionHolder.setAttribute("class", "aBigInfoDiv");
    textDecisionHolder.setAttribute("name", "textDecisionHolder" + theInt.toString());

    //Create a Div to send this LearnR for the User
    var userLearnRSender = document.createElement("div");
    userLearnRSender.setAttribute("id", "userLearnRSender" + theInt.toString());
    userLearnRSender.setAttribute("class", "aBigInfoDiv");
    userLearnRSender.setAttribute("name", "userLearnRSender" + theInt.toString());
    //Initially set to hidden; will be unhidden with 'textDropDownDiv'
    userLearnRSender.style.display = "none";

    //Add the inputs for the userLearnRSender div
    //Send User Name
    var theFieldDiv = document.createElement("div");
    theFieldDiv.setAttribute("id", "theFieldDiv" + theInt.toString() + "1");
    theFieldDiv.setAttribute("class", "aBigInfoDiv");
    theFieldDiv.setAttribute("name", "theFieldDiv" + theInt.toString() + "1");
    //The Desc
    var fieldsideDiv = document.createElement("div");
    fieldsideDiv.setAttribute("id", "fieldsideDiv" + theInt.toString() + "1");
    fieldsideDiv.setAttribute("class", "fieldsideDiv");
    fieldsideDiv.setAttribute("name", "fieldsideDiv" + theInt.toString() + "1");
    var fieldsideDescP = document.createElement("p");
    fieldsideDescP.setAttribute("id", "fieldsideDescP" + theInt.toString() + "1");
    fieldsideDescP.setAttribute("class", "fieldP");
    fieldsideDescP.setAttribute("name", "fieldsideDescP" + theInt.toString() + "1");
    fieldsideDescP.innerHTML = "Enter the name of the person you want to send this to...";
    //Append the values
    fieldsideDiv.appendChild(fieldsideDescP);
    theFieldDiv.appendChild(fieldsideDiv);
    //The Input
    var fieldsideDiv = document.createElement("div");
    fieldsideDiv.setAttribute("id", "fieldsideDiv" + theInt.toString() + "2");
    fieldsideDiv.setAttribute("class", "fieldsideDiv");
    fieldsideDiv.setAttribute("name", "fieldsideDiv" + theInt.toString() + "2");
    var fieldinputPersonName = document.createElement("input");
    fieldinputPersonName.setAttribute("id", "fieldinputPersonName" + theInt.toString() + "2");
    fieldinputPersonName.setAttribute("class", "fieldInput");
    fieldinputPersonName.setAttribute("name", "fieldinputPersonName" + theInt.toString() + "2");
    fieldinputPersonName.setAttribute("type", "text");
    fieldinputPersonName.setAttribute("maxlength", "20");
    fieldinputPersonName.setAttribute("placeholder", "What is this person's name?");
    //Append the values
    fieldsideDiv.appendChild(fieldinputPersonName);
    theFieldDiv.appendChild(fieldsideDiv);

    //Attach this field
    userLearnRSender.appendChild(theFieldDiv);

    //Send User PhoneNumber
    var theFieldDiv = document.createElement("div");
    theFieldDiv.setAttribute("id", "theFieldDiv" + theInt.toString() + "2");
    theFieldDiv.setAttribute("class", "aBigInfoDiv");
    theFieldDiv.setAttribute("name", "theFieldDiv" + theInt.toString() + "2");
    //The Desc
    var fieldsideDiv = document.createElement("div");
    fieldsideDiv.setAttribute("id", "fieldsideDiv" + theInt.toString() + "2");
    fieldsideDiv.setAttribute("class", "fieldsideDiv");
    fieldsideDiv.setAttribute("name", "fieldsideDiv" + theInt.toString() + "2");
    var fieldsideDescP = document.createElement("p");
    fieldsideDescP.setAttribute("id", "fieldsideDescP" + theInt.toString() + "2");
    fieldsideDescP.setAttribute("class", "fieldP");
    fieldsideDescP.setAttribute("name", "fieldsideDescP" + theInt.toString() + "2");
    fieldsideDescP.innerHTML = "Enter the phone number of this person, like so, (area code in front, no hyphens): '13783434567'"
    //Append the values
    fieldsideDiv.appendChild(fieldsideDescP);
    theFieldDiv.appendChild(fieldsideDiv);
    //The Input
    var fieldsideDiv = document.createElement("div");
    fieldsideDiv.setAttribute("id", "fieldsideDiv" + theInt.toString() + "3");
    fieldsideDiv.setAttribute("class", "fieldsideDiv");
    fieldsideDiv.setAttribute("name", "fieldsideDiv" + theInt.toString() + "3");
    var fieldinputPersonPN = document.createElement("input");
    fieldinputPersonPN.setAttribute("id", "fieldinputPersonPN" + theInt.toString() + "3");
    fieldinputPersonPN.setAttribute("class", "fieldInput");
    fieldinputPersonPN.setAttribute("name", "fieldinputPersonPN" + theInt.toString() + "3");
    fieldinputPersonPN.setAttribute("type", "number");
    fieldinputPersonPN.setAttribute("maxlength", "11");
    fieldinputPersonPN.setAttribute("minlength", "11");
    fieldinputPersonPN.setAttribute("placeholder", "E.g. 13459780123");
    //Append the values
    fieldsideDiv.appendChild(fieldinputPersonPN);
    theFieldDiv.appendChild(fieldsideDiv);

    //Attach this field
    userLearnRSender.appendChild(theFieldDiv);

    //Send Introduction for User
    var theFieldDiv = document.createElement("div");
    theFieldDiv.setAttribute("id", "theFieldDiv" + theInt.toString() + "4");
    theFieldDiv.setAttribute("class", "aBigInfoDiv");
    theFieldDiv.setAttribute("name", "theFieldDiv" + theInt.toString() + "4");
    //The Desc
    var fieldsideDiv = document.createElement("div");
    fieldsideDiv.setAttribute("id", "fieldsideDiv" + theInt.toString() + "4");
    fieldsideDiv.setAttribute("class", "fieldsideDiv");
    fieldsideDiv.setAttribute("name", "fieldsideDiv" + theInt.toString() + "4");
    var fieldsideDescP = document.createElement("p");
    fieldsideDescP.setAttribute("id", "fieldsideDescP" + theInt.toString() + "4");
    fieldsideDescP.setAttribute("class", "fieldP");
    fieldsideDescP.setAttribute("name", "fieldsideDescP" + theInt.toString() + "4");
    fieldsideDescP.innerHTML = "What would you like to say to this person? Remember to be kind, it's the best way to be persuasive!"
    //Append the values
    fieldsideDiv.appendChild(fieldsideDescP);
    theFieldDiv.appendChild(fieldsideDiv);
    //The Input
    var fieldsideDiv = document.createElement("div");
    fieldsideDiv.setAttribute("id", "fieldsideDiv" + theInt.toString() + "5");
    fieldsideDiv.setAttribute("class", "fieldsideDiv");
    fieldsideDiv.setAttribute("name", "fieldsideDiv" + theInt.toString() + "5");
    var fieldinputIntroduction = document.createElement("textarea");
    fieldinputIntroduction.setAttribute("id", "fieldinputIntroduction" + theInt.toString() + "5");
    fieldinputIntroduction.setAttribute("class", "fieldTextAreaInput");
    fieldinputIntroduction.setAttribute("name", "fieldinputIntroduction" + theInt.toString() + "5");
    fieldinputIntroduction.setAttribute("maxlength", "120");
    fieldinputIntroduction.setAttribute("minlength", "1");
    fieldinputIntroduction.setAttribute("placeholder", "What would you like to say to this person to let them know what they're reading?");
    //Append the values
    fieldsideDiv.appendChild(fieldinputIntroduction);
    theFieldDiv.appendChild(fieldsideDiv);

    //Attach this field
    userLearnRSender.appendChild(theFieldDiv);

    //Send LearnRButton
    var theFieldDiv = document.createElement("div");
    theFieldDiv.setAttribute("id", "theFieldDiv" + theInt.toString() + "5");
    theFieldDiv.setAttribute("class", "aBigInfoDiv");
    theFieldDiv.setAttribute("name", "theFieldDiv" + theInt.toString() + "5");
    //Add Result P
    var sendLearnRResult = document.createElement("p");
    sendLearnRResult.setAttribute("id", "sendLearnRResult" + theInt.toString() + "6");
    sendLearnRResult.setAttribute("class", "resultInput");
    sendLearnRResult.setAttribute("name", "sendLearnRResult" + theInt.toString() + "6");
    //Append Result P
    theFieldDiv.appendChild(sendLearnRResult);
    //Add Button
    var sendLearnRButton = document.createElement("button");
    sendLearnRButton.setAttribute("id", "sendLearnRButton" + theInt.toString() + "7");
    sendLearnRButton.setAttribute("class", "sendButton");
    sendLearnRButton.setAttribute("name", "sendLearnRButton" + theInt.toString() + "7");
    sendLearnRButton.innerHTML = "Send LearnR";
    sendLearnRButton.disabled = true; //Initially set as disabled
    sendLearnRButton.addEventListener('click', function(){
        submitDisablers[0] = true;
        submitDisablers[1] = true;
        submitDisablers[2] = true;
        sendLearnRButton.disabled = true; //Disable Button for no further sends
        var OurJSON = {
            TheUser: TheUser,
            TheLearnR: learnrArray[theInt],
            TheLearnRInfo: {},
            PersonName: String(fieldinputPersonName.value),
            PersonPhoneNum: String(fieldinputPersonPN.value.toString()),
            Introduction: String(fieldinputIntroduction.value)
        };
        /* Start Loading Bar */
        sendLearnRResult.innerHTML = "Sending LearnR, please wait..."
        //Send Ajax
        var jsonString = JSON.stringify(OurJSON); //Stringify Data
        //Send Request to change page
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/canSendLearnR', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                var ReturnData = JSON.parse(item);
                if (ReturnData.SuccessNum == 0){
                    /* Successful LearnR Send start. Update our result, delay, then reload page */
                    submitDisablers[0] = true;
                    submitDisablers[1] = true;
                    submitDisablers[2] = true;
                    sendLearnRButton.disabled = true; //Disable Button for no further sends
                    sendLearnRResult.innerHTML = ReturnData.Message; //Set success message
                    setTimeout(() => { navigateHeader(3); }, 5000); //Delay 5, then reload page
                } else {
                    /* Sending text to User unsuccessful. Inform User */
                    sendLearnRButton.disabled = true;
                    sendLearnRResult.innerHTML = "Failed to send LearnR! " + ReturnData.Message;
                    setTimeout(() => { sendButtonFix(); }, 5000); //Delay 5, then un-disable button
                }
                function sendButtonFix(){
                    sendLearnRButton.disabled = false;
                    submitDisablers[0] = false;
                    submitDisablers[1] = false;
                    submitDisablers[2] = false;
                }
            }
        });
        xhr.send(jsonString);
    });
    //Append Button
    theFieldDiv.appendChild(sendLearnRButton);

    //Attach this field
    userLearnRSender.appendChild(theFieldDiv);

    /* Add event listeners that will disable our button above if they have the wrong in put */
    //Phone Number
    fieldinputPersonPN.addEventListener('input', function(){
        var theText = fieldinputPersonPN.value.toString();
        //Another check to see if numbers are too long
        var theText = fieldinputPersonPN.value.toString();
        if (theText.length > 11 || theText.length < 1) {
            submitDisablers[0] = true;
            sendLearnRButton.disabled = true;
        } else {
            submitDisablers[0] = false;
            if (submitDisablers[1] == false && submitDisablers[2] == false){
                sendLearnRButton.disabled = false;
            } else {
                sendLearnRButton.disabled = true;
            }
        }
        //Check for illegal characters
        if (theText.includes("-") || theText.includes("+") || theText.includes(" ") || theText.includes(".") || theText.includes(",")) {
            console.log("Removing bad character.");
            theText = theText.replace('-', '');
            theText = theText.replace('+', '');
            theText = theText.replace(' ', '');
            theText = theText.replace('.', '');
            theText = theText.replace(',', '');
            fieldinputPersonPN.value = Number(theText);
        }
    });
    //Person Name
    fieldinputPersonName.addEventListener('input', function(){
        if (fieldinputPersonName.value.length >= 1 && fieldinputPersonName.value.length <= 20){
            submitDisablers[1] = false;
            if (submitDisablers[0] == false && submitDisablers[2] == false){
                sendLearnRButton.disabled = false;
            } else {
                sendLearnRButton.disabled = true;
            }
        } else {
            submitDisablers[1] = true;
            sendLearnRButton.disabled = true;
        }
    });
    //Person Introduction
    fieldinputIntroduction.addEventListener('input', function(){
        if (fieldinputIntroduction.value.length >= 1 && fieldinputIntroduction.value.length <= 120){
            submitDisablers[2] = false;
            if (submitDisablers[0] == false && submitDisablers[1] == false){
                sendLearnRButton.disabled = false;
            } else {
                sendLearnRButton.disabled = true;
            }
        } else {
            submitDisablers[2] = true;
            sendLearnRButton.disabled = true;
        }
    });


    //Add the userLearnRSender to this hidden div
    textDecisionHolder.appendChild(userLearnRSender);

    //Make div to hold all texts for this LearnR,(will start as hidden). Parent ==> textDecisionHolder
    var allTextHolder = document.createElement("div");
    allTextHolder.setAttribute("id", "allTextHolder" + theInt.toString());
    allTextHolder.setAttribute("class", "aBigInfoDiv");
    allTextHolder.setAttribute("name", "allTextHolder" + theInt.toString());
    //Initially set to hidden; will be unhidden with 'textDropDownDiv'
    allTextHolder.style.display = "none";

    //Got texts, add allTextHolder to textDecisionHolder
    textDecisionHolder.appendChild(allTextHolder);

    //Loop thorough texts to add text divs/texts to the allTextHolder
    for (var k = 0; k < learnrArray[theInt].LearnRInforms.length; k++) {
        //Create Div to hold text. Parent ==> allTextHolder
        var aTextHolder = document.createElement("div");
        aTextHolder.setAttribute("id", "aTextHolder" + theInt.toString() + k.toString());
        aTextHolder.setAttribute("class", "textHolder");
        aTextHolder.setAttribute("name", "aTextHolder" + theInt.toString() + k.toString());

        //Create P with text in it. Parent ==> aTextHolder
        var aText = document.createElement("p");
        aText.setAttribute("id", "aText" + theInt.toString() + k.toString());
        aText.setAttribute("class", "textFont");
        aText.setAttribute("name", "aText" + theInt.toString() + k.toString());
        aText.innerHTML = learnrArray[theInt].LearnRInforms[k].TheInfo;
        
        //Add text to div
        aTextHolder.appendChild(aText);
        //Add to allTextHolder
        allTextHolder.appendChild(aTextHolder);
    }

    //Add div for drop down. Parent ==> textDecisionHolder
    var textDropDownDiv = document.createElement("button");
    textDropDownDiv.setAttribute("id", "textDropDownDiv" + theInt.toString());
    textDropDownDiv.setAttribute("class", "interiorBigInfoDiv");
    textDropDownDiv.setAttribute("name", "textDropDownDiv" + theInt.toString());
    //textDropDownDiv.style.backgroundImage = 'url(static/images/svg/downarrow.svg)'; //Set image
    textDropDownDiv.innerHTML = "Click to see texts";
    
    //Add this button to div first
    textDecisionHolder.appendChild(textDropDownDiv); 

    /* textDecisionHolder assembled, add it to infolearnrHolder */
    infolearnrHolder.appendChild(textDecisionHolder);

    /* All infolearnrHolder parts assembled. Add it to 'resultLearnrHolder' */
    resultLearnrHolder.appendChild(infolearnrHolder);

    /* All elements have been added to the learnr. Add to learnrHolderDiv */
    learnrHolderDiv.appendChild(resultLearnrHolder);

    //Add event listener for this button
    textDropDownDiv.addEventListener('click', function(){ 
        if (allTextHolder.style.display === "none"){
            //textDropDownDiv.style.backgroundImage = 'url(static/images/svg/uparrow.svg)'; //Set Image
            allTextHolder.style.display = "flex";
            userLearnRSender.style.display = "flex";
            //console.log("DEBUG: Showing this 'allTExtHolder': " + allTextHolder.getAttribute("id"));
        } else {
            //textDropDownDiv.style.backgroundImage = 'url(static/images/svg/downarrow.svg)'; //Set Image
            allTextHolder.style.display = "none";
            userLearnRSender.style.display = "none";
        }
    });

    
    /* DEBUG PRINTING */
}

//Set a custom delimiter for templates
app.config(function($interpolateProvider) {
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
});

//Main Controller
app.controller('myCtrl', function($scope, $timeout) {
    /* Use for casedata loading */
    $scope.caseData = null;
    //Learnr/LearnRInforms Declarations
    $scope.jsLearnRArray = [];
    $scope.jsLearnInformArray = [];
    /* LearnrSet
    Calls Ajax to get our Learnrs and put them into our jsLearnRArray */
    $scope.learnRSet = function() {
        var xhr = new XMLHttpRequest();
        xhr.open('GET', '/giveAllLearnrDisplay', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                var SuccessMSG = JSON.parse(item);
                if (SuccessMSG.SuccessNum === 0){
                    $scope.jsLearnRArray = SuccessMSG.TheDisplayLearnrs;
                    //console.log("DEBUG: Here is our jsLearnRArray: " + JSON.stringify($scope.jsLearnRArray));
                    $scope.caseData = 'hey!';
                    //Pass it on to Javascript to add data
                    addlearnRVisuals($scope.jsLearnRArray);
                } else {
                    console.log("Failed to reach out to giveAllLearnrDisplay");
                }
            }
        });
        xhr.send("testsend");
    };
    $scope.learnRSet();
    //mimic a delay in getting the data from $http
    $timeout(function () {
        $scope.caseData = 'hey!';
    }, 1000);
});

//Javascript stuff to call Angular and vice versa
window.addEventListener('DOMContentLoaded', function(){
    
});

//Used to control the search for LearnRs
function learnRSearch(){
    var learnRNameInput = document.getElementById("learnRNameInput");
    var learnRTagInput = document.getElementById("learnRTagInput");
    var resultThing = document.getElementById("resultThing");

    
    var SearchJSON = {
        TheNameInput: String(learnRNameInput.value),
        TheTagInput: String(learnRTagInput.value)
    };
    
    //console.log("DEBUG: Here is special cases: " + TheSpecialCases);
    //Send Ajax
    var jsonString = JSON.stringify(SearchJSON); //Stringify Data
    //Send Request to change page
    
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/searchLearnRs', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            var ReturnData = JSON.parse(item);
            if (ReturnData.SuccessNum == 0){
                learnRNameInput.value = "";
                learnRTagInput.value = "";
                //Take action if nothing is returned
                if (ReturnData.ReturnLearnRs != null && ReturnData.ReturnLearnRs){
                    //Repopulate learnrs
                    rePopulateLearnRs(ReturnData.ReturnLearnRs);
                } else {
                    //Nothing returned
                    resultThing.innerHTML = "No LearnRs returned from search!";
                }
            } else {
                resultThing.innerHTML = "Error finding those LearnRs! " + ReturnData.Message;
                learnRNameInput.value = "";
                learnRTagInput.value = "";
            }
        }
    });
    xhr.send(jsonString);
}

