isArray = Array.isArray
isString = (v)-> typeof v == 'string'
isObject = (v)-> v && typeof v == 'object'

getClass = (aClassName, aScope, aValues)->
  if isArray aScope
    if !isArray aValues
      aValues = aScope
      aScope = aValues.map (k)->
        result = k.name
        throw TypeError 'No Scope Name for ' + k unless result?
        result
  else if isObject aScope
    vKeys = Object.keys(aScope)
    aValues = vKeys.map (k)->aScope[k]
    aScope = vKeys
  else
  Function aScope, 'return ' + aClassName
  .apply null, aValues

module.exports = (aClassName, aScope, aValues)->
  if aClassName?
    if typeof aClassName is 'function'
      result = aClassName
    else if typeof aClassName is 'string'
      result = getClass(aClassName, aScope, aValues) unless /[, {}();.]+/.test aClassName
      #result = eval aClassName
  result