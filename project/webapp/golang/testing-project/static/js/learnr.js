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
    var informtextPLearnOrg = document.getElementById("informtextPLearnOrg");
    var submitLearnR = document.getElementById("submitLearnR");

    /* Intitially set this to disabled so User needs to input values
    for this LearnR */
    submitLearnR.disabled = true;

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
                    var inputLearnRInfo = document.getElementById("inputLearnRInfo");
                    inputLearnRInfo.innerHTML = '';
                    submitLearnR.disabled = false;
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
                    informtextPLearnOrg.innerHTML = 'Please tell us why you want to make this LearnR Organization';
                    submitLearnR.disabled = true;
                } else if (item == 'TooLong'){
                    informtextPLearnOrg.innerHTML = 'LearnR Organization about section must be under 400 characters';
                    submitLearnR.disabled = true;
                } else if (item == 'ContainsLanguage'){
                    informtextPLearnOrg.innerHTML = 'This section contains innapropriate content; please contact our help center for more information.';
                    submitLearnR.disabled = true;
                } else if (item == 'okay') {
                    informtextPLearnOrg.innerHTML = '';
                    submitLearnR.disabled = false;
                } else {
                    informtextPLearnOrg.innerHTML = 'Error checking your LearnR Organiztion about section';
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
        console.log("Here is our key: " + key + " and here is our value: " + value);
    }
    learnrTags = learnrTags + 1; //Needed to interact with our map and other values
}

/* Called when '' is clicked. Create a new inform and add it to our 
learnrInforms. These can also be deleted at any time */
function addInform(){
    /* Declare our variables for later use */
    var textDesc = document.getElementById("textDesc");
    var timewait = document.getElementById("timewait");
    var timewaiting = document.getElementById("timewaiting");
    var resultHolderInformDiv = document.getElementById("resultHolderInformDiv");

    /* Create our LearnrInform object */
    var learnrInform = {
        ID: 0,
        Name: "",
        LearnRID: "",
        LearnRName: "",
        Order: 0,
        TheInfo: String(textDesc.value),
        ShouldWait: false,
        WaitTime: 0,
        DateCreated: "",
        DateUpdated: ""
    };
    /* Add time wait if it is selected */
    if (timewait.value != "false") {
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
    learnrInforms.set(Number(learnrInformsCount), String(textDesc.value));
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