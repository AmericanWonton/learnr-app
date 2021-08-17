var learnrTags = 0; //This increases every time we click our buttons,(used mostly for naming)
let learnrTagStrings = new Map(); //This contains all of our tags

var learnrInformsCount = 0; //This increases every time we click our buttons,(used mostly for naming)
let learnrInforms = new Map(); //This contains all our LearnRInforms

//Add our listener events to window loading
window.addEventListener('DOMContentLoaded', function(){
    var theLearnR = {
        ID: 0,
        InfoID: 0,
        OrgID: 0,
        Name: "",
        Tags: [],
        Description: [],
        PhoneNums: [],
        LearnRInforms: [],
        Active: true,
        DateCreated: "",
        DateUpdated: ""
    };

    /* When 'Submit' is clicked contact Ajax to create profile for LearnROrg;
    We can check to see if the LearnROrg they're making is taken or not */
    var learnrname = document.getElementById("learnrname");
    var textareaTellMe = document.getElementById("textareaTellMe");
    var learnrorgs = document.getElementById("learnrorgs");
    var informtextPLearnr = document.getElementById("informtextPLearnr");
    var informtextPLearnrDesc = document.getElementById("informtextPLearnrDesc");
    var submitLearnR = document.getElementById("submitLearnR");

    /* Intitially set this to disabled so User needs to input values
    for this LearnR */
    submitLearnR.disabled = true;
    /* Initially set the orgid, in case User dosen't change their selection */
    theLearnR.OrgID = Number(learnrorgs.value); //Adjust orgID value

    /* Check the database for LearnR Name when the key is pressed! */
    learnrname.addEventListener('input', function(){
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/checkLearnRNames', true);
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                if (item == 'TooShort') {
                    var inputLearnRInfo = document.getElementById("inputLearnRInfo");
                    inputLearnRInfo.innerHTML = 'Please enter a name for your LearnR';
                    submitLearnR.disabled = true;
                } else if (item == 'TooLong'){
                    var inputLearnRInfo = document.getElementById("inputLearnRInfo");
                    inputLearnRInfo.innerHTML = 'LearnR Name must be under 40 characters';
                    submitLearnR.disabled = true;
                } else if (item == 'ContainsLanguage'){
                    var inputLearnRInfo = document.getElementById("inputLearnRInfo");
                    inputLearnRInfo.innerHTML = 'This name contains innapropriate content; please contact our help center for more information.';
                    submitLearnR.disabled = true;
                } else if (item == 'true') {
                    var inputLearnRInfo = document.getElementById("inputLearnRInfo");
                    inputLearnRInfo.innerHTML = 'LearnR Name taken...try another name!';
                    submitLearnR.disabled = true;
                } else {
                    //Check to see if LearnRName is good
                    var inputLearnRInfo = document.getElementById("inputLearnRInfo");
                    inputLearnRInfo.innerHTML = '';
                    var goodString = checkInputLearnRName(learnrname.value);
                    if (goodString === true){
                        //learnRName is good
                        inputLearnRInfo.innerHTML = '';
                        submitLearnR.disabled = false;
                    } else {
                        inputLearnRInfo.innerHTML = 'LearnR Name contains illegal characters... ';
                        submitLearnR.disabled = true;
                    }
                }
            }
        });
        xhr.send(learnrname.value);
    });

    /* Check our data base to see if this about LearnR section is okay */
    textareaTellMe.addEventListener('input', function(){
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/checkOrgAbout', true);
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                if (item == 'TooShort') {
                    informtextPLearnrDesc.innerHTML = 'Please tell us why you want to make this LearnR Organization';
                    submitLearnR.disabled = true;
                } else if (item == 'TooLong'){
                    informtextPLearnrDesc.innerHTML = 'LearnR Organization about section must be under 400 characters';
                    submitLearnR.disabled = true;
                } else if (item == 'ContainsLanguage'){
                    informtextPLearnrDesc.innerHTML = 'This section contains innapropriate content; please contact our help center for more information.';
                    submitLearnR.disabled = true;
                } else if (item == 'okay') {
                    informtextPLearnrDesc.innerHTML = '';
                    submitLearnR.disabled = false;
                } else {
                    informtextPLearnrDesc.innerHTML = 'Error checking your LearnR Organiztion about section';
                    submitLearnR.disabled = true;
                }
            }
        });
        xhr.send(textareaTellMe.value);
    });

    /* Changes our LearnR object to the organization ID when dropdown menu is selected */
    learnrorgs.addEventListener('change', function(){
        theLearnR.OrgID = Number(learnrorgs.value); //Adjust orgID value
    });

    /* Disables the 'timewaiting' obect based on whether or not
    'timewait' is true or false */
    timewait.addEventListener('change', function(){
        var timewaiting = document.getElementById("timewaiting");

        if (timewait.value === "true") {
            timewaiting.disabled = false;
        } else {
            timewaiting.disabled = true;
        }
    });

    /* Submit this LearnR for review! */
    submitLearnR.addEventListener('click', function(){
        submitLearnR.disabled = true; //Disable this until Ajax comes back
        //Add all our LearnInforms to this array
        let learnInformArray = [];
        for (const [key, value] of learnrInforms.entries()){
            learnInformArray.push(value);
        }
        //Add our tags to an array as well
        let learnrTagArray = [];
        for (const [key, value] of learnrTagStrings.entries()){
            learnrTagArray.push(value);
        }
        //Add our new variables to our LearnR Array
        theLearnR.LearnRInforms = learnInformArray;
        theLearnR.Tags = learnrTagArray;
        theLearnR.Description.push(String(textareaTellMe.value));
        theLearnR.Name = String(learnrname.value);
        //Use Ajax to send this information
        //Declare Full JSON to send, with our UserID
        var SendJSON = {
            TheLearnr: theLearnR,
            OurUser: TheUser,
        };

        var jsonString = JSON.stringify(SendJSON);
        console.log("DEBUG: Here is our LearnR: " + JSON.stringify(theLearnR));
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/createLearnR', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                var SuccessMSG = JSON.parse(item);
                if (SuccessMSG.SuccessNum === 0){

                    informtextPLearnr.innerHTML = "LearnR succesfully created. Returning to mainpage...";
                    setTimeout(() => { navigateHeader(3); }, 4000);
                } else {
                    submitLearnR.disabled = false;
                    console.log("DEBUG: We have an error: " + SuccessMSG.SuccessNum + " " +
                    SuccessMSG.Message);
                    document.getElementById("informtextPLearnr").innerHTML = SuccessMSG.Message;
                    document.getElementById("informtextPLearnr").value = SuccessMSG.Message;
                }
            }
        });
        xhr.send(jsonString);
    });

});

