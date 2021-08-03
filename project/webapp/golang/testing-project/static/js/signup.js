//Add our listening events to the window loading
window.addEventListener('DOMContentLoaded', function(){
    /* When 'Sign Up' is clicked contact Ajax to create profile for User; then
    we can log them in with new 'User' cookie created */
    var signUpB = document.getElementById("submitSignUpButton");
    var username = document.getElementById("username");
    var firstname = document.getElementById("firstname");
    var lastname = document.getElementById("lastname");
    var password = document.getElementById("password");
    var passwordRetype = document.getElementById("passwordRetype");
    var primaryPhoneNums = document.getElementById("primaryPhoneNums");
    var textareaTellMe = document.getElementById("textareaTellMe");
    var email = document.getElementById("email");
    var emailOkay = document.getElementById("emailOkay");
    var usernameErr = document.getElementById("form-input-info");
    var passwordErr = document.getElementById("form-input-info2");
    /* Used for informing User of the results */
    var informativeDivSignUp = document.getElementById("informativeDivSignUp");
    var informtextPSignUp = document.getElementById("informtextPSignUp");

    /* Check the database for Username when the key is pressed! */
    username.addEventListener('input', function(){
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/checkUsername', true);
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                if (item == 'TooShort') {
                    usernameErr.textContent = 'Please enter a Username';
                    signUpB.disabled = true;
                } else if (item == 'TooLong'){
                    usernameErr.textContent = 'Username must be under 20 characters';
                    signUpB.disabled = true;
                } else if (item == 'ContainsLanguage'){
                    usernameErr.textContent = 'Username is innapropriate';
                    signUpB.disabled = true;
                } else if (item == 'true') {
                    usernameErr.textContent = 'Username taken - Try another name!';
                    signUpB.disabled = true;
                } else {
                    //Check to see if this Username has the 'wrong characters'
                    var goodString = checkInput(username.value);
                    if (goodString === true){
                        //Username is good
                        usernameErr.textContent = '';
                        signUpB.disabled = false;
                    } else {
                        usernameErr.textContent = 'Username contains illegal characters... ';
                        signUpB.disabled = true;
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
            passwordErr.textContent = 'Please enter a password';
            signUpB.disabled = true;
        } else if (passString.length > 20){
            passwordErr.textContent = 'Password must be under 20 characters.';
            signUpB.disabled = true;
        } else {
            //Check to see if this Password has the 'wrong characters'
            var goodString = checkInput(password.value);
            if (goodString === true){
                //Password is good
                passwordErr.textContent = '';
                signUpB.disabled = false;
            } else {
                passwordErr.textContent = 'Password contains illegal characters... ';
                signUpB.disabled = true;
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
                    emailOkay.textContent = 'Please enter an email';
                    signUpB.disabled = true;
                } else if (item == 'TooLong'){
                    emailOkay.textContent = 'Email must be under 50 characters';
                    signUpB.disabled = true;
                } else if (item == 'ContainsLanguage'){
                    emailOkay.textContent = 'Email is innapropriate';
                    signUpB.disabled = true;
                } else if (item == 'true') {
                    emailOkay.textContent = 'Email taken - try another!';
                    signUpB.disabled = true;
                } else {
                    emailOkay.textContent = '';
                    signUpB.disabled = false;
                }
            }
        });
        xhr.send(String(email.value));
    });

    /* Check to see if the password matches the password re-type */
    passwordRetype.addEventListener('input', function(){
        if (passwordRetype.value != password.value){
            //Passwords don't match, inform user
            signUpB.disabled = true;
            passwordErr.textContent = 'Password must match password re-type';
        } else {
            //Check to see if this Password Re-type has the 'wrong characters'
            var goodString = checkInput(passwordRetype.value);
            if (goodString === true){
                //Password is good
                passwordErr.textContent = '';
                signUpB.disabled = false;
            } else {
                passwordErr.textContent = 'Password contains illegal characters... ';
                signUpB.disabled = true;
            }
        }
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
            YourNameInput: String(userName.value),
            YourEmailInput: String(email.value),
            YourUserID: Number(0),
            YourUser: newUser
        };

        var jsonString = JSON.stringify(MessageInfo);
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/sendVerificationEmail', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                var SuccessMSG = JSON.parse(item);
                if (SuccessMSG.SuccessNum === 0){
                    console.log("DEBUG: User successfully created: " + SuccessMSG.Message);
                    informtextPSignUp.innerHTML = "User succesfully created. Returning to login page...";
                    setTimeout(() => { navigateHeader(1); }, 4000);
                } else {
                    console.log("DEBUG: We have an error: " + SuccessMSG.SuccessNum + " " +
                    SuccessMSG.Message);
                    document.getElementById("informFormP").innerHTML = SuccessMSG.Message;
                    document.getElementById("informFormP").value = SuccessMSG.Message;
                    document.getElementById("informFormDiv").style.display = "block";
                }
            }
        });
        xhr.send(jsonString);
    });
});


/* This checks the verification code the User enters */
function checkCode(){
    /* Values for our verification code entry */
    var codeOkay = document.getElementById("codeOkay");
    var verifCode = document.getElementById("verifCode");
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
        Code: Number(verifCode.value)
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
                informtextPSignUp.innerHTML = "User succesfully created. Returning to login page...";
                setTimeout(() => { navigateHeader(1); }, 4000);
            } else {
                console.log("DEBUG: We have an error: " + SuccessMSG.SuccessNum + " " +
                    SuccessMSG.Message);
                document.getElementById("informFormP").innerHTML = SuccessMSG.Message;
                document.getElementById("informFormP").value = SuccessMSG.Message;
                document.getElementById("informFormDiv").style.display = "block";
                //Reload page after displaying error
                setTimeout(() => {  window.location.assign("/"); }, 3000);
            }
        }
    });
    xhr.send(jsonString);

}