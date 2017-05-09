(function() {
  var getClass, isArray, isObject, isString;

  isArray = Array.isArray;

  isString = function(v) {
    return typeof v === 'string';
  };

  isObject = function(v) {
    return v && typeof v === 'object';
  };

  getClass = function(aClassName, aScope, aValues) {
    var vKeys;
    if (isArray(aScope)) {
      if (!isArray(aValues)) {
        aValues = aScope;
        aScope = aValues.map(function(k) {
          var result;
          result = k.name;
          if (result == null) {
            throw TypeError('No Scope Name for ' + k);
          }
          return result;
        });
      }
    } else if (isObject(aScope)) {
      vKeys = Object.keys(aScope);
      aValues = vKeys.map(function(k) {
        return aScope[k];
      });
      aScope = vKeys;
    } else {

    }
    return Function(aScope, 'return ' + aClassName).apply(null, aValues);
  };

  module.exports = function(aClassName, aScope, aValues) {
    var result;
    if (aClassName != null) {
      if (typeof aClassName === 'function') {
        result = aClassName;
      } else if (typeof aClassName === 'string') {
        if (!/[, {}();.]+/.test(aClassName)) {
          result = getClass(aClassName, aScope, aValues);
        }
      }
    }
    return result;
  };

}).call(this);

//# sourceMappingURL=get-class-by-name.js.map
