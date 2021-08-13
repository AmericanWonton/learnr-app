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
    /*
    var form = document.getElementById("excelForm");
    form.onsubmit = function(){
        form.reset();
        console.log("Form cleared");
    }
    */
});

function formSubmit(){
    /* Get file name for form submission */
    var filename;
    var fullPath = document.getElementById('excelFileInput').value;
    if (fullPath) {
        var startIndex = (fullPath.indexOf('\\') >= 0 ? fullPath.lastIndexOf('\\') : fullPath.lastIndexOf('/'));
        filename = fullPath.substring(startIndex);
        if (filename.indexOf('\\') === 0 || filename.indexOf('/') === 0) {
            filename = filename.substring(1);
        }
    }
    //Get Form Data
    var data = new FormData();
    data.append('excel-file', document.getElementById("excelFileInput").value, String(filename));
    data.append('learnR', document.getElementById("learnR").value);
    data.append('hiddenFormValue', document.getElementById("hiddenFormValue").value);


    //Ajax
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/bulksend', true);
    xhr.onload = function(){
        console.log("Form submitted: " + this.response);
        //MANUAL RESET
        document.getElementById("excelFileInput").value = "";
        document.getElementById("learnR").value = "";
    };
    xhr.send(data);

    // (C) STOP DEFAULT FORM SUBMIT/PAGE RELOAD
    return false;
}