/* Called when 'addLearnRTagButton' is clicked. Add the string to our
learnrTagStrings, then make a Div that will be appended to resultHolderTagDiv.
These can be deleted anytime. */
function addLearnRTag(){
    /* Declare variables */
    var tagDesc = document.getElementById("tagDesc");
    var resultHolderTagDiv = document.getElementById("resultHolderTagDiv");

    /* Create our divs to append to resultHolderTagDiv */
    var resultTagDiv = document.createElement("div");
    resultTagDiv.setAttribute("id", "resultTagDiv" + learnrTags.toString());
    resultTagDiv.setAttribute("class", "resultTagDiv");
    resultTagDiv.setAttribute("name", "resultTagDiv" + learnrTags.toString());
    
    var resultTagP = document.createElement("p");
    resultTagP.setAttribute("id", "resultTagP" + learnrTags.toString());
    resultTagP.setAttribute("class", "resultTagP");
    resultTagP.setAttribute("name", "resultTagP" + learnrTags.toString());
    resultTagP.innerHTML = String(tagDesc.value);

    /* Check to see if we can add this tag,(appropriate, good length, etc.) */
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/checkLearnRNames', true);
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            if (item == 'TooShort') {
                tagDesc.value = "";
                tagDesc.innerHTML = "";
                tagDesc.setAttribute("placeholder", "Enter something for your tag... at least 1 character!");
            } else if (item == 'TooLong'){
                tagDesc.value = "";
                tagDesc.innerHTML = "";
                tagDesc.setAttribute("placeholder", "LearnR Tag must be under 20 characters!");
            } else if (item == 'ContainsLanguage'){
                tagDesc.value = "";
                tagDesc.innerHTML = "";
                tagDesc.setAttribute("placeholder", "This LearnR Tag contains bad language; consult our team for more information.");
            } else if (item == 'true') {
                //Do Nothing, it's a tag
            } else {
                //Check to see if LearnRtag is good
                var goodString = checkInput(tagDesc.value);
                if (goodString === true){
                    //learnRtag is good
                    tagAdder();
                    tagDesc.value = "";
                    tagDesc.innerHTML = "";
                    tagDesc.setAttribute("placeholder", "What word describes this LearnR?");
                    //Good tag, add it to our display
                } else {
                    tagDesc.value = "";
                    tagDesc.innerHTML = "";
                    tagDesc.setAttribute("placeholder", "LearnR Tag contains illegal characters...");
                }
            }
        }
    });
    xhr.send(tagDesc.value);
    

    function tagAdder(){
        /* Add first elements to each other */
        resultTagDiv.appendChild(resultTagP);
        /* Add the tag to our current map */
        learnrTagStrings.set(Number(learnrTags), String(tagDesc.value));
        /* add the appropriate function on click */
        var thePosition = Number(learnrTags); //Used for deleteing tags
        resultTagDiv.addEventListener("click", function(){
            //Remove from this position
            learnrTagStrings.delete(thePosition);
            resultTagDiv.remove();
        });

        /* Add to result holder div for display */
        resultHolderTagDiv.appendChild(resultTagDiv);

        /* Clear the text box after entry */
        tagDesc.value = "";
        tagDesc.innerHTML = "";

        for (const [key, value] of learnrTagStrings.entries()){
            
        }
        learnrTags = learnrTags + 1; //Needed to interact with our map and other values
    }
}

