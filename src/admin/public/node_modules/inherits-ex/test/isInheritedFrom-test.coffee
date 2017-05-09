chai            = require 'chai'
sinon           = require 'sinon'
sinonChai       = require 'sinon-chai'
assert          = chai.assert
should          = chai.should()
chai.use(sinonChai)

inherits        = require '../src/inherits'
isInheritedFrom = require '../src/isInheritedFrom'


describe "isInheritedFrom", ->
  it "should check self circular", ->
    class A
    isInheritedFrom(A, A, false).should.be.equal true
    isInheritedFrom(A, 'A', false).should.be.equal true
  it "should check dead circular", ->
    class A
    class B
    class C
    class D
    A.super_ = B
    B.super_ = C
    C.super_ = A
    isInheritedFrom(B, D, false).should.be.equal true
    isInheritedFrom(A, D, false).should.be.equal true
    isInheritedFrom(B, 'D', false).should.be.equal true
    isInheritedFrom(A, 'D', false).should.be.equal true

