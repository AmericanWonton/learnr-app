var app = angular.module('mymainpageApp', []);

var displayedTexts = [];

/* Changed at Golang templating to give ourselves special search display criteria */
var specialLearnRSearch = 0;

/* This sets our special learnr search variable above to determine if we
need to display special searched learnrs*/
/* SET LEARNR FUNC BEGINNING */
function setLearnnRSearchNew(variablepassed){
    specialLearnRSearch = Number(variablepassed);
}

/* SET LEARNR FUNC ENDING */

//Set a custom delimiter for templates
app.config(function($interpolateProvider) {
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
});

//Main Controller
app.controller('myCtrl', function($scope, $http) {
    $scope.hasCompleted = false; // Do not show data until http gets back with data
    $scope.LearnRArray = []; //The LearnRArray used for display
    $scope.LearnRCounter = 0; //Some variable for counting
    $scope.showDivMap={};
    $scope.showTextDivMap={};
    $scope.buttonDisableMap={};
    /* Call HTTP to get our data. If specialLearnRSearch == false, then
    we search for all LearnRs. If specialLearnRSearch == true, call HTTP to get the LearnRs
    we already searched for */
    if (specialLearnRSearch == 1){
        $http({
            method: 'GET',
            url: '/getSpecialLearnRAngular'
            }).then(function successCallback(response) {
            // this callback will be called asynchronously
            // when the response is available
            $scope.theCheck = response.data.LearnRArray;
            /* Initial Check to see if value is greater than 0 */
            if (!$scope.theCheck){
                //Data has not loaded or no search results were found!
                var errDataLoadP = document.getElementById("errDataLoadP");
                errDataLoadP.innerHTML = "No data returned for search!";
                $scope.hasCompleted = true; //Data load complete, we can show data in template
            } else {
                for (var i = 0; i < response.data.LearnRArray.length; i++){
                    $scope.LearnRArray.push(response.data.LearnRArray[i]);
                }
                $scope.hasCompleted = true; //Data load complete, we can show data in template
            }
            
            }, function errorCallback(response) {
            // called asynchronously if an error occurs
            // or server returns response with an error status.
            console.log("Error with returned Data! " + String(response));
            //Data has not loaded or no search results were found!
            var errDataLoadP = document.getElementById("errDataLoadP");
            errDataLoadP.innerHTML = "Error finding LearnR data!\n" + String(response);
            $scope.hasCompleted = true; //Data load complete, we can show data in template
        });
    } else {
        $http({
            method: 'GET',
            url: '/getLearnRAngular'
            }).then(function successCallback(response) {
            // this callback will be called asynchronously
            // when the response is available
            $scope.theCheck = response.data.LearnRArray;
            /* Initial Check to see if value is greater than 0 */
            if (!$scope.theCheck){
                //Data has not loaded or no search results were found!
                var errDataLoadP = document.getElementById("errDataLoadP");
                errDataLoadP.innerHTML = "No data returned for search!";
                $scope.hasCompleted = true; //Data load complete, we can show data in template
            } else {
                for (var i = 0; i < response.data.LearnRArray.length; i++){
                    $scope.LearnRArray.push(response.data.LearnRArray[i]);
                }
                $scope.hasCompleted = true; //Data load complete, we can show data in template
            }
            }, function errorCallback(response) {
            // called asynchronously if an error occurs
            // or server returns response with an error status.
            console.log("Error with returned Data! " + String(response));
            //Data has not loaded or no search results were found!
            var errDataLoadP = document.getElementById("errDataLoadP");
            errDataLoadP.innerHTML = "Error finding LearnR data!\n" + String(response);
            $scope.hasCompleted = true; //Data load complete, we can show data in template
        });
    }

    //Increment the Counter
    $scope.incrementCounter = function(aLearnR){
        $scope.LearnRCounter++;
        //console.log("The Return counter is: " + $scope.LearnRCounter + ". The LearnR Name is: " + aLearnR.Name);
        //Add LearnRID to map to have the divs shown as false
        $scope.showDivMap[aLearnR.ID] = false;
        $scope.showTextDivMap[aLearnR.ID] = false;
        $scope.buttonDisableMap[aLearnR.ID] = false;
    }

    //Compile LearnR description to one big string
    $scope.giveLearnRDescription = function(arrayODesc) {
        let bigString = "";
        for (var n = 0; n < arrayODesc.length; n++){
            bigString = bigString + arrayODesc[n] + " ";
        }
        return bigString;
    }

    //Return the counter
    $scope.returnCounter = function(){
        return $scope.LearnRCounter;
    }
    
    //Return if this page is LearnR sending info is visible
    $scope.returnVisibleLearnR = function(learnRID){
        return $scope.showDivMap[learnRID];
    }

    //Return if this page is LearnR text info is visible
    $scope.returnVisibleLearnRTexts = function(learnRID){
        return $scope.showTextDivMap[learnRID];
    }
    //Show a div based on a click
    $scope.divTextClicker = function(learnRID){
        if ($scope.showDivMap[learnRID] == true){
            $scope.showDivMap[learnRID] = false;
        } else {
            $scope.showDivMap[learnRID] = true;
        }
    }

    //Show LearnR Texts for a div based on a click
    $scope.divLearnRTextClicker = function(learnRID){
        if ($scope.showTextDivMap[learnRID] == true){
            $scope.showTextDivMap[learnRID] = false;
        } else {
            $scope.showTextDivMap[learnRID] = true;
        }
    }
    //Show a button based on a click
    $scope.buttonDisable = function(learnRID){
        $scope.buttonDisableMap[learnRID] = true;
    }
    //Return a unique ID based on counter
    $scope.uniqueIDInput = function(){
        return "fieldinputPersonName" + String($scope.LearnRCounter);
    }
    $scope.uniqueIDForm = function(){
        console.log("We should be sending this: " + "learnRSendForm" + String($scope.LearnRCounter));
        return "learnRSendForm" + String($scope.LearnRCounter);
    }

    //Used for sending finished form
    $scope.LearnRSending = function(learnR, learnRID){
        var fieldinputPersonName = document.getElementById("fieldinputPersonName" + String(learnRID));
        var fieldinputPersonPN = document.getElementById("fieldinputPersonPN" + String(learnRID));
        var fieldinputIntroduction = document.getElementById("fieldinputIntroduction" + String(learnRID));
        var sendLearnRButton = document.getElementById("sendLearnRButton" + String(learnRID));
        var sendLearnRResult = document.getElementById("sendLearnRResult" + String(learnRID));

        var OurJSON = {
            TheUser: TheUser,
            TheLearnR: learnR,
            TheLearnRInfo: {},
            PersonName: String(fieldinputPersonName.value),
            PersonPhoneNum: String(fieldinputPersonPN.value.toString()),
            Introduction: String(fieldinputIntroduction.value)
        };

        $http({
            method: 'POST',
            url: '/canSendLearnR',
            data: JSON.stringify(OurJSON)
            }).then(function successCallback(response) {
            // this callback will be called asynchronously
            // when the response is available
            sendLearnRButton.disabled = true;
            console.log(response.data);
            console.log(response.data.ResultMsg);
            sendLearnRResult.innerHTML = "Your LearnR has successfully been started!";
            setTimeout(() => {navigateHeader(3);}, 5000);
            }, function errorCallback(response) {
            // called asynchronously if an error occurs
            // or server returns response with an error status.
            alert("Error returned: " + String(response.data.Message));
            console.log("Error with returned Data! " + String(response.data.Message));
            //Report Error
            sendLearnRResult.innerHTML = "Error: "+ String(response.data.Message);
            sendLearnRResult.value = "Error: "+ String(response.data.Message);
            //Clear input data
            fieldinputPersonName.value = "";
            fieldinputPersonName.innerHTML = "";
            fieldinputPersonPN.value = "";
            fieldinputPersonPN.innerHTML = "";
            fieldinputIntroduction.value = "";
            fieldinputIntroduction.innerHTML = "";
        });
    }

    //Used for searching LearnRs based off tag or name
    $scope.searchLearnRs = function(){
        var learnRNameInput = document.getElementById("learnRNameInput");
        var learnRTagInput = document.getElementById("learnRTagInput");
        var resultThing = document.getElementById("resultThing");

        var SearchJSON = {
            TheNameInput: String(learnRNameInput.value),
            TheTagInput: String(learnRTagInput.value)
        };

        $http({
            method: 'POST',
            url: '/searchLearnRs',
            data: JSON.stringify(SearchJSON)
            }).then(function successCallback(response) {
            // this callback will be called asynchronously
            // when the response is available
            
            //Determine response returned
            if (response.data.SuccessNum == 0){
                //Reload page, server should have new learnrs
                navigateHeader(3, 0);
            } else {
                //Error show error
                resultThing.innerHTML = "Error finding those LearnRs! " + String(response.data.Message);
                learnRNameInput.value = "";
                learnRTagInput.value = "";
            }
            }, function errorCallback(response) {
            // called asynchronously if an error occurs
            // or server returns response with an error status.
            resultThing.innerHTML = "Error finding those LearnRs! " + String(response.data.Message);
            learnRNameInput.value = "";
            learnRTagInput.value = "";
        });
    }

    //Handle the printed LearnRstuff
    $scope.LearnRPageAdd = function(){

    }
});

