var inheritsDirectly  = require('./inheritsDirectly');
var isInheritedFrom   = require('./isInheritedFrom');
var isMixinedFrom     = require('./isMixinedFrom');
var defineProperty    = require('./defineProperty');

var getOwnPropertyNames = Object.getOwnPropertyNames;

/**
 *  Mixin multi classes to ctor.
 *  mixin(Class, ParentClass1, ParentClass2, ...)
 *  + __mixin_ctors__ array to keep the mixined super ctors
 *  + anonymous mixin_ctor hook to super_.
 *  inject into methods to implement inherit.
 *
 *  A11 -> A1 -> A -> Root
 *  B11 -> B1 -> B -> Root
 *  C11 -> C1 -> C -> Root
 *  mixin B11, C
 *  clone C.prototype to MixinCtor.prototype
 *  for k,method of C.prototype
 *    originalMethod = MixinCtor.prototype[k]
 *    if isFunction(originalMethod) and originalMethod.__mixin_prototype__
 *      #B11.__super__ is MixinCtor.prototype
 *      method  = ->
 *        B11.__super__ = originalMethod.__mixin_prototype__
 *        method.apply this, arguments
 *        B11.__super__ = MixinCtor
 *      method.__mixin_prototype__ = C.prototype
 *  B11 -> MixinCtor -> B1 -> B -> Root
 *
mixin the exists method: the new mixin method will oerwrite the old one.

```coffee
class Root
  m: ->
    console.log 'root'
    console.log '----'
class C
  inherits C, Root
  m: ->
    console.log "c"
    super
class B
  inherits B, Root
  m: ->
    console.log "b"
    super
class B11
  inherits B11, B
  m: ->
    console.log 'b11'
    super

b = new B11
b.m()
mixin B11, C
b = new B11
b.m()

# The console results:
# b11
# b
# root
# ----
# b11
# c
# root
# ----

```


 *
 */


function isSuperInFunction(aMethod) {
  return (typeof aMethod === 'function') && aMethod.__mixin_super__ &&
    aMethod.toString().indexOf('__super__') >= 0;
}

function mixin(ctor, superCtor) {
  function clonePrototype(dest, src) {
    var sp = src.prototype;
    var dp = dest.prototype;
    var names = getOwnPropertyNames(sp);
    function _mixin_gen_method(origM, newM, src) {
      var oldSuper = src.__super__;
      return function() {
        src.__super__ = origM.__mixin_super__;
        var result = newM.apply(this, arguments);
        src.__super__ = oldSuper;
        return result;
      };
    }
    for (var i = 1; i < names.length; i++ ) { //i = 1 to skip constructor property
      var k = names[i];
      var method = sp[k];
      var originalMethod = dp[k];
      if (isSuperInFunction(originalMethod) && sp !== originalMethod.__mixin_super__) {
        method = _mixin_gen_method(originalMethod, method, src);
      }
      if (typeof method === 'function') method.__mixin_super__ = sp;
      dp[k] = method;
    }
  }
  var v  = ctor.super_;
  var result = false;
  if (!isMixinedFrom(ctor, superCtor) && !isInheritedFrom(ctor, superCtor) && !isInheritedFrom(superCtor, ctor)) {
    var mixinCtor = ctor.mixinCtor_;
    var mixinCtors = ctor.mixinCtors_;
    if (!mixinCtor) {
      mixinCtor = function MixinCtor_(){};
      defineProperty(ctor, 'mixinCtor_', mixinCtor);
      if (v) inheritsDirectly(mixinCtor, v);
    }
    if (!mixinCtors) {
      mixinCtors = [];
      defineProperty(ctor, 'mixinCtors_', mixinCtors);
    }
    mixinCtors.push(superCtor);//quickly check in isMixinedFrom.
    clonePrototype(mixinCtor, superCtor);
    inheritsDirectly(ctor, mixinCtor);
    result = true;
  }
  return result;
}

module.exports = function(ctor, superCtors, options) {
  if (typeof superCtors === 'function') return mixin(ctor, superCtors, options);
  for (var i = 0; i < superCtors.length; i++) {
    var superCtor = superCtors[i];
    if (!mixin(ctor, superCtor, options)) return false;
  }
  return true;
};
