chai            = require 'chai'
sinon           = require 'sinon'
sinonChai       = require 'sinon-chai'
assert          = chai.assert
should          = chai.should()
chai.use(sinonChai)

getPrototypeOf  = require '../src/getPrototypeOf'

describe "getPrototypeOf", ->
  it "should get PrototypeOf object", ->
    class Abc
    obj = new Abc
    result = getPrototypeOf(obj)
    should.exist result
    result.should.be.equal Abc::