//Javascript stuff to call Angular and vice versa
window.addEventListener('DOMContentLoaded', function(){
    
});

//Used to control the search for LearnRs
function learnRSearch(){
    var learnRNameInput = document.getElementById("learnRNameInput");
    var learnRTagInput = document.getElementById("learnRTagInput");
    var resultThing = document.getElementById("resultThing");

    
    var SearchJSON = {
        TheNameInput: String(learnRNameInput.value),
        TheTagInput: String(learnRTagInput.value)
    };
    
    //Send Ajax
    var jsonString = JSON.stringify(SearchJSON); //Stringify Data
    //Send Request to change page
    
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/searchLearnRs', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            var ReturnData = JSON.parse(item);
            if (ReturnData.SuccessNum == 0){
                learnRNameInput.value = "";
                learnRTagInput.value = "";
                //Take action if nothing is returned
                if (ReturnData.ReturnLearnRs != null && ReturnData.ReturnLearnRs.length >= 1){
                    //Repopulate learnrs
                    
                } else {
                    //Nothing returned
                    resultThing.innerHTML = "No LearnRs returned from search!";
                }
            } else {
                resultThing.innerHTML = "Error finding those LearnRs! " + ReturnData.Message;
                learnRNameInput.value = "";
                learnRTagInput.value = "";
            }
        }
    });
    xhr.send(jsonString);
}

