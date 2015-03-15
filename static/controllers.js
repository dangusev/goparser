/**
 * Created by dan on 15.03.15.
 */

var app = angular.module('Goparser', []);
app.controller('QueriesController', ['$scope', '$http', function($scope,$http) {
    $http.get("/queries/").success(function (response) {
        $scope.Queries = response;
    });
}]);
