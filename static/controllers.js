
var app = angular.module('Goparser', ["ui.router"]);

app.config(function ($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise("/");

    $stateProvider
        .state("main", {
            url: "/",
            templateUrl: "/templates/?tname=ajax/queries.html",
            controller: "QueriesController"
        })
        .state("items", {
            url: "/queries/{queryId}/items/",
            templateUrl: "/templates/?tname=ajax/items.html",
            controller: "ItemsController"
        });
});


app.controller('QueriesController', ['$scope', '$http', function($scope,$http) {
    $http.get("/api/queries/").success(function (response) {
        $scope.Queries = response.queries;
    });
}])
    .controller('ItemsController', ['$scope', '$http', '$stateParams', function($scope, $http, $stateParams){

        $http.get('/api/queries/'+$stateParams.queryId+'/items/').success(function (response) {
            $scope.Items = response.items;
            $scope.Query = response.query
        });
}]);
