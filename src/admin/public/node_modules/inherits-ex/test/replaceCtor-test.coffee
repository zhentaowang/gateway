chai            = require 'chai'
sinon           = require 'sinon'
sinonChai       = require 'sinon-chai'
assert          = chai.assert
should          = chai.should()
chai.use(sinonChai)

replaceCtor     = require '../src/replaceCtor'

describe "replaceCtor", ->
  it "should replace an object's constructor", ->
    class A
    class B
    result = new A
    replaceCtor result, B
    result.should.have.property 'Class', B
    result.should.be.instanceOf B
