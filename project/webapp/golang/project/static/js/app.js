var TheUser = {
    Username: "",
    Password: "",
    Firstname: "",
    Lastname: "",
    PhoneNums: new Array(),
    UserID: 0,
    Email: new Array(),
    Whoare: "",
    AdminOrgs: new Array(),
    OrgMember: new Array(),
    Banned: false,
    DateCreated: "",
    DateUpdated: ""
};

var isBanned = false;

function setUser(theUser){
    TheUser = theUser;
}

function setBanned(banned){
    isBanned = banned;
}

//Handles User clicking log out to delete their session
function logOut(){
    var jsonString = JSON.stringify(TheUser);
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/logUserOut', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            var SuccessMSG = JSON.parse(item);
            if (SuccessMSG.SuccessNum === 0){
                navigateHeader(6,0);
            } else {
                console.log("DEBUG: We have an error: " + SuccessMSG.SuccessNum + " " +
                SuccessMSG.Message);
                navigateHeader(6,0);
            }
        }
    });
    xhr.send(jsonString);
}