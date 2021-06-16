var learnrTags = 0; //This increases every time we click our buttons,(used mostly for naming)
let learnrTagStrings = new Map(); //This contains all of our tags

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
                    submitLearnROrg.disabled = true;
                } else if (item == 'TooLong'){
                    informtextPLearnOrg.innerHTML = 'LearnR Organization about section must be under 400 characters';
                    submitLearnROrg.disabled = true;
                } else if (item == 'ContainsLanguage'){
                    informtextPLearnOrg.innerHTML = 'This section contains innapropriate content; please contact our help center for more information.';
                    submitLearnROrg.disabled = true;
                } else if (item == 'okay') {
                    informtextPLearnOrg.innerHTML = '';
                    submitLearnROrg.disabled = false;
                } else {
                    informtextPLearnOrg.innerHTML = 'Error checking your LearnR Organiztion about section';
                    submitLearnROrg.disabled = true;
                }
            }
        });
        xhr.send(textareaTellMe.value);
    });

    /* Changes our LearnR object to the organization ID when dropdown menu is selected */
    learnrorgs.addEventListener('onchange', function(){
        theLearnR.OrgID = Number(learnrorgs.value); //Adjust orgID value
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