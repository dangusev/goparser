
var app = angular.module('Goparser', ["ui.router", "ui.bootstrap"]);

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
                "queries_add_view": {
                    templateUrl : "/templates/?tname=ajax/modal_form.html",
                    controller: "QueriesAddController"
                }
            }
        })
        .state("queries.update", {
            url: "queries/{queryId}/update/",
            views: {
                "queries_add_view": {
                    templateUrl : "/templates/?tname=ajax/modal_form.html",
                    controller: "QueriesUpdateController"
                }
            }
        })
        .state("queries.delete", {
            url: "queries/{queryId}/delete/",
            views: {
                "queries_add_view": {
                    controller: "QueriesDeleteController"
                }
            }
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
    .controller('QueriesAddController', ['$scope', '$http', '$state', function($scope, $http, $state) {
        $scope.Query = {};
        $scope.SubmitQuery = function (item, event) {
            $http.post(url="/api/queries/", data=$scope.Query).success(function (response) {
                    $state.go('queries',null, {reload: true});
                });
        };
}])
    .controller('QueriesUpdateController', ['$scope', '$http', '$state', '$stateParams', function($scope, $http, $state, $stateParams) {
        var url = "/api/queries/" + $stateParams.queryId + '/';
        $http.get(url).success(function (response){
            $scope.Query = response.query;
        });

        $scope.SubmitQuery = function (item, event) {
            $http.post(url, $scope.Query).success(function (response) {
                $state.reload();
            });
        };
}])
    .controller('QueriesDeleteController', ['$scope', '$http', '$state', '$stateParams', function($scope, $http, $state, $stateParams) {
        var url = "/api/queries/" + $stateParams.queryId + '/';
        $http.delete(url).success(function (response) {
            $state.go('queries', null , {reload:true});
        });
}])
    .controller('ItemsController', ['$scope', '$http', '$stateParams', function($scope, $http, $stateParams){

        $http.get('/api/queries/'+$stateParams.queryId+'/items/').success(function (response) {
            $scope.Items = response.items;
            $scope.Query = response.query;
        });
}]);