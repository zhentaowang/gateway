"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.setDirname = setDirname;

exports.default = function (loc, relative) {
  var result = undefined;
  [relative].concat(global.babel_core_resolve_enhance_dirnames || []).forEach(function (dirname) {
    if (!result) {
      result = oldResolve(loc, dirname);
    }
  });
  return result;
};

var _module = require("module");

var _module2 = _interopRequireDefault(_module);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _typeof(obj) { return obj && typeof Symbol !== "undefined" && obj.constructor === Symbol ? "symbol" : typeof obj; }

var relativeModules = {};

function oldResolve(loc, relative) {
  // we're in the browser, probably
  if ((typeof _module2.default === "undefined" ? "undefined" : _typeof(_module2.default)) === "object") return null;

  var relativeMod = relativeModules[relative];

  if (!relativeMod) {
    relativeMod = new _module2.default();
    relativeMod.paths = _module2.default._nodeModulePaths(relative);
    relativeModules[relative] = relativeMod;
  }

  try {
    return _module2.default._resolveFilename(loc, relativeMod);
  } catch (err) {
    return null;
  }
}

function setDirname(dirname) {
  global.babel_core_resolve_enhance_dirnames = global.babel_core_resolve_enhance_dirnames || [];
  global.babel_core_resolve_enhance_dirnames.push(dirname);
}