//Add our listening events to the window loading
window.addEventListener('DOMContentLoaded', function(){
    //Define elements of the 'Ajax Form'
    var inputTextMobileUN = document.getElementById("inputTextMobileUN");
    var inputTextMobilePW = document.getElementById("inputTextMobilePW");
    var submitLoginButton = document.getElementById("submitLoginButton");
    var passwordErr = document.getElementById("password-err");
    var divformDivLogin = document.getElementById("divformDivLogin");

    /* When clicked, submit the results of the login items;
    if successful, it will redirect to the mainpage with your newly made cookie.
    */
    submitLoginButton.addEventListener("click", function(){
        var LoginData = {
            Username: String(inputTextMobileUN.value),
            Password: String(inputTextMobilePW.value)
        };

        var jsonString = JSON.stringify(LoginData); //stringify JSON

        //Call Ajax to see if password/username are correct
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/canLogin', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.addEventListener('readystatechange', function(){
            if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var Response = xhr.responseText;
                var SuccessMSG = JSON.parse(Response);
                if (Number(SuccessMSG.SuccessNum) === 0){
                    //Successful User Search
                    //Increase size of Div
                    divformDivLogin.style.height = "400px";
                    //Clear passwordErr
                    passwordErr.innerHTML = "";
                    passwordErr.innerHTML = SuccessMSG.Message;
                    //This User should have their cookie. Send them to the choice page
                    navigateHeader(3);
                } else if (Number(SuccessMSG.SuccessNum) === 1){
                    //Failed User Search
                    //Increase size of Div
                    divformDivLogin.style.height = "400px";
                    //Clear passwordErr
                    passwordErr.innerHTML = "";
                    passwordErr.innerHTML = SuccessMSG.Message;
                } else {
                    //Failed User Search
                    //Increase size of Div
                    divformDivLogin.style.height = "400px";
                    //Clear informtextPSignIn
                    passwordErr.innerHTML = "";
                    passwordErr.innerHTML = "Wrong information returned from server";
                    console.log("Wrong number returned from canLogin: " + SuccessMSG.SuccessNum);
                }
            }
        });
        xhr.send(jsonString);
    });
});