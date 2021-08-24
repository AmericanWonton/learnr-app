var app = angular.module('mymainpageApp', []);

var displayedTexts = [];

/* This takes the learnr array we've created and begins to list it on our page.
Divs will be created, being added into 'learnrHolderDiv'*/


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
    /* Call HTTP to get our data */
    $http({
        method: 'GET',
        url: '/getLearnRAngular'
        }).then(function successCallback(response) {
        // this callback will be called asynchronously
        // when the response is available
        console.log(response.data);
        console.log(response.data.ResultMsg);
        for (var i = 0; i < response.data.LearnRArray.length; i++){
            console.log("This is LearnR: " + response.data.LearnRArray[i].Name);
            $scope.LearnRArray.push(response.data.LearnRArray[i]);
        }

        $scope.hasCompleted = true; //Data load complete, we can show data in template
        }, function errorCallback(response) {
        // called asynchronously if an error occurs
        // or server returns response with an error status.
        console.log("Error with returned Data! " + String(response));
    });

    //Increment the Counter
    $scope.incrementCounter = function(aLearnR){
        $scope.LearnRCounter++;
        console.log("The Return counter is: " + $scope.LearnRCounter + ". The LearnR Name is: " + aLearnR.Name);
        //Add LearnRID to map to have the divs shown as false
        $scope.showDivMap[aLearnR.ID] = false;
        //Debug Print
        /*
        angular.forEach($scope.showDivMap, function(value, key){
            console.log("Map Key: " + key +  " Map Value: " + value);
        });
        */
        //console.log("Map is currently: " + $scope.showDivMap);
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

    //Show a div based on a click
    $scope.divTextClicker = function(learnRID){
        if ($scope.showDivMap[learnRID] == true){
            $scope.showDivMap[learnRID] = false;
        } else {
            $scope.showDivMap[learnRID] = true;
        }
    }

    //Return a unique ID based on counter
    $scope.uniqueIDInput = function(){
        return "fieldinputPersonName" + String($scope.LearnRCounter);
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
    
    //console.log("DEBUG: Here is special cases: " + TheSpecialCases);
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
                if (ReturnData.ReturnLearnRs != null && ReturnData.ReturnLearnRs){
                    //Repopulate learnrs
                    rePopulateLearnRs(ReturnData.ReturnLearnRs);
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

