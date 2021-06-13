


//Add our listener events to window loading
window.addEventListener('DOMContentLoaded', function(){
    /* When 'Submit' is clicked contact Ajax to create profile for LearnROrg;
    We can check to see if the LearnROrg they're making is taken or not */
    var signUpB = document.getElementById("signUpSubmitB");
    var learnrorgname = document.getElementById("learnrorgname");
    var textareaTellMe = document.getElementById("textareaTellMe");
    var informtextPLearnOrg = document.getElementById("informtextPLearnOrg");
    var submitLearnROrg = document.getElementById("submitLearnROrg");

    /* Check the database for LearnRORg Name when the key is pressed! */
    learnrorgname.addEventListener('input', function(){
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/checkLearnROrgNames', true);
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                if (item == 'TooShort') {
                    var inputOrgInfo = document.getElementById("inputOrgInfo");
                    inputOrgInfo.innerHTML = 'Please enter a name for your LearnR Organization';
                    signUpB.disabled = true;
                } else if (item == 'TooLong'){
                    var inputOrgInfo = document.getElementById("inputOrgInfo");
                    inputOrgInfo.innerHTML = 'LearnR Organization Name must be under 25 characters';
                    signUpB.disabled = true;
                } else if (item == 'ContainsLanguage'){
                    var inputOrgInfo = document.getElementById("inputOrgInfo");
                    inputOrgInfo.innerHTML = 'This name contains innapropriate content; please contact our help center for more information.';
                    signUpB.disabled = true;
                } else if (item == 'true') {
                    var inputOrgInfo = document.getElementById("inputOrgInfo");
                    inputOrgInfo.innerHTML = 'LearnR Organization Name taken...try another name!';
                    signUpB.disabled = true;
                } else {
                    var inputOrgInfo = document.getElementById("inputOrgInfo");
                    inputOrgInfo.innerHTML = '';
                    signUpB.disabled = false;
                }
            }
        });
        xhr.send(learnrorgname.value);
    });

    /* Check our data base to see if this about LearnR Org section is okay */
    textareaTellMe.addEventListener('input', function(){
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/checkOrgAbout', true);
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                if (item == 'TooShort') {
                    informtextPLearnOrg.innerHTML = 'Please tell us why you want to make this LearnR Organization';
                    signUpB.disabled = true;
                } else if (item == 'TooLong'){
                    informtextPLearnOrg.innerHTML = 'LearnR Organization about section must be under 400 characters';
                    signUpB.disabled = true;
                } else if (item == 'ContainsLanguage'){
                    informtextPLearnOrg.innerHTML = 'This section contains innapropriate content; please contact our help center for more information.';
                    signUpB.disabled = true;
                } else if (item == 'okay') {
                    informtextPLearnOrg.innerHTML = '';
                    signUpB.disabled = false;
                } else {
                    informtextPLearnOrg.innerHTML = 'Error checking your LearnR Organiztion about section';
                    signUpB.disabled = true;
                }
            }
        });
        xhr.send(textareaTellMe.value);
    });

    submitLearnROrg.addEventListener('click', function(){
        var newLearnOrg = {
            OrgID: 0,
            Name: String(learnrorgname.value),
            OrgGoals: new Array(),
            UserList: new Array(),
            AdminList: new Array(),
            LearnrList: new Array(),
            DateCreated: "",
            DateUpdated: ""
        };
        newLearnOrg.OrgGoals.push(String(textareaTellMe.value));

        //Declare Full JSON to send, with our UserID
        var SendJSON = {
            TheLearnOrg: newLearnOrg,
            OurUser: TheUser,
        };

        var jsonString = JSON.stringify(SendJSON);
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/createLearnROrg', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                var SuccessMSG = JSON.parse(item);
                if (SuccessMSG.SuccessNum === 0){
                    informtextPLearnOrg.innerHTML = "LearnROrg succesfully created. Returning to mainpage...";
                    setTimeout(() => { navigateHeader(3); }, 4000);
                } else {
                    console.log("DEBUG: We have an error: " + SuccessMSG.SuccessNum + " " +
                    SuccessMSG.Message);
                    document.getElementById("informtextPLearnOrg").innerHTML = SuccessMSG.Message;
                    document.getElementById("informtextPLearnOrg").value = SuccessMSG.Message;
                }
            }
        });
        xhr.send(jsonString);



    });

});