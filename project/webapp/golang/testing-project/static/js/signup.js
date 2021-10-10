
/* The original height of our div that holds form values.
Will be adjusted as the dom is resized or errors appear/dissapear.
Should be changed if css files changes for this div: divformDivSignUp */
let originalHeight = 700;
/* Order: Username, Password, RetypePassword, Email
1 is not okay, 0 is okay*/
let okayFields = [0,0,0,0];
let errorMessages = ["", "", "", ""];

//Add our listening events to the window loading
window.addEventListener('DOMContentLoaded', function(){
    /* When 'Sign Up' is clicked contact Ajax to create profile for User; then
    we can log them in with new 'User' cookie created */

    /* When an error is displayed, we must increase the height of 'divformDivSignUp'
    to properly display the error; when the erro goes away, we can decrease it.*/
    var username = document.getElementById("username");
    var firstname = document.getElementById("firstname");
    var lastname = document.getElementById("lastname");
    var password = document.getElementById("password");
    var passwordRetype = document.getElementById("passwordRetype");
    var primaryPhoneNums = document.getElementById("primaryPhoneNums");
    var textareaTellMe = document.getElementById("textareaTellMe");
    var email = document.getElementById("email");
    var signUpB = document.getElementById("submitSignUpButton");
    var verificationDiv = document.getElementById("verificationDiv");
    /* Used for informing User of the results */
    var informtextPSignUp = document.getElementById("informtextPSignUp");

    /* Check the database for Username when the key is pressed! */
    username.addEventListener('input', function(){
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/checkUsername', true);
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                if (item == 'TooShort') {
                    errorDisplay("Please enter a Username", 0);
                    
                } else if (item == 'TooLong'){
                    errorDisplay("Username must be under 20 characters",0);
                    
                } else if (item == 'ContainsLanguage'){
                    errorDisplay("Username is innapropriate",0);
                    
                } else if (item == 'true') {
                    errorDisplay("Username is taken...try another name!",0);
                    
                } else {
                    //Check to see if this Username has the 'wrong characters'
                    var goodString = checkInput(username.value);
                    if (goodString === true){
                        //Username is good
                        errorDisplay("",0);
                        signUpB.disabled = false;
                    } else {
                        errorDisplay("Username contains illegal characters...",0);
                        
                    }
                }
            }
        });
        xhr.send(username.value);
    });

    /*Check to see if password is an appropriate length! */
    password.addEventListener('input', function(){
        passString = password.value;
        if (passString.length <= 0) {
            errorDisplay("Please enter a password",1);
            
        } else if (passString.length > 20){
            errorDisplay("Password must be under 20 characters",1);
            
        } else {
            //Check to see if this Password has the 'wrong characters'
            var goodString = checkInput(password.value);
            if (goodString === true){
                //Password is good
                errorDisplay("",1);
                signUpB.disabled = false;
            } else {
                errorDisplay("Password contains illegal characters",1);
                
            }
        }
    });

    /* Check to see if the password matches the password re-type */
    passwordRetype.addEventListener('input', function(){
        if (passwordRetype.value != password.value){
            //Passwords don't match, inform user
            errorDisplay("Passwords do not match!",2);
            
        } else {
            //Check to see if this Password Re-type has the 'wrong characters'
            var goodString = checkInput(passwordRetype.value);
            if (goodString === true){
                //Password is good
                errorDisplay("",2);
                signUpB.disabled = false;
            } else {
                errorDisplay("Password contains illegal characters",2);
                
            }
        }
    });

    /* Check to see if email is already in use */
    email.addEventListener('input', function(){
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/checkEmail', true);
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                if (item == 'TooShort') {
                    errorDisplay("Please enter an email",3);
                    
                } else if (item == 'TooLong'){
                    errorDisplay("Email address must be under 50 characters",3);
                    
                } else if (item == 'ContainsLanguage'){
                    errorDisplay("Email is innapropriate",3);
                    
                } else if (item == 'true') {
                    errorDisplay("Email is taken...try another!",3);
                    
                } else {
                    errorDisplay("",3);
                    signUpB.disabled = false;
                }
            }
        });
        xhr.send(String(email.value));
    });

    /* Send Email to User when fields are filled out...
    THIS IS ONLY FOR VERIFICATION */
    signUpB.addEventListener("click", function(){
        var newUser = {
            UserName: String(username.value),
            Password: String(password.value),
            Firstname: String(firstname.value),
            Lastname: String(lastname.value),
            PhoneNums: [String(primaryPhoneNums.value)],
            UserID:   0,
            Email: [String(email.value)],
            Whoare: String(textareaTellMe.value),
            AdminOrgs: new Array(),
            OrgMember: new Array(),
            Banned: false,
            DateCreated: "",
            DateUpdated: "",
        };

        var MessageInfo = {
            YourNameInput: String(username.value),
            YourEmailInput: String(email.value),
            YourUserID: Number(0),
            YourUser: newUser
        };

        /* Disable the inputs so they can't be changed */
        
        username.disabled = true;
        firstname.disabled = true;
        lastname.disabled = true;
        password.disabled = true;
        passwordRetype.disabled = true;
        primaryPhoneNums.disabled = true;
        textareaTellMe.disabled = true;
        email.disabled = true;

        var jsonString = JSON.stringify(MessageInfo);
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/sendVerificationEmail', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                var ReturnMessage = JSON.parse(item);
                if (ReturnMessage.SuccOrFail === 0){
                    console.log("DEBUG: Email verification sent: " + ReturnMessage.ResultMsg);
                    var promptAnswer = prompt("Please enter your email verification", "");
                    if (promptAnswer != null || promptAnswer != ""){
                        //User answered, evaulate the code
                        checkCode(promptAnswer);
                    } else {
                        //Did not detect User input, reset page
                        alert("User input not detected, please reload page and try again...");
                    }
                } else {
                    console.log("DEBUG: We have an error: " + SuccessMSG.SuccOrFail + " " +
                    SuccessMSG.TheErr);
                    document.getElementById("informFormP").innerHTML = SuccessMSG.TheErr;
                    document.getElementById("informFormP").value = SuccessMSG.TheErr;
                    document.getElementById("informFormDiv").style.display = "block";
                    setTimeout(() => { navigateHeader(1); }, 4000);
                }
            }
        });
        xhr.send(jsonString);
    });
});


