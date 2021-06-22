var app = angular.module('mymainpageApp', []);

//Set a custom delimiter for templates
app.config(function($interpolateProvider) {
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
});

//Main Controller
app.controller('myCtrl', function($scope) {
    $scope.jsLearnRArray = [];
    $scope.names = ["Emil", "Tobias", "Linus"];
    $scope.testlearnrs = [
        {ID: 0, Name: "Cool1", Description: ["cool", "based"]},
        {ID: 5555, Name: "Cool2", Description: ["not cool", "not based"]}
    ];
    //Add our learnrs to this array above
    $scope.loadLearnrs = function(){
        //Add map values to this array
        for (const [key, value] of learnrMap.entries()){
            this.jsLearnRArray.push(value);
        }
        for (var q = 0; q < this.jsLearnRArray.length; q++){
            console.log("Here is jsleanrr at spot " +
            q + ": " + this.jsLearnRArray[q].Name);
        }
        //console.log("Here is our JSON for the array: " + JSON.stringify(learnrArray));
    };
    $scope.loadLearnrs(); //Call the function directly above
});

