/* When 'Sign Up' is clicked contact Ajax to create profile for User; then
we can log them in with new 'User' cookie created */
var signUpB = document.getElementById("signUpSubmitB");
var userName = document.getElementById("username");
var password = document.getElementById("password");
var passwordRetype = document.getElementById("passwordRetype");
var primaryPhoneNums = document.getElementById("primaryPhoneNums");
var textareaTellMe = document.getElementById("textareaTellMe");
var email = document.getElementById("email");
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
                usernameErr.textContent = '';
                signUpB.disabled = false;
            }
        }
    });
    xhr.send(userName.value);
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
        passwordErr.textContent = '';
        signUpB.disabled = false;
    }
});

/* Check to see if the password matches the password re-type */
passwordRetype.addEventListener('input', function(){
    if (passwordRetype.value != password.value){
        //Passwords don't match, inform user
        console.log("DEBUG: Passwords don't match: " + passwordRetype.value + " and " + password.value);
        signUpB.disabled = true;
        passwordErr.textContent = 'Password must match password re-type';
    } else {
        //Passwords match
        passwordErr.textContent = '';
        signUpB.disabled = false;
    }
});

/* Send Email to User when fields are filled out */
signUpB.addEventListener("click", function(){
    var newUser = {
        UserName: "",
        Password: "",
        UserID:   0,
        Email: "",
        PostsMade: 0,
        RepliesMade: 0,
        DateCreated: "",
        DateUpdated: "",
    };
    newUser.UserName = String(userName.value);
    newUser.Password = String(password.value);
    newUser.Email = String(email.value);

    var jsonString = JSON.stringify(newUser);
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/createUser', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            var SuccessMSG = JSON.parse(item);
            if (SuccessMSG.SuccessNum === 0){
                location.reload();
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