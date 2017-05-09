
module.exports = getProtoChain = (ctor, deepth)->
  deepth?=0
  if deepth >= getProtoChain.maxDeepth
    throw new Error('The maximum depth of nesting is reached.')
  result = while ctor
    if lastCtor and mctors=lastCtor.mixinCtors_
      mctors = mctors.map (m)->getProtoChain m, ++deepth
      name = mctors
    else
      if ctor is Object
        name = "Base"
      else
        name = ctor.name
    lastCtor = ctor
    ctor = ctor.super_
    name
  result.reverse()

getProtoChain.maxDeepth = 10