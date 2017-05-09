setPrototypeOf = Object.setPrototypeOf
unless setPrototypeOf
  setPrototypeOf = (obj, prototype)->obj.__proto__ = prototype

module.exports = setPrototypeOf