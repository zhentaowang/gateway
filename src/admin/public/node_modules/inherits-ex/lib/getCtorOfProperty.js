(function() {
  var getPrototypeOf;

  getPrototypeOf = require('./getPrototypeOf');

  module.exports = function(aClass, aPropertyName) {
    var result, vPrototype;
    vPrototype = aClass.prototype;
    while (vPrototype && !vPrototype.hasOwnProperty(aPropertyName)) {
      vPrototype = getPrototypeOf(vPrototype);
    }
    if (vPrototype) {
      if (vPrototype.hasOwnProperty('Class')) {
        result = vPrototype.Class;
      } else {
        result = vPrototype.constructor;
      }
    }
    return result;
  };

}).call(this);

//# sourceMappingURL=getCtorOfProperty.js.map
