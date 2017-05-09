var inherits = require('./inherits');
var getPrototypeOf = require('./getPrototypeOf');
var setPrototypeOf = require('./setPrototypeOf');

//make sure the aClass.prototype hook to the aObject instance.

module.exports = function(aObject, aClass) {
  // ES6: Object.getPrototypeOf / Object.setPrototypeOf
  var vOldProto = getPrototypeOf(aObject);
  var result = false;
  if ( vOldProto !== aClass.prototype) {
    inherits(aClass, vOldProto.constructor);
    setPrototypeOf(aObject, aClass.prototype);
    result = true;
  }
  return result;
};
