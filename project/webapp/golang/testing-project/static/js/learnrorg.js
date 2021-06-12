


//Add our listener events to window loading
window.addEventListener('DOMContentLoaded', function(){
    /* When 'Submit' is clicked contact Ajax to create profile for LearnROrg;
    We can check to see if the LearnROrg they're making is taken or not */
    var signUpB = document.getElementById("signUpSubmitB");
    var learnrorgname = document.getElementById("learnrorgname");
    var textareaTellMe = document.getElementById("textareaTellMe");
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
        xhr.send(username.value);
    });

    /* Check our data base to see if this message is okay */
    textareaTellMe.addEventListener('input', function(){
        
    })

});