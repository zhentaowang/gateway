var isEmptyFunction = require('./isEmptyFunction')

//get latest non-empty constructor function through inherits link:
module.exports = function (ctor) {
  var result = ctor;
  var isEmpty = isEmptyFunction(result);
  while (isEmpty && (result.super_)) {
    result  = result.super_;
    isEmpty = isEmptyFunction(result);
  }
  //if (isEmpty) result = null;
  return result;
}

