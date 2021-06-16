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