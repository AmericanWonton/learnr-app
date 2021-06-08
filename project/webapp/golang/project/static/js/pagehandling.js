
//Used to control which link to send our user to; also creates game session
function navigateHeader(whichLink, whichgame) {
    console.log("whichgame: " + whichgame);
    switch (whichLink) {
        case 1:
            //Go to Login Page
            window.location.assign("/login");
            break;
        case 2:
            //Go to Google
            window.location.assign("https://www.google.com");
            break;
        case 3:
            //Go to mainpage
            window.location.assign("/mainpage");
            break;
        case 4:
            //Go to aboutme page
            window.location.assign("/aboutme");
            break;
        case 5:
            //Go to messageme page
            window.location.assign("/messageme");
            break;
        case 6:
            //Go to Index
            window.location.assign("/");
            break;
        case 7:
            //Go to Rules
            window.location.assign("/rules");
            break;
        case 8:
            //Game Page
            window.location.assign("/gamepage");
            break;
        case 9:
            //Game Page
            window.location.assign("/admin");
            break;
        default:
            console.log("Error, wrong whichLink entered: " + whichLink);
            break;
    }
}