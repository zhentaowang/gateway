chai            = require 'chai'
sinon           = require 'sinon'
sinonChai       = require 'sinon-chai'
assert          = chai.assert
should          = chai.should()
chai.use(sinonChai)

createCtor      = require '../src/createCtor'
isInheritedFrom = require '../src/isInheritedFrom'

describe "createCtor", ->
  it "should create a contructor function", ->
    ctor = createCtor "MyClass"
    should.exist ctor, "ctor"
    ctor.should.have.property 'name', 'MyClass'
    ctor.toString().should.have.string 'return MyClass.__super__.constructor.apply(this, arguments)'
  it "should create a contructor function with arguments", ->
    ctor = createCtor "MyClass", ['arg1', 'arg2']
    should.exist ctor, "ctor"
    ctor.should.have.property 'name', 'MyClass'
    s = ctor.toString()
    s.should.have.string 'return MyClass.__super__.constructor.apply(this, arguments)'
    s.should.have.string 'arg1, arg2'
  it "should create a contructor function with body", ->
    ctor = createCtor "MyClass", "return 13414;"
    should.exist ctor, "ctor"
    ctor.should.have.property 'name', 'MyClass'
    s = ctor.toString()
    s.should.have.string 'return 13414;'
  it "should create a contructor function with arguments and body", ->
    ctor = createCtor "MyClass", ['arg1', 'arg2'], "return 13414;"
    should.exist ctor, "ctor"
    ctor.should.have.property 'name', 'MyClass'
    s = ctor.toString()
    s.should.have.string 'return 13414;'
    s.should.have.string 'arg1, arg2'