/* Called when '' is clicked. Create a new inform and add it to our 
learnrInforms. These can also be deleted at any time */
function addInform(){
    /* Declare our variables for later use */
    var textDesc = document.getElementById("textDesc");
    var timewait = document.getElementById("timewait");
    var learnrname = document.getElementById("learnrname");
    var timewaiting = document.getElementById("timewaiting");
    var resultHolderInformDiv = document.getElementById("resultHolderInformDiv");

    /* Create our LearnrInform object */
    var learnrInform = {
        ID: 0,
        Name: "",
        LearnRID: 0,
        LearnRName: String(learnrname.value),
        Order: 0,
        TheInfo: String(textDesc.value),
        ShouldWait: false,
        WaitTime: 0,
        DateCreated: "",
        DateUpdated: ""
    };
    /* Add time wait if it is selected */
    if (timewait.value != "false") {
        console.log("DEBUG: We switched this learnrinform to true.");
        learnrInform.ShouldWait = true;
        learnrInform.WaitTime = Number(timewaiting.value);
    } 

    /* Create visual representation and add to our array */
    /* Create our divs to append to resultHolderTagDiv */
    var resultLearnrInformDiv = document.createElement("div");
    resultLearnrInformDiv.setAttribute("id", "resultLearnrInformDiv" + learnrInformsCount.toString());
    resultLearnrInformDiv.setAttribute("class", "resultLearnrInformDiv");
    resultLearnrInformDiv.setAttribute("name", "resultLearnrInformDiv" + learnrInformsCount.toString());
    
    var resultInformP = document.createElement("p");
    resultInformP.setAttribute("id", "resultInformP" + learnrInformsCount.toString());
    resultInformP.setAttribute("class", "resultInformP");
    resultInformP.setAttribute("name", "resultInformP" + learnrInformsCount.toString());
    resultInformP.innerHTML = String(textDesc.value).substr(0,20) + "..." + "//ShouldWait: " +
    timewait.value + "//TimeWait: " + timewaiting.value.toString();

    /* Add first elements to each other */
    resultLearnrInformDiv.appendChild(resultInformP);
    /* Add the LearnRInform to our current map */
    learnrInforms.set(Number(learnrInformsCount), learnrInform);
    /* add the appropriate function on click */
    var thePosition = Number(learnrInformsCount); //Used for deleteing tags
    resultLearnrInformDiv.addEventListener("click", function(){
        //Remove from this position
        learnrInforms.delete(thePosition);
        resultLearnrInformDiv.remove();
    });
    /* Add to result holder div for display */
    resultHolderInformDiv.appendChild(resultLearnrInformDiv);
    /* Clear the text box after entry */
    textDesc.value = "";
    textDesc.innerHTML = "";
    timewaiting.value = 0;

    learnrInformsCount = learnrInformsCount + 1; //Increment this for future naming values
}