function routeAPI($routeProvider) {
    $routeProvider.when("/apis", {
        "templateUrl": "html/api/list.html",
        "controller": APIController
    }).when("/apis/:url/:method", {
        "templateUrl": "html/api/update.html",
        "controller": APIUpdateController
    }).when("/new/api", {
        "templateUrl": "html/api/new.html",
        "controller": APICreateController
    });
}

function APIUpdateController($scope, $routeParams, $http, $location, $route) {
    $http.get("api/clusters").success(function (data) {
        $scope.clusters = data.value;
    });

    $http.get("api/apis/" + $routeParams.url + "?method=" + $routeParams.method).success(function (data) {
        $scope.api = data.value;
        $scope.oldMethod = data.value.method;
        $scope.oldUrl = data.value.url;
    });


    $scope.resetNode = function () {
        $scope.newClusterName = "";
        $scope.newNodeUrl = "";
        $scope.newNodeAttrName = "";
        $scope.newNodeRewrite = "";
        $scope.newValidation = [];
    }

    $scope.addNode = function () {
        $scope.api.nodes.push({
            "clusterName": $scope.newClusterName,
            "url": $scope.newNodeUrl,
            "attrName": $scope.newNodeAttrName,
            "rewrite": $scope.newNodeRewrite,
            "validations": $scope.newValidation,
        });

        $scope.resetNode();
    }

    $scope.delete = function (node) {
        ns = []

        for (var i = 0; i < $scope.api.nodes.length; i++) {
            if ($scope.api.nodes[i] != node) {
                ns.push($scope.api.nodes[i]);
            }
        }

        $scope.api.nodes = ns;
    }

    $scope.update = function () {
        $scope.api.status = parseInt($scope.api.status)
        $http.put('api/apis', $scope.api).success(function (data) {
            if ($scope.oldMethod != $scope.api.method || $scope.oldUrl != $scope.api.url) {
                $http.delete('api/apis/' + Base64.encodeURI($scope.oldUrl) + "?method=" + $scope.oldMethod).success(function (data) {
                    $location.path("/apis");
                    $route.reload();
                });
            } else {
                $location.path("/apis");
                $route.reload();
            }
        });
    }
}


function APICreateController($scope, $routeParams, $http, $location, $route) {
    $http.get("api/clusters").success(function (data) {
        $scope.clusters = data.value;
    });

    $scope.newUrl = "";
    $scope.newDesc = "";
    $scope.newName = "";
    $scope.newMock = {};
    $scope.newAccessControl = {};
    $scope.newNodes = [];
    $scope.newStatus = 0;


    $scope.resetNode = function () {
        $scope.newClusterName = "";
        $scope.newNodeUrl = "";
        $scope.newNodeAttrName = "";
        $scope.newNodeRewrite = "";
        $scope.newValidation = [];
    }

    $scope.addNode = function () {
        $scope.newNodes.push({
            "clusterName": $scope.newClusterName,
            "url": $scope.newNodeUrl,
            "attrName": $scope.newNodeAttrName,
            "rewrite": $scope.newNodeRewrite,
            "validations": $scope.newValidation,
        });

        $scope.resetNode();
    }

    $scope.delete = function (node) {
        ns = []

        for (var i = 0; i < $scope.newNodes.length; i++) {
            if ($scope.newNodes[i] != node) {
                ns.push($scope.newNodes[i]);
            }
        }

        $scope.newNodes = ns;
    }

    $scope.add = function () {
        d = {
            "name": $scope.newName,
            "url": $scope.newUrl,
            "method": $scope.newMethod,
            "status": parseInt($scope.newStatus),
            "desc": $scope.newDesc,
            "mock": $scope.newMock,
            "accessControl": $scope.newAccessControl,
            "nodes": $scope.newNodes,
        }

        $http.post('api/apis', d).success(function (data) {
            $location.path("/apis");
            $route.reload();
        });
    }
}

function APIController($scope, $routeParams, $http, $location, $route) {
    $http.get("api/apis").success(function (data) {
        $scope.apis = data.value;

        for (var i = 0; i < $scope.apis.length; i++) {
            $scope.apis[i].encodeURL = Base64.encodeURI($scope.apis[i].url);
        }
    });

    $scope.create = function () {
        $location.path("/new/api");
    }

    $scope.delete = function (url, method) {
        $http.delete('api/apis/' + url + "?method=" + method).success(function (data) {
            $location.path("/apis");
            $route.reload();
        });
    }

    $scope.up = function (api) {
        if (api.status == 1) {
            return;
        }

        var old = api.status;
        api.status = 1;
        $http.put('api/apis', api).error(function (err) {
            api.status = old;
        });
    }

    $scope.down = function (api) {
        if (api.status == 0) {
            return;
        }

        api.status = 0;
        $http.put('api/apis', api).error(function (err) {
            api.status = old;
        });
    }
}