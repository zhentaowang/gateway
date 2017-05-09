var isArray           = Array.isArray;
var isInheritedFrom   = require('./isInheritedFrom');
var inheritsDirectly  = require('./inheritsDirectly');
/**
 * Inherit the prototype methods from one constructor into another.
 *
 *
 * The Function.prototype.inherits from lang.js rewritten as a standalone
 * function (not on Function.prototype). NOTE: If this file is to be loaded
 * during bootstrapping this function needs to be rewritten using some native
 * functions as prototype setup using normal JavaScript does not work as
 * expected during bootstrapping (see mirror.js in r114903).
 *
 * @param {function} ctor Constructor function which needs to inherit the
 *     prototype.
 * @param {function} superCtor Constructor function to inherit prototype from.
 */
function inherits(ctor, superCtor) {
  var v  = ctor.super_;
  var mixinCtor = ctor.mixinCtor_;
  if (mixinCtor && v === mixinCtor) {
    ctor = mixinCtor;
    v = ctor.super_;
  }
  var result = false;
  if (!isInheritedFrom(ctor, superCtor) && !isInheritedFrom(superCtor, ctor)) {
    inheritsDirectly(ctor, superCtor);
    while (v != null && superCtor !== v) {
      ctor = superCtor;
      superCtor = v;
      inheritsDirectly(ctor, superCtor);
      v = ctor.super_;
    }
    result = true;
  }
  return result;
}

module.exports = function(ctor, superCtors) {
  if (!isArray(superCtors)) return inherits(ctor, superCtors);
  for (var i = superCtors.length - 1; i >= 0; i--) {
    if (!inherits(ctor, superCtors[i])) return false;
  }
  return true;
}
