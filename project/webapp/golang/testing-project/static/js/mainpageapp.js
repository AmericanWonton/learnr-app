var app = angular.module('mymainpageApp', []);

var displayedTexts = [];

/* This takes the learnr array we've created and begins to list it on our page.
Divs will be created, being added into 'learnrHolderDiv'*/
function addlearnRVisuals(learnrArray){
    console.log("DEBUG: Getting learnr Visuals added.");
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
        for (var j = 0; j < learnrArray[n].Description.length; j++){
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


        console.log("DEBUG: About to add some infotext to our allTextHolder: " +  learnrArray[n].LearnRInforms[0].TheInfo);
        //Loop thorough texts to add text divs/texts to the allTextHolder
        for (var k = 0; k < learnrArray[n].LearnrInforms.length; k++) {
            console.log("DEBUG: Here is this infotext: " + learnrArray[n].LearnRInforms[k].TheInfo);
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
        var textDropDownDiv = document.createElement("button");
        textDropDownDiv.setAttribute("id", "textDropDownDiv" + n.toString());
        textDropDownDiv.setAttribute("class", "interiorBigInfoDiv");
        textDropDownDiv.setAttribute("name", "textDropDownDiv" + n.toString());
        //textDropDownDiv.style.backgroundImage = 'url(static/images/svg/downarrow.svg)'; //Set image
        textDropDownDiv.innerHTML = "Click to see texts";
        //Add event listener for this button
        textDropDownDiv.addEventListener('click', function(){
            //Evaluate 'allTextHolder' to see if it's hidden
            if (allTextHolder.style.display === "none") {
                //textDropDownDiv.style.backgroundImage = 'url(static/images/svg/uparrow.svg)'; //Set Image
                allTextHolder.style.display = "flex";
            } else {
                //textDropDownDiv.style.backgroundImage = 'url(static/images/svg/downarrow.svg)'; //Set Image
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

        /* DEBUG PRINTING */
        console.log("DEBUG: Added our first value learnr to webpage: " + learnrArray[n].Name);
    }
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

