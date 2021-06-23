var app = angular.module('mymainpageApp', []);

var allLearnrs = [];

//Set a custom delimiter for templates
app.config(function($interpolateProvider) {
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
});

//Main Controller
app.controller('myCtrl', function($scope) {
    //Learnr/LearnRInforms Declarations
    $scope.jsLearnRArray = [];
    $scope.jsLearnInformArray = [];
    $scope.learnrAssemble = {
        ID: 0,
        InfoID: 0,
        OrgID: 0,
        Name: "",
        Tags: [],
        Description: [],
        PhoneNums: [],
        LearnrInforms: [],
        Active: true,
        DateCreated: "",
        DateUpdated: ""
    };
    $scope.learnrInforms = {
        ID: 0,
        Name: "",
        LearnRID: 0,
        LearnRName: "",
        Order: 0,
        TheInfo: "",
        ShouldWait: true,
        WaitTime: 0,
        DateCreated: "",
        DateUpdated: ""
    };
    //LearnrSet
    $scope.learnRSet = function() {
        var testJSONVar = String(getLearnrsAjax());
        $scope.SuccessMSG = JSON.parse(testJSONVar);
        $scope.jsLearnRArray = $scope.SuccessMSG.TheDisplayLearnrs;
        //$scope.jsLearnRArray = item;
        console.log("Here is our jsLearnRArray: " + JSON.stringify($scope.jsLearnRArray) + "\n");
        //return 1;
    };
    $scope.learnRSet();
    //Set functions for learnrInform
    $scope.informAdder = function(value) {
        console.log("Trying to add informAdder");
        $scope.learnrAssemble.LearnrInforms.push(value);
        return $scope.learnrAssemble.LearnrInforms;
    };
    //Set the learnr to array
    $scope.learnrAdder = function(value) {
        this.jsLearnRArray.push(value);
    };
    //Give all Learnrs in an array 
    $scope.giveLearnrs = function() {
        return this.jsLearnRArray;
    };
    $scope.names = ["Emil", "Tobias", "Linus"];
    $scope.testlearnrs = [
        {ID: 0, Name: "Cool1", Description: ["cool", "based"]},
        {ID: 5555, Name: "Cool2", Description: ["not cool", "not based"]}
    ];
    $scope.callLoadLearners = function() {
        addLearnerService.LearnrAdd();
    }
    //Test function for Javascript to call
    $scope.javatoAngularHello = function(value) {
        console.log("Hello from Angular. Here's a name: " + $scope.names[0] + ". And here's the value we were given: " +
        value);
        return $scope.names[0];
    };
    $scope.htmltoAngularHello = function(value) {
        console.log("Hello from Angular, HTML. Here's a name: " + $scope.names[0] + ". And here's the value we were given: " +
        value);
    }
});

//Add learnr add service
app.factory('addLearnerService', function(){
    return {
        LearnrAdd: function(){
            return 22;
        }
    };
});

//Javascript stuff to call Angular and vice versa
window.addEventListener('DOMContentLoaded', function(){
    //getLearnrsAjax(); //Get all Learnsrs loaded to Angular
});

function getLearnrsAjax(){
    var xhr = new XMLHttpRequest();
    xhr.open('GET', '/giveAllLearnrDisplay', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            var SuccessMSG = JSON.parse(item);
            if (SuccessMSG.SuccessNum === 0){
                return item;
            } else {
                console.log("Failed to reach out to giveAllLearnrDisplay");
                return item;
            }
        }
    });
    xhr.send("testsend");
}

/* Send the created LearnRInform into the learnrAssemble Array 
with LearnrInforms inside it */
function sendLearnRInform(){
    //Add this learnInform to the inform array in Angular
    var scope = angular.element(document.getElementById('theMainBody')).scope();
    var returnvalue = scope.informAdder(learnrInforms);
    console.log("Here is the return value: " + JSON.stringify(returnvalue));
    //Reset LearnrInform
    var newLearnrInform = {
        ID: 0,
        Name: "",
        LearnRID: 0,
        LearnRName: "",
        Order: 0,
        TheInfo: "",
        ShouldWait: true,
        WaitTime: 0,
        DateCreated: "",
        DateUpdated: ""
    };
    learnrInforms = newLearnrInform;
}

/* Add the learnr to our array once it's assembled */
function sendLearnR(){
    var scope = angular.element(document.getElementById('theMainBody')).scope();
    scope.learnrAdder(learnrAssemble);
    //Reset Learnr
    var newLearnr = {
        ID: 0,
        InfoID: 0,
        OrgID: 0,
        Name: "",
        Tags: [],
        Description: [],
        PhoneNums: [],
        LearnrInforms: [],
        Active: true,
        DateCreated: "",
        DateUpdated: ""
    };
    learnrAssemble = newLearnr;
}

/* DEBUG: GET THE LEARNRS FROM JAVASCRIPT */
function getLearnrs(){
    var allLearnrs = [];

    var scope = angular.element(document.getElementById('theMainBody')).scope();
    allLearnrs = scope.giveLearnrs();
    console.log("Here is the return value: " + JSON.stringify(allLearnrs));
}

/* Add LearnR values */
function setid(thevalue){
    learnrAssemble.ID = Number(thevalue);
}

function setinfoid(thevalue){
    learnrAssemble.InfoID = Number(thevalue);
}

function setorgid(thevalue){
    learnrAssemble.OrgID = Number(thevalue);
}

function setlearnrname(thevalue){
    learnrAssemble.Name = String(thevalue);
}

function setlearnractive(thevalue){
    learnrAssemble.Active = Boolean(thevalue);
}

function setlearnrdatecreated(thevalue){
    learnrAssemble.DateCreated = String(thevalue);
}

function setlearnrupdated(thevalue){
    learnrAssemble.DateUpdated = String(thevalue);
}

function settags(thevalue){
    learnrAssemble.Tags = String(thevalue);
}

function setdescription(thevalue){
    learnrAssemble.Description = String(thevalue);
}

function setlearnrphonenums(thevalue){
    learnrAssemble.PhoneNums = String(thevalue);
}


/* Add LearnRInform values */

function sendlearnrinformid(thevalue){
    learnrInforms.ID = Number(thevalue);
}

function setlearnrinformname(thevalue){
    learnrInforms.Name = String(thevalue);
}

function setlearnrnameinform(thevalue){
    learnrInforms.LearnRName = String(thevalue);
}

function setlearnrinformorder(thevalue){
    learnrInforms.Order = Number(thevalue);
}

function setlearnrinfo(thevalue){
    learnrInforms.TheInfo = String(thevalue);
}

function setshouldwaitinform(thevalue){
    learnrInforms.ShouldWait = Boolean(thevalue);
}

function setwaittimeinform(thevalue){
    learnrInforms.WaitTime = Number(thevalue);
}

function setdatecreatedinform(thevalue){
    learnrInforms.DateCreated = String(thevalue);
}

function setdateupdatedinform(thevalue){
    learnrInforms.DateUpdated = String(thevalue);
}