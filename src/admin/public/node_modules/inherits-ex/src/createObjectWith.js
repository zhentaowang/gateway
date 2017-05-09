var defineProperty = Object.defineProperty;

module.exports = function(aClass, aArguments) {
  var args = [aClass];
  if (aArguments)
    args = args.concat(Array.prototype.slice.call(aArguments));
  var result = new (Function.prototype.bind.apply(aClass, args));
  if (aClass !== Object && aClass !== Array && aClass !== RegExp) {
    if (!aClass.prototype.hasOwnProperty('Class')) {
      defineProperty(aClass.prototype, 'Class', {
        value: aClass,
        configurable: true
      });
    }
    if (aClass !== aClass.prototype.constructor)
      aClass.prototype.constructor.apply(result, aArguments);
  }
  return result;
};
