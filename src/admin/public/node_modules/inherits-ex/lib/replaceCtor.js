var getPrototypeOf = require('./getPrototypeOf');
var setPrototypeOf = require('./setPrototypeOf');
var defineProperty = Object.defineProperty;

//just replace the object's constructor
module.exports = function(aObject, aClass) {
  var vOldProto = getPrototypeOf(aObject);
  var result = false;
  if ( vOldProto && vOldProto !== aClass.prototype) {
    if (!aClass.prototype.hasOwnProperty('Class')) {
      defineProperty(aClass.prototype, 'Class', {
        value: aClass,
        configurable: true
      });
    }
    setPrototypeOf(aObject, aClass.prototype);
    result = true;
  }
  return result;
}

