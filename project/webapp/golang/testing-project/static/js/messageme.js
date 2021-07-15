/* A function to send me an email when the fields are filled out */
function messageMe(){
    //Declare variables
    var YourMessageInput = document.getElementById("mesMeTextArea");
    var messageResponseP = document.getElementById("messageResponseP");

    var MessageInfo = {
        YourNameInput: String(TheUser.Username),
        YourEmailInput: String(TheUser.Email),
        YourMessageInput: String(YourMessageInput.value),
        YourUserID: Number(TheUser.UserID),
        YourUser: TheUser
    };
    var jsonString = JSON.stringify(MessageInfo); //Stringify Data
    //Send Request to user message update page
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/emailMe', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            var ReturnMessage = JSON.parse(item);
            if (ReturnMessage.SuccOrFail == 0){
                messageResponseP.innerHTML = "Email sent. I'll write back soon!";
            } else {
                messageResponseP.innerHTML = "Sorry! I couldn't get your messsage.... :(";
            }
        }
    });
    xhr.send(jsonString);
}