var app = angular.module('mymainpageApp', []);

var allLearnrs = [];

//Set a custom delimiter for templates
app.config(function($interpolateProvider) {
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
});

//Main Controller
app.controller('myCtrl', function($scope, $timeout) {
    /* Use for casedata loading */
    $scope.caseData = null;
    //Learnr/LearnRInforms Declarations
    $scope.jsLearnRArray = [];
    $scope.jsLearnInformArray = [];
    /* LearnrSet
    Calls Ajax to get our Learnrs and put them into our jsLearnRArray */
    $scope.learnRSet = function() {
        var xhr = new XMLHttpRequest();
        xhr.open('GET', '/giveAllLearnrDisplay', true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.addEventListener('readystatechange', function(){
            if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
                var item = xhr.responseText;
                var SuccessMSG = JSON.parse(item);
                if (SuccessMSG.SuccessNum === 0){
                    $scope.jsLearnRArray = SuccessMSG.TheDisplayLearnrs;
                    console.log("DEBUG: Here is our jsLearnRArray: " + JSON.stringify($scope.jsLearnRArray));
                    $scope.caseData = 'hey!';
                } else {
                    console.log("Failed to reach out to giveAllLearnrDisplay");
                }
            }
        });
        xhr.send("testsend");
    };
    $scope.learnRSet();
    //mimic a delay in getting the data from $http
    $timeout(function () {
        $scope.caseData = 'hey!';
    }, 1000);
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
    
});