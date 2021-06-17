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

/* SET USER FUNC BEGINNING */
function setUsername(username){
    TheUser.Username = String(username);
}

function setPassword(password){
    TheUser.Password = String(password);
}

function setFirstname(firstname){
    TheUser.Firstname = String(firstname);
}

function setLastname(lastname){
    TheUser.Lastname = String(lastname);
}

function setPhoneNums(phonenums){
    TheUser.PhoneNums.push(phonenums);
}

function setUserID(userid){
    TheUser.UserID = userid;
}

function setEmail(emails){
    TheUser.Email.push(emails);
}

function setWhoAre(whoare){
    TheUser.Whoare = whoare;
}

function setAdminOrgs(adorgs){
    TheUser.AdminOrgs.push(adorgs);
}

function setOrgMember(themembersorg){
    TheUser.OrgMember.push(themembersorg);
}

function setBanned(banned){
    isBanned = banned;
    TheUser.Banned = banned;
}

function setDateCreated(datecreated){
    TheUser.DateCreated = datecreated;
}

function setDateUpdated(dateupdated){
    TheUser.DateUpdated = dateupdated;
}

/* SET USER FUNC ENDING */

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