/* This checks the verification code the User enters */
function checkCode(theCode){
    /* Values for our verification code entry */
    var codeOkay = document.getElementById("codeOkay");
    /* Values for building our User */
    var username = document.getElementById("username");
    var firstname = document.getElementById("firstname");
    var lastname = document.getElementById("lastname");
    var password = document.getElementById("password");
    var primaryPhoneNums = document.getElementById("primaryPhoneNums");
    var textareaTellMe = document.getElementById("textareaTellMe");
    var email = document.getElementById("email");

    var newUser = {
        UserName: String(username.value),
        Password: String(password.value),
        Firstname: String(firstname.value),
        Lastname: String(lastname.value),
        PhoneNums: [String(primaryPhoneNums.value)],
        UserID:   0,
        Email: [String(email.value)],
        Whoare: String(textareaTellMe.value),
        AdminOrgs: new Array(),
        OrgMember: new Array(),
        Banned: false,
        DateCreated: "",
        DateUpdated: "",
    };

    var NewCreation = {
        NewUser: newUser,
        Code: Number(theCode)
    };

    var jsonString = JSON.stringify(NewCreation);
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/createUser', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            var SuccessMSG = JSON.parse(item);
            if (SuccessMSG.SuccOrFail === 0){
                /* User succussfully created. Delay, then send to home page */
                errorSignUpDiv.style.backgroundColor = "green";
                signInErrorP.innerHTML = SuccessMSG.Message;
                signInErrorP.innerText = SuccessMSG.Message;
                setTimeout(() => { navigateHeader(1); }, 4000);
            } else {
                console.log("DEBUG: We have an error: " + SuccessMSG.SuccessNum + " " +
                SuccessMSG.Message);
                signInErrorP.innerText = SuccessMSG.Message;
                errorSignUpDiv.style.backgroundColor = "red";
                //Reload page after displaying error
                setTimeout(() => {  window.location.assign("/"); }, 3000);
            }
        }
    });
    xhr.send(jsonString);

}

/* This displays an error based on User input*/
function errorDisplay(theError, whichField) {
    
    var signUpB = document.getElementById("submitSignUpButton");
    var signInErrorP = document.getElementById("signInErrorP");
    var errorSignUpDiv = document.getElementById("errorSignUpDiv");

    if (theError == ""){
        /* Error removed from this field. Leave color blank, 
        then turn their field to okay. Errors may still display if there are more errors*/
        okayFields[whichField] = 0;
        errorMessages[whichField] = theError;
        if ((okayFields[0] == 1) || (okayFields[1] == 1) || (okayFields[2] == 1) || (okayFields[3] == 1)){
            errorSignUpDiv.style.backgroundColor = "red";
            signUpB.disabled = true;
            //Find the first error still occuring
            for (var x = 0; x < (okayFields.length); x++){
                if (okayFields[x] == 1) {
                    signInErrorP.innerHTML = errorMessages[x];
                    signInErrorP.innerText = errorMessages[x];
                    break;
                }
            }
        } else {
            signInErrorP.innerHTML = "";
            signInErrorP.innerText = "";
            errorSignUpDiv.style.backgroundColor = "white";
            signUpB.disabled = false;
        }
    } else {
        /* Change the background color, have the newewst error displayed */
        errorSignUpDiv.style.backgroundColor = "red";
        errorMessages[whichField] = theError;
        okayFields[whichField] = 1;
        signUpB.disabled = true;
        //Find the first error still occuring
        for (var x = 0; x < (okayFields.length); x++){
            if (okayFields[x] == 1) {
                signInErrorP.innerHTML = errorMessages[x];
                signInErrorP.innerText = errorMessages[x];
                break;
            }
        }
    }
}