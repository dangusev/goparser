/**
 * Created by dan on 15.03.15.
 */

var app = angular.module('Goparser', ["ui.router"]);

app.config(function ($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise("/");

    $stateProvider
        .state("main", {
            url: "/",
            templateUrl: "/templates/main.html",
            controller: "QueriesController"
        });
    //$routeProvider.when("/", {
    //    templateUrl: "/templates/main.html",
    //    controller: "QueriesController"
    //})
    //.when("/queries/:param/items/", {
    //    templateUrl: "/templates/items.html",
    //    controller: "ItemsController"
    //});
    //.otherwise({
    //    redirectTo: "/"
    //})
});


app.controller('QueriesController', ['$scope', '$http', function($scope,$http) {
    $http.get("/api/queries/").success(function (response) {
        $scope.Queries = response;
    });
}]);
