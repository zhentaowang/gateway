module.exports = function(ctor, superStr, throwError) {
  if (ctor.name === superStr) {
    if (throwError)
      throw new Error('Circular inherits found!');
    else
      return true;
  }
  var result  =  ctor.super_ != null && ctor.super_.name === superStr;
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
    result = ctor.super_ != null && ctor.super_.name === superStr;
  }
  if (result) {
    result = ctor;
    ctor = checkeds[0];
    if (ctor.mixinCtor_ === result) result = ctor;
  }

  return result;
};
