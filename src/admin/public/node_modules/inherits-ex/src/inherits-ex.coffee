###
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

###

getClassByName  = require './get-class-by-name'
inherits        = require './inherits'

isFunction      = (value)->typeof value is 'function'
isString        = (value)->typeof value is 'string'
isObject        = (value)-> typeof value is 'object'
isArray         = Array.isArray

module.exports  = class InheritsEx
  @requireClass: getClassByName
  @scope: {}
  @setScope: (aScope)->
    if isArray aScope
      for k in aScope
        vName = k.name
        throw TypeError 'No Scope Name for ' + k unless vName?
        InheritsEx.scope[vName] = k
    else if isObject aScope
      InheritsEx.scope = aScope
    return
  # get the class from scope.
  @getClass: (aClassName, aScope, aValues)->
    requireClass = InheritsEx.requireClass
    result = requireClass(aClassName, aScope, aValues) if aScope?
    result = requireClass(aClassName, aScope) if !result and aScope = InheritsEx.scope
    result
  @execute = (ctor, superCtors, aScope, aValues)->
    getClass  = InheritsEx.getClass
    ctor      = getClass ctor, aScope, aValues if isStrCtor = isString ctor
    if isString superCtors
      superCtors    = getClass superCtors, aScope, aValues
    else if isArray superCtors
      superCtors = for i in superCtors
        if isString(i) then getClass(i, aScope, aValues) else i
    result = inherits ctor, superCtors
    InheritsEx.scope[ctor.name] = ctor if result and !isStrCtor
    result
  constructor: (aDefaultRequire)->
    InheritsEx.requireClass = aDefaultRequire if isFunction aDefaultRequire
    return InheritsEx.execute

