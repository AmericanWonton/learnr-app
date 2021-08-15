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

/* This funciton disables or button, until we recieve a response from the
server */
function disableSubmitting(){
    var buttonSubmit = document.getElementById("buttonSubmit");

    buttonSubmit.disabled = 'true';
}

/* This clears our form upon error or page reload */
window.addEventListener('DOMContentLoaded', function(){
    // Access the form element...
    const form = document.getElementById("excelForm");

    /* Define function for sending form data */
    function sendData(){
        const XHR = new XMLHttpRequest();
        const FD = new FormData(form);

        // Define what happens on successful data submission
        XHR.addEventListener("load",function(event){
            var item = XHR.responseText;
            var SuccessMSG = JSON.parse(item);
            if (SuccessMSG.SuccessNum === 0){
                //Good file submit, alert User
                alert(String(SuccessMSG.Message));
                form.reset();
            } else {
                //Bad file submit, alert User
                alert(String(SuccessMSG.Message));
                form.reset();
            }
        });

        // Define what happens in case of error
        XHR.addEventListener("error",function(event) {
            alert('Oops! Something went wrong.');
        });

        // Set up our request
        XHR.open("POST","/bulksend");

        // The data sent is what the user provided in the form
        XHR.send(FD);
    }
    

    //Add listening event to the form when submitted
    form.addEventListener("submit", function(event){
        event.preventDefault();

        sendData(); //Send Form data from JS
    });
});