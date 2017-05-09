chai            = require 'chai'
sinon           = require 'sinon'
sinonChai       = require 'sinon-chai'
assert          = chai.assert
should          = chai.should()

isEmptyFunction = require "../src/isEmptyFunction"

chai.use(sinonChai)

describe "isEmptyFunction", ->

  it "should test empty functions", ->
    emptyFunc = ->
    isEmptyFunction(emptyFunc).should.be.true "emptyFunc"
    emptyFunc = (abc, ase)->
    isEmptyFunction(emptyFunc).should.be.true "emptyFunc2"
    emptyFunc = Function('arg1', 'arg2', '\n;')
    isEmptyFunction(emptyFunc).should.be.true
    emptyFunc = Function('arg1', 'arg2', 'arg3', 'abs;')
    isEmptyFunction(emptyFunc).should.not.be.true

  it "should support istanbul hooked empty function", ->
    `function Test(location){__cov_jOpjHgc_dBxBGFBAhmJ5rg.f['5']++;}`
    isEmptyFunction(Test, true).should.be.true "istanbul"
