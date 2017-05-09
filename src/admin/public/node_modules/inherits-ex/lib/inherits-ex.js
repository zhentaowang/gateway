
/*
The enhanced dynamical `inherits` implementation.

```js
  var InheritsEx = require('inherits-ex/lib/inherits-ex')
  var defaultRequire = InheritsEx.requireClass;
  // You should return the proper class(ctor) here.
  InheritsEx.requireClass = function(className, scope){return defaultRequire.apply(null, arguments)};
  var inherits = InheritsEx()
```

The enhanced dynamical `inherits` implementation.

+ load the class via dynamical name.
* requireClass *(Function)*:
* scope *(Object)*: collects the register classes.
  * the inherits ctor will be added into the scope automatically.

The default requireClass is `getClassByName`.

usage:

```coffee
  requireClass = (aClassName, aScope)->
    getClassViaName aClassName
  InheritsEx = require('inherits-ex/lib/inherits-ex')
  inherits   = InheritsEx(requireClass)

  class RootClass
  class Parent
  InheritsEx.setScope RootClass:RootClass, Parent:Parent

  class MyClass3
    inherits MyClass3, RootClass

```
 */

(function() {
  var InheritsEx, getClassByName, inherits, isArray, isFunction, isObject, isString;

  getClassByName = require('./get-class-by-name');

  inherits = require('./inherits');

  isFunction = function(value) {
    return typeof value === 'function';
  };

  isString = function(value) {
    return typeof value === 'string';
  };

  isObject = function(value) {
    return typeof value === 'object';
  };

  isArray = Array.isArray;

  module.exports = InheritsEx = (function() {
    InheritsEx.requireClass = getClassByName;

    InheritsEx.scope = {};

    InheritsEx.setScope = function(aScope) {
      var j, k, len, vName;
      if (isArray(aScope)) {
        for (j = 0, len = aScope.length; j < len; j++) {
          k = aScope[j];
          vName = k.name;
          if (vName == null) {
            throw TypeError('No Scope Name for ' + k);
          }
          InheritsEx.scope[vName] = k;
        }
      } else if (isObject(aScope)) {
        InheritsEx.scope = aScope;
      }
    };

    InheritsEx.getClass = function(aClassName, aScope, aValues) {
      var requireClass, result;
      requireClass = InheritsEx.requireClass;
      if (aScope != null) {
        result = requireClass(aClassName, aScope, aValues);
      }
      if (!result && (aScope = InheritsEx.scope)) {
        result = requireClass(aClassName, aScope);
      }
      return result;
    };

    InheritsEx.execute = function(ctor, superCtors, aScope, aValues) {
      var getClass, i, isStrCtor, result;
      getClass = InheritsEx.getClass;
      if (isStrCtor = isString(ctor)) {
        ctor = getClass(ctor, aScope, aValues);
      }
      if (isString(superCtors)) {
        superCtors = getClass(superCtors, aScope, aValues);
      } else if (isArray(superCtors)) {
        superCtors = (function() {
          var j, len, results;
          results = [];
          for (j = 0, len = superCtors.length; j < len; j++) {
            i = superCtors[j];
            if (isString(i)) {
              results.push(getClass(i, aScope, aValues));
            } else {
              results.push(i);
            }
          }
          return results;
        })();
      }
      result = inherits(ctor, superCtors);
      if (result && !isStrCtor) {
        InheritsEx.scope[ctor.name] = ctor;
      }
      return result;
    };

    function InheritsEx(aDefaultRequire) {
      if (isFunction(aDefaultRequire)) {
        InheritsEx.requireClass = aDefaultRequire;
      }
      return InheritsEx.execute;
    }

    return InheritsEx;

  })();

}).call(this);

//# sourceMappingURL=inherits-ex.js.map
