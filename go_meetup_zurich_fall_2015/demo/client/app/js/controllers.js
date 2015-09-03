'use strict';

/* Controllers */

function getComments(toGet, scope, Comment) {
    if (toGet.length > 0) {
        Comment.get({commentId: toGet[0]}, function(comment) {
            scope.comments.push(comment);
            toGet.pop();
            getComments(toGet, scope, comment)
        });
    }
}

var phonecatControllers = angular.module('phonecatControllers', []);

phonecatControllers.controller('PhoneListCtrl', ['$scope', 'Pic',
  function($scope, Pic) {
    $scope.phones = Pic.query();
    $scope.orderProp = 'age';
  }]);

phonecatControllers.controller('PhoneDetailCtrl',
    ['$scope', '$routeParams', '$http', '$window', 'Pic', 'Comment', 'Like',
  function($scope, $routeParams, $http, $window, Pic, Comment, Like) {
    $scope.currentUserId = undefined;

    if ($window.sessionStorage.username) {
        $scope.currentUserId = $window.sessionStorage.username;
    }

    $scope.pic = Pic.get({picId: $routeParams.phoneId}, function(pic) {
      $scope.mainImageUrl = pic.image_url;
      $scope.comments = [];

      $scope.getComments(pic.comments);
    });
    $scope.getComments = function(toGet) {
        if (toGet.length > 0) {
            Comment.get({commentId: toGet[0]}, function(comment) {
                $scope.comments.push(comment);
                toGet.shift();
                $scope.getComments(toGet)
            });
        }
    }

    $scope.newComment = new Comment();

    $scope.addComment = function() {
        $scope.newComment.userid = $scope.currentUserId;
        $scope.newComment.picid = parseInt($routeParams.phoneId);
        $scope.newComment.$save(function(comment) {
            $scope.comments.push(comment);
            $scope.newComment = new Comment();
            console.log("ok posted", comment);
        });
    }

    $scope.like = new Like();

    $scope.likePic = function() {
        $scope.like.$save({picId: $routeParams.phoneId}, function(like) {
            $scope.pic.liked.push($scope.currentUserId);
        });
    }

    $scope.checkLiked = function() {
        if ($scope.currentUserId && $scope.pic.liked) {
            for (var i = 0; i != $scope.pic.liked.length; ++i) {
                if ($scope.pic.liked[i] == $scope.currentUserId) {
                    return true;
                }
            }

            return false;
        }

        return false;
    }

    $scope.logout = function() {
        delete $window.sessionStorage.token;
        delete $window.sessionStorage.username;
        $scope.currentUserId = undefined;
    }

    $scope.submit = function () {
        $http
        .post('/api/login', $scope.user)
        .success(function (data, status, headers, config) {
            $window.sessionStorage.token = data.token;
            $window.sessionStorage.username = $scope.user.username;
            $scope.currentUserId = $scope.user.username;
        })
        .error(function (data, status, headers, config) {
            $scope.logout();
        })
    }
}]);
