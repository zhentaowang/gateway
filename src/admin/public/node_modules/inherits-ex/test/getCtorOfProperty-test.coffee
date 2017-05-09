chai            = require 'chai'
sinon           = require 'sinon'
sinonChai       = require 'sinon-chai'
assert          = chai.assert
should          = chai.should()
chai.use(sinonChai)

getCtorOfProperty = require '../src/getCtorOfProperty'
isInheritedFrom   = require '../src/isInheritedFrom'

describe "getCtorOfProperty", ->
  it "should get the constructor which owned the property", ->
    class A
      a: 1
    class B extends A
      b: 2
    class C extends B
      c: 3
    result = getCtorOfProperty C, 'c'
    result.should.be.equal C
    result = getCtorOfProperty C, 'b'
    result.should.be.equal B
    result = getCtorOfProperty C, 'a'
    result.should.be.equal A
    result = getCtorOfProperty C, 'no'
    should.not.exist result
