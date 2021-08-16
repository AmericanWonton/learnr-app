
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
            //Go to learnmore page
            window.location.assign("/learnmore");
            break;
        case 5:
            //Go to Send Help Page
            window.location.assign("/sendhelp");
            break;
        case 6:
            //Go to Index
            //Need to log User out

            window.location.assign("/");
            break;
        case 7:
            //Go to Rules
            window.location.assign("/signup");
            break;
        case 8:
            //Make Organization
            window.location.assign("/makeorg");
            break;
        case 9:
            //Admin Page
            window.location.assign("/admin");
            break;
        case 10:
            //Make LearnR Page
            window.location.assign("/learnr");
        case 11:
            //Bulk LearnR Send
            window.location.assign("/bulksend");
            break;
        default:
            console.log("Error, wrong whichLink entered: " + whichLink);
            break;
    }
}