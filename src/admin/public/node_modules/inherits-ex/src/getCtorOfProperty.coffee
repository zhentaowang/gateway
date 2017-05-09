getPrototypeOf  = require './getPrototypeOf'

module.exports = (aClass, aPropertyName)->
  vPrototype = aClass::
  while vPrototype and !vPrototype.hasOwnProperty aPropertyName
    vPrototype = getPrototypeOf vPrototype
  if vPrototype
    if vPrototype.hasOwnProperty 'Class'
      result = vPrototype.Class
    else
      result = vPrototype.constructor
  result
