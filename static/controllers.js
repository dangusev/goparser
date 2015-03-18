
var app = angular.module('Goparser', ["ui.router"]);

app.config(function ($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise("/");

    $stateProvider
        .state("queries", {
            url: "/",
            templateUrl: "/templates/?tname=ajax/queries.html",
            controller: "QueriesController"
        }).
        state("queries.add", {
            url: "queries/add/",
            views: {
                "queries_add_view": {templateUrl : "/templates/?tname=ajax/queries_form.html"}
            },
            controller: "QueriesAddController"
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
    .controller('QueriesAddController', ['$scope', '$http', function($scope,$http) {
        $scope.Query = {};
        $scope.Query.SubmitQuery = function (item, event) {
            $http.post(url="/api/queries/", data=$scope.Query).success(function (response) {
                    console.log(response)
                });
        };
}])
    .controller('ItemsController', ['$scope', '$http', '$stateParams', function($scope, $http, $stateParams){

        $http.get('/api/queries/'+$stateParams.queryId+'/items/').success(function (response) {
            $scope.Items = response.items;
            $scope.Query = response.query;
        });
}]);

