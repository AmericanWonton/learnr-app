/* This function displays our hidden instructions, or hides it */
function displayInstructions(){
    var infoHidden = document.getElementById("infoHidden");
    var seeThingsB = document.getElementById("seeThingsB");

    if (infoHidden.style.display === "none"){
        //Reveal the information
        infoHidden.style.display = "flex";
        seeThingsB.innerHTML = "Hide Instructions";
    } else {
        //Hide the information
        infoHidden.style.display = "none";
        seeThingsB.innerHTML = "Show Instructions";
    }
}