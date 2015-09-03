'use strict';

/* Services */

var phonecatServices = angular.module('phonecatServices', ['ngResource']);

phonecatServices.factory('Pic', ['$resource',
  function($resource){
    return $resource('api/pics/:picId', {}, {
      query: {method:'GET', isArray:true}
    });
  }]);

phonecatServices.factory('Comment', ['$resource',
  function($resource){
    return $resource('api/comments/:commentId', {}, {
      query: {method:'GET', isArray:true}
    });
  }]);

phonecatServices.factory('Like', ['$resource',
  function($resource){
    return $resource('api/pics/:picId/like', {}, {
      query: {method:'GET', isArray:true}
    });
  }]);
