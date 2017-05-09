getPrototypeOf = Object.getPrototypeOf
unless getPrototypeOf
  getPrototypeOf = (obj)->obj.__proto__

module.exports = getPrototypeOf