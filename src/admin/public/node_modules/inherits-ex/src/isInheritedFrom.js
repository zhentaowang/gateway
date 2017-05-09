var isInheritedFromStr = require('./isInheritedFromStr');

module.exports = function(ctor, superCtor, throwError) {
  if (typeof superCtor === 'string') return isInheritedFromStr(ctor, superCtor, throwError);
  if (ctor === superCtor) {
    if (throwError)
      throw new Error('Circular inherits found!');
    else
      return true;
  }
  var result  = ctor.super_ === superCtor;
  var checkeds = [];
  checkeds.push(ctor);
  while (!result && ((ctor = ctor.super_) != null)) {
    if (checkeds.indexOf(ctor) >= 0) {
      if (throwError)
        throw new Error('Circular inherits found!');
      else
        return true;
    }
    checkeds.push(ctor);
    result = ctor.super_ === superCtor;
  }
  if (result) {
    result = ctor;
    ctor = checkeds[0];
    if (ctor.mixinCtor_ === result) result = ctor;
  }

  return result;
}
